package app

import (
	"context"
	"strconv"

	"github.com/crazy-max/swarm-cronjob/internal/docker"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/crazy-max/swarm-cronjob/internal/scheduler"
	"github.com/crazy-max/swarm-cronjob/internal/worker"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type (
	Jobs map[string]scheduler.Uid

	// SwarmCronjob represents an active swarm-cronjob object
	SwarmCronjob struct {
		docker    docker.Client
		scheduler *scheduler.Scheduler
		jobs      Jobs
		runOnce   map[string]bool
	}
)

// New setup new swarm-cronjob instance
func New() (*SwarmCronjob, error) {
	log.Debug().Msg("Creating Docker API client")
	d, err := docker.NewEnvClient()

	scheduler, err := scheduler.NewScheduler(true)
	if err != nil {
		return nil, err
	}

	return &SwarmCronjob{
		docker:    d,
		scheduler: scheduler,
		jobs:      make(Jobs),
		runOnce:   make(map[string]bool),
	}, err
}

// Run starts swarm-cronjob process
func (sc *SwarmCronjob) Run() error {
	// Find scheduled services
	services, err := sc.docker.ServiceList(&model.ServiceListArgs{
		Labels: []string{
			"swarm.cronjob.enable",
			"swarm.cronjob.schedule",
		},
	})
	if err != nil {
		return err
	}
	log.Debug().Msgf("%d scheduled services found through labels", len(services))

	// Add services as cronjobs
	for _, service := range services {
		if _, err := sc.crudJob(service.Name, false); err != nil {
			log.Error().Err(err).Msgf("Cannot manage job for service %s", service.Name)
		}
	}

	// Start cron routine
	log.Debug().Msg("Starting the cron scheduler")
	sc.scheduler.Start()

	// Listen Docker events
	log.Debug().Msg("Listening docker events...")
	filter := filters.NewArgs()
	filter.Add("type", "service")

	msgs, errs := sc.docker.Events(context.Background(), events.ListOptions{
		Filters: filter,
	})

	var (
		event  model.ServiceEvent
		deploy bool
	)
	for {
		select {
		case err := <-errs:
			log.Fatal().Err(err).Msg("Event channel failed")
		case msg := <-msgs:
			if err := mapstructure.Decode(msg.Actor.Attributes, &event); err != nil {
				log.Warn().Msgf("Cannot decode event, %v", err)
				continue
			}

			log.Debug().
				Str("service", event.Service).
				Str("newstate", event.UpdateState.New).
				Str("oldstate", event.UpdateState.Old).
				Msg("Event triggered")

			if event.UpdateState.New == "" && event.UpdateState.Old == "" {
				deploy = true
			}
			processed, err := sc.crudJob(event.Service, deploy)
			if err != nil {
				log.Error().Str("service", event.Service).Err(err).Msg("Cannot manage job")
				continue

			} else if processed {
				log.Debug().Msgf("Number of cronjob tasks: %d", sc.scheduler.CountJobs())
			}
		}
	}
}

// crudJob adds, updates or removes cron job service
func (sc *SwarmCronjob) crudJob(serviceName string, deploy bool) (bool, error) {
	// Find existing job
	jobID, jobFound := sc.jobs[serviceName]

	// Check service exists
	service, err := sc.docker.Service(serviceName)
	if err != nil {
		if jobFound {
			log.Info().Str("service", serviceName).Msg("Remove cronjob")
			sc.removeJob(serviceName, jobID)
			return true, nil
		}
		log.Debug().Str("service", serviceName).Msg("Service does not exist (removed)")
		return false, nil
	}

	// Cronjob worker
	wc := &worker.Client{
		Docker: sc.docker,
		Job: model.Job{
			Name:        service.Name,
			Enable:      false,
			SkipRunning: false,
			Replicas:    1,
		},
	}

	// Seek swarm.cronjob labels
	for labelKey, labelValue := range service.Labels {
		switch labelKey {
		case "swarm.cronjob.enable":
			if wc.Job.Enable, err = strconv.ParseBool(labelValue); err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			}

		case "swarm.cronjob.schedule":
			wc.Job.Schedule = labelValue

		case "swarm.cronjob.skip-running":
			if wc.Job.SkipRunning, err = strconv.ParseBool(labelValue); err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			}

		case "swarm.cronjob.replicas":
			if wc.Job.Replicas, err = strconv.ParseUint(labelValue, 10, 64); err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			} else if wc.Job.Replicas < 1 {
				log.Error().Str("service", service.Name).Msgf("%s must be greater than or equal to one", labelKey)
			}

		case "swarm.cronjob.registry-auth":
			if wc.Job.RegistryAuth, err = strconv.ParseBool(labelValue); err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			}

		case "swarm.cronjob.query-registry":
			queryRegistry, err := strconv.ParseBool(labelValue)
			if err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			}
			wc.Job.QueryRegistry = &queryRegistry

		case "swarm.cronjob.scaledown":
			if labelValue == "true" {
				log.Debug().Str("service", service.Name).Msg("Scale down detected. Skipping cronjob")
				return false, nil
			}

		case "swarm.cronjob.run-once":
			if labelValue == "true" {
				sc.runOnce[service.Name] = true
				log.Debug().Str("service", service.Name).Msgf("Enabled run once for the job %s", service.Name)
			} else {
				sc.runOnce[service.Name] = false
			}
		}
	}

	// Check if is disabled or is a non-cron service
	if !wc.Job.Enable {
		if jobFound {
			log.Info().Str("service", service.Name).Msg("Disable cronjob")
			sc.removeJob(serviceName, jobID)
			return true, nil
		}
		log.Debug().Str("service", service.Name).Msg("Cronjob disabled")
		return false, nil
	}

	// Add/Update job
	if jobFound {
		// check if is to run job only once
		if sc.runOnce[service.Name] {
			// check if the defined service got an update
			if deploy {
				sc.removeJob(serviceName, jobID)
			}
			log.Info().Str("service", service.Name).Msgf("Job %s only scheduled to run once, skipping", wc.Job.Name)
			return true, err
		}

		sc.removeJob(serviceName, jobID)
		log.Debug().Str("service", service.Name).Msgf("Update cronjob with schedule %s", wc.Job.Schedule)
	} else {
		if sc.runOnce[service.Name] {
			log.Info().Str("service", service.Name).Msgf("Add one time job to be run after %s seconds", wc.Job.Schedule)
		} else {
			log.Info().Str("service", service.Name).Msgf("Add cronjob with schedule %s", wc.Job.Schedule)
		}
	}

	var job scheduler.Job
	// check if current service is configured to run only once
	if sc.runOnce[service.Name] {
		if job, err = sc.scheduler.OneTimeJob(wc.Job.Schedule, func() {
			defer wc.Run()
			log.Debug().Str("service", service.Name).Msgf("Triggered one time job %s ...", wc.Job.Name)
		}); err != nil {
			return false, err
		}
	} else {

		// by default set service as a Cron Job
		if job, err = sc.scheduler.CronJob(wc.Job.Schedule, func() {
			wc.Run()
		}); err != nil {
			return false, err
		}
	}

	sc.jobs[serviceName] = job.ID()
	return true, err
}

// Close stops swarm-cronjob jobs to be processed
func (sc *SwarmCronjob) Close() {
	if sc.scheduler != nil {
		sc.scheduler.Stop()
	}
}

func (sc *SwarmCronjob) removeJob(serviceName string, id scheduler.Uid) {
	defer sc.scheduler.RemoveJob(id)
	delete(sc.jobs, serviceName)
}
