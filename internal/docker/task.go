package docker

import (
	"context"
	"sort"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/moby/moby/client"
)

// TaskList return all running tasks of a service.
func (c *DockerClient) TaskList(service string) ([]*model.TaskInfo, error) {
	tasksFilters := make(client.Filters).Add("service", service)
	tasksRes, err := c.api.TaskList(context.Background(), client.TaskListOptions{
		Filters: tasksFilters,
	})
	if err != nil {
		return nil, err
	}
	tasks := tasksRes.Items
	if len(tasks) == 0 {
		return nil, nil
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
	})

	nodes := make(map[string]string)
	for _, t := range tasks {
		if _, ok := nodes[t.NodeID]; !ok {
			if result, e := c.api.NodeInspect(context.Background(), t.NodeID, client.NodeInspectOptions{}); e == nil {
				if result.Node.Spec.Name == "" {
					nodes[t.NodeID] = result.Node.Description.Hostname
				} else {
					nodes[t.NodeID] = result.Node.Spec.Name
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
