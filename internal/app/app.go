package app

import (
	"context"
	"slices"
	"strconv"
	"time"

	"github.com/crazy-max/swarm-cronjob/internal/docker"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/crazy-max/swarm-cronjob/internal/worker"
	"github.com/go-viper/mapstructure/v2"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
	cron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// Periodic reconciliation keeps the in-memory schedule aligned with Swarm
// even if a service update event is missed or lacks enough attributes.
const reconcileInterval = time.Minute

// SwarmCronjob represents an active swarm-cronjob object
type SwarmCronjob struct {
	docker docker.Client
	cron   *cron.Cron
	jobs   map[string]cron.EntryID
}

// New creates new swarm-cronjob instance
func New() (*SwarmCronjob, error) {
	log.Debug().Msg("Creating Docker API client")
	d, err := docker.NewEnvClient()

	return &SwarmCronjob{
		docker: d,
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		)),
		jobs: make(map[string]cron.EntryID),
	}, err
}

func (sc *SwarmCronjob) Run(ctx context.Context) error {
	if err := sc.reconcileJobs(); err != nil {
		return err
	}

	log.Debug().Msg("Starting the cron scheduler")
	sc.cron.Start()
	defer func() {
		<-sc.cron.Stop().Done()
	}()
	sc.logScheduledJobs()

	log.Debug().Msg("Listening docker events...")
	filter := make(client.Filters).Add("type", "service")
	reconcileTicker := time.NewTicker(reconcileInterval)
	defer reconcileTicker.Stop()

	msgs, errs := sc.docker.Events(ctx, client.EventsListOptions{
		Filters: filter,
	})

	for msgs != nil || errs != nil {
		select {
		case <-ctx.Done():
			return nil
		case <-reconcileTicker.C:
			if err := sc.reconcileJobs(); err != nil {
				log.Error().Err(err).Msg("Cannot reconcile cronjobs")
			}
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err == nil {
				continue
			}
			if errors.Is(err, context.Canceled) && context.Cause(ctx) != nil {
				return nil
			}
			return errors.Wrap(err, "event channel failed")
		case msg, ok := <-msgs:
			if !ok {
				msgs = nil
				continue
			}
			event := model.ServiceEvent{}
			err := mapstructure.Decode(msg.Actor.Attributes, &event)
			if err != nil {
				log.Warn().Msgf("Cannot decode event, %v", err)
				continue
			}
			if event.Service == "" {
				log.Debug().Msg("Service event missing name, reconciling all cronjobs")
				if err := sc.reconcileJobs(); err != nil {
					log.Error().Err(err).Msg("Cannot reconcile cronjobs")
				}
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
			}
			if processed {
				log.Debug().Msgf("Number of cronjob tasks: %d", len(sc.cron.Entries()))
			}
		}
	}

	return nil
}

func (sc *SwarmCronjob) reconcileJobs() error {
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

	desired := make(map[string]struct{}, len(services))
	for _, service := range services {
		desired[service.Name] = struct{}{}
		if _, err := sc.crudJobWithService(service, false); err != nil {
			log.Error().Err(err).Msgf("Cannot manage job for service %s", service.Name)
		}
	}

	existing := make([]string, 0, len(sc.jobs))
	for serviceName := range sc.jobs {
		if _, ok := desired[serviceName]; !ok {
			existing = append(existing, serviceName)
		}
	}
	slices.Sort(existing)
	for _, serviceName := range existing {
		if _, err := sc.crudJob(serviceName); err != nil {
			log.Error().Err(err).Msgf("Cannot manage job for service %s", serviceName)
		}
	}

	return nil
}

func (sc *SwarmCronjob) crudJob(serviceName string) (bool, error) {
	jobID, jobFound := sc.jobs[serviceName]

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

	return sc.crudJobWithService(service, true)
}

