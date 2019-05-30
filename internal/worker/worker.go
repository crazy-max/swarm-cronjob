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
	service, _, err := c.Docker.ServiceInspectWithRaw(context.Background(), c.Job.Name, types.ServiceInspectOptions{})
	if err != nil {
		log.Error().Str("service", c.Job.Name).Err(err).Msg("Cannot inspect")
	}

	status := c.Docker.ServiceTaskStatus(service.ID)

	if c.Job.SkipRunning && status == "running" {
		log.Warn().Str("service", c.Job.Name).Str("last_status", status).Msg("Skip job")
		return
	}

	if service.Spec.Mode.Replicated == nil {
		log.Error().Str("service", c.Job.Name).Err(err).Msg("Only replicated mode is supported")
		return
	}

	log.Info().Str("service", c.Job.Name).Str("last_status", status).Msg("Start job")

	*service.Spec.Mode.Replicated.Replicas = 1                    // Only 1 replica is necessary
	service.Spec.TaskTemplate.ForceUpdate = service.Version.Index // Set ForceUpdate with Version to ensure update
	_, err = c.Docker.ServiceUpdate(context.Background(), service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		log.Error().Str("service", c.Job.Name).Err(err).Msg("Cannot update")
	}
}
