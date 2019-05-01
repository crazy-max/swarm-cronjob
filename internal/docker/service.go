package docker

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/rs/zerolog/log"
)

// ServiceEvent represents attributes of a Docker service event
type ServiceEvent struct {
	Service     string `mapstructure:"name"`
	UpdateState struct {
		Old string `mapstructure:"updatestate.old"`
		New string `mapstructure:"updatestate.new"`
	} `mapstructure:",squash"`
}

// ScheduledServices returns the list of scheduled Docker services based on swarm-cronjob labels
func (c *Client) ScheduledServices() ([]swarm.Service, error) {
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
func (c *Client) Service(name string) (swarm.Service, error) {
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

// ServiceTaskStatus returns latest service task status
func (c *Client) ServiceTaskStatus(id string) string {
	taskFilter := filters.NewArgs()
	taskFilter.Add("service", id)

	tasks, _ := c.TaskList(context.Background(), types.TaskListOptions{
		Filters: taskFilter,
	})

	ltask := swarm.Task{}
	exitCode := -1
	status := "n/a"
	stopStates := []swarm.TaskState{
		swarm.TaskStateComplete,
		swarm.TaskStateFailed,
		swarm.TaskStateRejected,
	}

	for _, task := range tasks {
		// skip completed tasks
		if task.Status.State == swarm.TaskStateComplete {
			continue
		}
		ltask = task
		status = string(task.Status.State)
		stop := false
		for _, stopState := range stopStates {
			if task.Status.State == stopState {
				stop = true
				break
			}
		}
		if stop {
			exitCode = task.Status.ContainerStatus.ExitCode
			if exitCode == 0 && task.Status.State == swarm.TaskStateRejected {
				exitCode = 255 // force non-zero exit for task rejected
			}
		}
		break
	}

	log.Debug().
		Str("node_id", ltask.NodeID).
		Str("service_id", id).
		Str("task_id", ltask.ID).
		Int("exit_code", exitCode).
		Str("status_message", ltask.Status.Message).
		Str("status_state", status).
		Msg("Service task")

	return status
}