func (sc *SwarmCronjob) crudJobWithService(service *model.ServiceInfo, noopCountsAsProcessed bool) (bool, error) {
	jobID, jobFound := sc.jobs[service.Name]
	var err error

	wc := &worker.Client{
		Docker: sc.docker,
		Job: model.Job{
			Name:        service.Name,
			Enable:      false,
			SkipRunning: false,
			Replicas:    1,
		},
	}

	for labelKey, labelValue := range service.Labels {
		switch labelKey {
		case "swarm.cronjob.enable":
			wc.Job.Enable, err = strconv.ParseBool(labelValue)
			if err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			}
		case "swarm.cronjob.schedule":
			wc.Job.Schedule = labelValue
		case "swarm.cronjob.skip-running":
			wc.Job.SkipRunning, err = strconv.ParseBool(labelValue)
			if err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			}
		case "swarm.cronjob.replicas":
			wc.Job.Replicas, err = strconv.ParseUint(labelValue, 10, 64)
			if err != nil {
				log.Error().Str("service", service.Name).Err(err).Msgf("Cannot parse %s value of label %s", labelValue, labelKey)
			} else if wc.Job.Replicas < 1 {
				log.Error().Str("service", service.Name).Msgf("%s must be greater than or equal to one", labelKey)
			}
		case "swarm.cronjob.registry-auth":
			wc.Job.RegistryAuth, err = strconv.ParseBool(labelValue)
			if err != nil {
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
		}
	}

	if !wc.Job.Enable {
		if jobFound {
			log.Info().Str("service", service.Name).Msg("Disable cronjob")
			sc.removeJob(service.Name, jobID)
			return true, nil
		}
		log.Debug().Str("service", service.Name).Msg("Cronjob disabled")
		return false, nil
	}

	if jobFound {
		entry := sc.cron.Entry(jobID)
		if entry.Valid() {
			if workerClient, ok := entry.Job.(*worker.Client); ok && jobEqual(workerClient.Job, wc.Job) {
				return noopCountsAsProcessed, nil
			}
		}
		sc.removeJob(service.Name, jobID)
		log.Debug().Str("service", service.Name).Msgf("Update cronjob with schedule %s", wc.Job.Schedule)
	} else {
		log.Info().Str("service", service.Name).Msgf("Add cronjob with schedule %s", wc.Job.Schedule)
	}

	jobID, err = sc.cron.AddJob(wc.Job.Schedule, wc)
	if err != nil {
		return false, err
	}

	sc.jobs[service.Name] = jobID
	sc.logScheduledJob(service.Name, jobID)
	return true, err
}

func (sc *SwarmCronjob) removeJob(serviceName string, id cron.EntryID) {
	delete(sc.jobs, serviceName)
	sc.cron.Remove(id)
}

func (sc *SwarmCronjob) logScheduledJobs() {
	serviceNames := make([]string, 0, len(sc.jobs))
	for serviceName := range sc.jobs {
		serviceNames = append(serviceNames, serviceName)
	}
	slices.Sort(serviceNames)
	for _, serviceName := range serviceNames {
		sc.logScheduledJob(serviceName, sc.jobs[serviceName])
	}
}

func (sc *SwarmCronjob) logScheduledJob(serviceName string, id cron.EntryID) {
	entry := sc.cron.Entry(id)
	if !entry.Valid() {
		return
	}
	logger := log.Debug().Str("service", serviceName)
	if workerClient, ok := entry.Job.(*worker.Client); ok {
		logger = logger.Str("schedule", workerClient.Job.Schedule)
	}
	if entry.Next.IsZero() {
		return
	}
	logger.Time("next_run", entry.Next.UTC()).Msg("Cronjob scheduled")
}

func jobEqual(a, b model.Job) bool {
	if a.Name != b.Name ||
		a.Enable != b.Enable ||
		a.Schedule != b.Schedule ||
		a.SkipRunning != b.SkipRunning ||
		a.RegistryAuth != b.RegistryAuth ||
		a.Replicas != b.Replicas {
		return false
	}
	if a.QueryRegistry == nil || b.QueryRegistry == nil {
		return a.QueryRegistry == nil && b.QueryRegistry == nil
	}
	return *a.QueryRegistry == *b.QueryRegistry
}
