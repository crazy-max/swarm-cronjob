package app

import (
	"context"
	"strconv"
	"time"

	"github.com/crazy-max/cron"
	"github.com/crazy-max/swarm-cronjob/internal/docker"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/crazy-max/swarm-cronjob/internal/worker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

// SwarmCronjob represents an active swarm-cronjob object
type SwarmCronjob struct {
	docker   *docker.Client
	cron     *cron.Cron
	location *time.Location
}

// New creates new swarm-cronjob instance
func New(location *time.Location) (*SwarmCronjob, error) {
	log.Debug().Msg("Creating Docker API client")
	d, err := docker.NewEnvClient()

	log.Debug().Msg("Creating Cron job runner")
	c := cron.NewWithLocation(location)

	return &SwarmCronjob{
		docker:   d,
		cron:     c,
		location: location,
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
			log.Debug().Msgf("Event triggered for %s (newstate='%s' oldstate='%s')", event.Service, event.UpdateState.New, event.UpdateState.Old)
			processed, err := sc.crudJob(event.Service)
			if err != nil {
				log.Error().Err(err).Msgf("Cannot manage job for service %s", event.Service)
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
	var jobEntry *cron.Entry
	for _, entry := range sc.cron.Entries() {
		if entry.Name == serviceName {
			jobEntry = entry
			break
		}
	}

	// Check service exists
	service, err := sc.docker.Service(serviceName)
	if err != nil {
		if jobEntry != nil {
			log.Debug().Msgf("Remove cronjob for service %s", serviceName)
			return true, sc.cron.Remove(jobEntry.Name)
		}
		log.Debug().Msgf("Service %s does not exist (removed)", serviceName)
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
				log.Error().Err(err).Msgf("Cannot parse %s value of label swarm.cronjob.enable for service %s", labelKey, service.Spec.Name)
			}
		case "swarm.cronjob.schedule":
			wc.Job.Schedule = labelValue
		case "swarm.cronjob.skip-running":
			wc.Job.SkipRunning, err = strconv.ParseBool(labelValue)
			if err != nil {
				log.Error().Err(err).Msgf("Cannot parse %s value of label swarm.cronjob.skip-running for service %s", labelKey, service.Spec.Name)
			}
		}
	}

	// Disabled or non-cron service
	if !wc.Job.Enable {
		if jobEntry != nil {
			log.Debug().Msgf("Disable cronjob for service %s", service.Spec.Name)
			return true, sc.cron.Remove(jobEntry.Name)
		}
		log.Debug().Msgf("Cronjob disabled for service %s", service.Spec.Name)
		return false, nil
	}

	// Add/Update job
	if jobEntry != nil {
		if err := sc.cron.Remove(jobEntry.Name); err != nil {
			return true, err
		}
		log.Debug().Msgf("Update cronjob for service %s with schedule %s", service.Spec.Name, wc.Job.Schedule)
	} else {
		log.Info().Msgf("Add cronjob for service %s with schedule %s", service.Spec.Name, wc.Job.Schedule)
	}

	return true, sc.cron.AddJob(wc.Job.Schedule, wc, wc.Job.Name)
}

// Close closes swarm-cronjob
func (sc *SwarmCronjob) Close() {
	if sc.cron != nil {
		sc.cron.Stop()
	}
}
