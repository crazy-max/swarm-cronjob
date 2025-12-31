package docker

import (
	"context"
	"sort"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

// TaskList return all running tasks of a service.
func (c *DockerClient) TaskList(service string) ([]*model.TaskInfo, error) {
	tasksFilters := filters.NewArgs()
	tasksFilters.Add("service", service)
	tasks, err := c.api.TaskList(context.Background(), swarm.TaskListOptions{
		Filters: tasksFilters,
	})
	if err != nil || len(tasks) == 0 {
		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
	})
	nodes := make(map[string]string)
	for _, t := range tasks {
		if _, ok := nodes[t.NodeID]; !ok {
			if node, _, e := c.api.NodeInspectWithRaw(context.Background(), t.NodeID); e == nil {
				if node.Spec.Name == "" {
					nodes[t.NodeID] = node.Description.Hostname
				} else {
					nodes[t.NodeID] = node.Spec.Name
				}
			} else {
				nodes[t.NodeID] = ""
			}
		}
	}

	res := make([]*model.TaskInfo, len(tasks))
	for i, t := range tasks {
		res[i] = &model.TaskInfo{
			Task:        t,
			NodeName:    nodes[t.NodeID],
			ServiceName: service,
			Image:       normalizeImage(t.Spec.ContainerSpec.Image),
		}
	}

	return res, nil
}
