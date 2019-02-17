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
		log.Error().Err(err).Msgf("Cannot inspect service %s", c.Job.Name)
	}

	svcExitCode, svcStatus := c.Docker.ServiceStatus(service.ID)

	if c.Job.SkipRunning && svcStatus == "running" {
		log.Warn().Msgf("Skip %s (exit %d ; %s)", c.Job.Name, svcExitCode, svcStatus)
		return
	}

	log.Info().Msgf("Start %s (exit %d ; %s)", c.Job.Name, svcExitCode, svcStatus)

	serviceSpec := service.Spec
	*serviceSpec.Mode.Replicated.Replicas = 1                    // Only 1 replica is necessary
	serviceSpec.TaskTemplate.ForceUpdate = service.Version.Index // Set ForceUpdate with Version to ensure update
	_, err = c.Docker.ServiceUpdate(context.Background(), service.ID, service.Version, serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("Cannot update service %s", c.Job.Name)
	}
}
