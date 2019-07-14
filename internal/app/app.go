package app

import (
	"context"
	"strconv"
	"time"

	"github.com/crazy-max/swarm-cronjob/internal/docker"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/crazy-max/swarm-cronjob/internal/worker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/hako/durafmt"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// SwarmCronjob represents an active swarm-cronjob object
type SwarmCronjob struct {
	docker   *docker.Client
	cron     *cron.Cron
	location *time.Location
	jobs     map[string]cron.EntryID
}

// New creates new swarm-cronjob instance
func New(location *time.Location) (*SwarmCronjob, error) {
	log.Debug().Msg("Creating Docker API client")
	d, err := docker.NewEnvClient()

	return &SwarmCronjob{
		docker:   d,
		cron:     cron.New(cron.WithLocation(location), cron.WithSeconds()),
		location: location,
		jobs:     make(map[string]cron.EntryID),
	}, err
}

// Run starts swarm-cronjob process
func (sc *SwarmCronjob) Run() error {
	// Find scheduled services
	services, err := sc.docker.ScheduledServices()
	if err != nil {
		return err
	}
	log.Debug().Msgf("%d scheduled services found through labels", len(services))

	// Add services as cronjobs
	for _, service := range services {
		if _, err := sc.crudJob(service.Spec.Name); err != nil {
			log.Error().Err(err).Msgf("Cannot manage job for service %s", service.Spec.Name)
		}
	}

	// Start cron routine
	log.Debug().Msg("Starting the cron scheduler")
	sc.cron.Start()

	// Listen Docker events
	log.Debug().Msg("Listening docker events...")
	filter := filters.NewArgs()
	filter.Add("type", "service")

	msgs, errs := sc.docker.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})

	var event docker.ServiceEvent
	for {
		select {
		case err := <-errs:
			log.Fatal().Err(err).Msg("Event channel failed")
		case msg := <-msgs:
			err := mapstructure.Decode(msg.Actor.Attributes, &event)
			if err != nil {
				log.Warn().Msgf("Cannot decode event, %v", err)
				continue
			}
			log.Debug().
				Str("service", event.Service).
				Str("newstate", event.UpdateState.New).
				Str("oldstate", event.UpdateState.Old).
				Msg("Event triggered")
			processed, err := sc.crudJob(event.Service)
			if err != nil {
				log.Error().Str("service", event.Service).Err(err).Msg("Cannot manage job")
				continue
			} else if processed {
				log.Debug().Msgf("Number of cronjob tasks : %d", len(sc.cron.Entries()))
			}
		}
	}
}

// crudJob adds, updates or removes cron job service
func (sc *SwarmCronjob) crudJob(serviceName string) (bool, error) {
	// Find existing job
	jobID, jobFound := sc.jobs[serviceName]

	// Check service exists
	service, err := sc.docker.Service(serviceName)
	if err != nil {
		if jobFound {
			log.Debug().Str("service", serviceName).Msg("Remove cronjob")
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
			Name:        service.Spec.Name,
			Enable:      false,
			SkipRunning: false,
		},
	}

	// Seek swarm.cronjob labels
	for labelKey, labelValue := range service.Spec.Labels {
		switch labelKey {
		case "swarm.cronjob.enable":
			wc.Job.Enable, err = strconv.ParseBool(labelValue)
			if err != nil {
				log.Error().Str("service", service.Spec.Name).Err(err).Msgf("Cannot parse %s value of label swarm.cronjob.enable", labelKey)
			}
		case "swarm.cronjob.schedule":
			wc.Job.Schedule = labelValue
		case "swarm.cronjob.skip-running":
			wc.Job.SkipRunning, err = strconv.ParseBool(labelValue)
			if err != nil {
				log.Error().Str("service", service.Spec.Name).Err(err).Msgf("Cannot parse %s value of label swarm.cronjob.skip-running", labelKey)
			}
		}
	}

	// Disabled or non-cron service
	if !wc.Job.Enable {
		if jobFound {
			log.Debug().Str("service", service.Spec.Name).Msg("Disable cronjob")
			sc.removeJob(serviceName, jobID)
			return true, nil
		}
		log.Debug().Str("service", service.Spec.Name).Msg("Cronjob disabled")
		return false, nil
	}

	// Add/Update job
	if jobFound {
		sc.removeJob(serviceName, jobID)
		log.Debug().Str("service", service.Spec.Name).Msgf("Update cronjob with schedule %s", wc.Job.Schedule)
	} else {
		log.Info().Str("service", service.Spec.Name).Msgf("Add cronjob with schedule %s", wc.Job.Schedule)
	}

	sc.jobs[serviceName], err = sc.cron.AddJob(wc.Job.Schedule, wc)
	log.Info().Str("service", service.Spec.Name).Msgf("Next run in %s (%s)",
		durafmt.ParseShort(sc.cron.Entry(jobID).Next.Sub(time.Now())).String(),
		sc.cron.Entry(jobID).Next)

	return true, err
}

// Close closes swarm-cronjob
func (sc *SwarmCronjob) Close() {
	if sc.cron != nil {
		sc.cron.Stop()
	}
}

func (sc *SwarmCronjob) removeJob(serviceName string, id cron.EntryID) {
	delete(sc.jobs, serviceName)
	sc.cron.Remove(id)
}
