package worker

import (
	"context"

	"github.com/crazy-max/swarm-cronjob/internal/docker"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/docker/api/types"
	"github.com/rs/zerolog/log"
)

// Client represents an active worker object
type Client struct {
	Docker *docker.Client
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

	// Only 1 replica is necessary in replicated mode
	if service.Mode == model.ServiceModeReplicated {
		*serviceUp.Spec.Mode.Replicated.Replicas = 1
	}

	// Set ForceUpdate with Version to ensure update
	serviceUp.Spec.TaskTemplate.ForceUpdate = serviceUp.Version.Index

	// Update service
	_, err = c.Docker.Cli.ServiceUpdate(context.Background(), serviceUp.ID, serviceUp.Version, serviceUp.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		log.Error().Str("service", c.Job.Name).Err(err).Msg("Cannot update")
	}
}
