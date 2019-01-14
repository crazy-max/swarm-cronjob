package app

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// ServiceEvent represents attributes of a Docker service event
type ServiceEvent struct {
	Service     string `mapstructure:"name"`
	UpdateState struct {
		Old string `mapstructure:"updatestate.old"`
		New string `mapstructure:"updatestate.new"`
	} `mapstructure:",squash"`
}

// DockerEnvClient initializes a new Docker API client based on environment variables
func DockerEnvClient() (*client.Client, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	_, err = c.ServerVersion(context.Background())
	return c, err
}

// ScheduledServices returns the list of scheduled Docker services based on swarm-cronjob labels
func ScheduledServices(c *client.Client) ([]swarm.Service, error) {
	svcFilters := filters.NewArgs()
	svcFilters.Add("label", "swarm.cronjob.enable")
	svcFilters.Add("label", "swarm.cronjob.schedule")

	services, err := c.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: svcFilters,
	})
	if err != nil {
		return nil, err
	}

	return services, nil
}

// Service returns a Docker service
func Service(c *client.Client, name string) (swarm.Service, error) {
	svcFilters := filters.NewArgs()
	svcFilters.Add("name", name)

	services, err := c.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: svcFilters,
	})
	if services == nil || len(services) == 0 {
		return swarm.Service{}, errors.New("No matching service found for " + name)
	}

	return services[0], err
}

// RunService runs a cron based service
func RunService(c *client.Client, name string, skipRunning bool) {
	service, _, err := c.ServiceInspectWithRaw(context.Background(), name)
	if err != nil {
		Logger.Error().Err(err).Msgf("Cannot inspect service %s", name)
	}

	svcExitCode, svcStatus := ServiceStatus(c, service.ID)

	if skipRunning && svcStatus == "running" {
		Logger.Warn().Msgf("Skip %s (exit %d ; %s)", name, svcExitCode, svcStatus)
		return
	}

	Logger.Info().Msgf("Start %s (exit %d ; %s)", name, svcExitCode, svcStatus)

	serviceSpec := service.Spec
	*serviceSpec.Mode.Replicated.Replicas = 1                    // Only 1 replica is necessary
	serviceSpec.TaskTemplate.ForceUpdate = service.Version.Index // Set ForceUpdate with Version to ensure update
	_, err = c.ServiceUpdate(context.Background(), service.ID, service.Version, serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		Logger.Error().Err(err).Msgf("Cannot update service %s", name)
	}
}

// ServiceStatus returns service exit code and status
func ServiceStatus(c *client.Client, id string) (int, string) {
	taskFilter := filters.NewArgs()
	taskFilter.Add("service", id)

	tasks, _ := c.TaskList(context.Background(), types.TaskListOptions{
		Filters: taskFilter,
	})

	exitCode := 1
	status := ""
	stopStates := []swarm.TaskState{
		swarm.TaskStateComplete,
		swarm.TaskStateFailed,
		swarm.TaskStateRejected,
	}

	for _, task := range tasks {
		stop := false
		for _, stopState := range stopStates {
			if task.Status.State == stopState {
				stop = true
				break
			}
		}
		status = string(task.Status.State)
		if stop {
			exitCode = task.Status.ContainerStatus.ExitCode
			if exitCode == 0 && task.Status.State == swarm.TaskStateRejected {
				exitCode = 255 // force non-zero exit for task rejected
			}
		}
		break
	}

	return exitCode, status
}
