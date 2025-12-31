package docker

import (
	"context"
	"sort"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/pkg/errors"
)

// ServiceList return all services.
func (c *DockerClient) ServiceList(args *model.ServiceListArgs) ([]*model.ServiceInfo, error) {
	opts := swarm.ServiceListOptions{
		Filters: filters.NewArgs(),
	}
	if args.Name != "" {
		opts.Filters.Add("name", args.Name)
	}
	if len(args.Labels) > 0 {
		for _, label := range args.Labels {
			opts.Filters.Add("label", label)
		}
	}

	services, err := c.api.ServiceList(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].Spec.Name < services[j].Spec.Name
	})

	// nodes
	nodes, err := c.api.NodeList(context.Background(), swarm.NodeListOptions{})
	if err != nil {
		return nil, err
	}
	activeNodes := make(map[string]struct{})
	for _, node := range nodes {
		if node.Status.State != swarm.NodeStateDown {
			activeNodes[node.ID] = struct{}{}
		}
	}

	// tasks
	taskOpts := swarm.TaskListOptions{
		Filters: filters.NewArgs(),
	}
	for _, service := range services {
		taskOpts.Filters.Add("service", service.ID)
	}
	tasks, err := c.api.TaskList(context.Background(), taskOpts)
	if err != nil {
		return nil, err
	}

	// active tasks
	running, tasksNoShutdown := map[string]uint64{}, map[string]uint64{}
	for _, task := range tasks {
		if task.DesiredState != swarm.TaskStateShutdown {
			tasksNoShutdown[task.ServiceID]++
		}
		if _, nodeActive := activeNodes[task.NodeID]; nodeActive && task.Status.State == swarm.TaskStateRunning {
			running[task.ServiceID]++
		}
	}

	// res
	res := make([]*model.ServiceInfo, len(services))
	for i, service := range services {
		res[i] = &model.ServiceInfo{
			Raw:       service,
			ID:        service.ID,
			Name:      service.Spec.Name,
			Image:     normalizeImage(service.Spec.TaskTemplate.ContainerSpec.Image),
			Labels:    service.Spec.Labels,
			Actives:   running[service.ID],
			UpdatedAt: service.UpdatedAt.Local(),
			Rollback:  service.PreviousSpec != nil,
		}
		if service.UpdateStatus != nil {
			res[i].UpdateStatus = string(service.UpdateStatus.State)
		}
		if service.Spec.Mode.Replicated != nil && service.Spec.Mode.Replicated.Replicas != nil {
			res[i].Mode = model.ServiceModeReplicated
			res[i].Replicas = *service.Spec.Mode.Replicated.Replicas
		} else if service.Spec.Mode.Global != nil {
			res[i].Mode = model.ServiceModeGlobal
			res[i].Replicas = tasksNoShutdown[service.ID]
		}
	}

	return res, nil
}

// Service returns a service
func (c *DockerClient) Service(name string) (*model.ServiceInfo, error) {
	service, err := c.ServiceList(&model.ServiceListArgs{
		Name: name,
	})
	if err != nil {
		return nil, err
	} else if len(service) == 0 {
		return nil, errors.Errorf("%s service not found", name)
	}

	return service[0], nil
}
