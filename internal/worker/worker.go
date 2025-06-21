package worker

import (
	"context"

	"github.com/crazy-max/swarm-cronjob/internal/docker"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/rs/zerolog/log"
)

// Client represents an active worker object
type Client struct {
	Docker docker.Client
	Job    model.Job
}

// Run runs a cron based service
func (c *Client) Run() {
	service, err := c.Docker.Service(c.Job.Name)
	if err != nil {
		log.Error().Str("service", c.Job.Name).Err(err).Msg("Service not found")
		return
	}
	serviceUp := service.Raw

	tasks, err := c.Docker.TaskList(c.Job.Name)
	if err != nil {
		log.Error().Str("service", c.Job.Name).Err(err).Msg("Cannot list service tasks")
		return
	}
	for _, task := range tasks {
		log.Debug().
			Str("node", task.NodeName).
			Str("service", task.ServiceName).
			Str("task_id", task.ID).
			Str("status_message", task.Status.Message).
			Str("status_state", string(task.Status.State)).
			Msg("Service task")
	}

	if c.Job.SkipRunning && service.Actives > 0 {
		log.Warn().Str("service", c.Job.Name).
			Uint64("tasks_active", service.Actives).
			Msg("Skip running job")
		return
	}

	log.Info().Str("service", c.Job.Name).
		Uint64("tasks_active", service.Actives).
		Str("status", service.UpdateStatus).Msg("Start job")

	// Set number of replicas based on service mode
	if c.Job.Replicas > 1 {
		// Need to scale down service to 0 to fix an issue if replicas > 1
		// See https://github.com/crazy-max/swarm-cronjob/issues/16
		if serviceUp, err = c.scaleDown(serviceUp); err != nil {
			log.Error().Str("service", c.Job.Name).Err(err).Msg("Cannot scaled down")
		}
	}

	switch service.Mode {
	case model.ServiceModeReplicated:
		*serviceUp.Spec.Mode.Replicated.Replicas = c.Job.Replicas
	case model.ServiceModeReplicatedJob:
		*serviceUp.Spec.Mode.ReplicatedJob.MaxConcurrent = c.Job.Replicas
	}

	// Set ForceUpdate with Version to ensure update
	serviceUp.Spec.TaskTemplate.ForceUpdate = serviceUp.Version.Index

	// Update options
	updateOpts := types.ServiceUpdateOptions{}
	if c.Job.RegistryAuth {
		encodedAuth, err := c.Docker.RetrieveAuthTokenFromImage(context.Background(), serviceUp.Spec.TaskTemplate.ContainerSpec.Image)
		if err != nil {
			log.Error().Err(err).Str("service", c.Job.Name).Msg("Cannot retrieve auth token from service's image")
			return
		}
		if encodedAuth != "e30=" {
			updateOpts.EncodedRegistryAuth = encodedAuth
		}
	} else {
		updateOpts.RegistryAuthFrom = types.RegistryAuthFromSpec
	}
	if c.Job.QueryRegistry != nil {
		updateOpts.QueryRegistry = *c.Job.QueryRegistry
	}

	// Update service
	response, err := c.Docker.ServiceUpdate(context.Background(), serviceUp.ID, serviceUp.Version, serviceUp.Spec, updateOpts)
	if err != nil {
		log.Error().Str("service", c.Job.Name).Err(err).Msg("Cannot update")
	}
	for _, warn := range response.Warnings {
		log.Warn().Str("service", c.Job.Name).Msg(warn)
	}
}

func (c *Client) scaleDown(serviceRaw swarm.Service) (swarm.Service, error) {
	switch {
	case serviceRaw.Spec.Mode.Replicated != nil:
		*serviceRaw.Spec.Mode.Replicated.Replicas = 0
	case serviceRaw.Spec.Mode.ReplicatedJob != nil:
		*serviceRaw.Spec.Mode.ReplicatedJob.MaxConcurrent = 0
	}
	serviceRaw.Spec.Labels["swarm.cronjob.scaledown"] = "true"
	serviceRaw.Spec.TaskTemplate.ForceUpdate = serviceRaw.Version.Index

	_, err := c.Docker.ServiceUpdate(context.Background(), serviceRaw.ID, serviceRaw.Version, serviceRaw.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		return swarm.Service{}, err
	}

	service, err := c.Docker.Service(c.Job.Name)
	if err != nil {
		return swarm.Service{}, err
	}

	delete(service.Raw.Spec.Labels, "swarm.cronjob.scaledown")
	return service.Raw, nil
}
