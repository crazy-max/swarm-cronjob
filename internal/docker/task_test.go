package docker

import (
	"net/http"
	"testing"
	"time"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/require"
)

func TestTaskList(t *testing.T) {
	t.Parallel()

	newest := time.Date(2026, time.April, 22, 16, 0, 0, 0, time.UTC)
	older := newest.Add(-5 * time.Minute)

	nodeInspectHits := map[string]int{}

	dockerClient, _ := newDockerTestClient(t, map[string]http.HandlerFunc{
		"/tasks": func(w http.ResponseWriter, r *http.Request) {
			if !assertQueryHasFilterValues(t, r.URL.Query(), "service", []string{"backup"}) {
				return
			}

			tasks := []swarm.Task{
				{
					ID:     "older-task",
					NodeID: "node-a",
					Status: swarm.TaskStatus{
						State:   swarm.TaskStateRunning,
						Message: "older",
					},
					Meta: swarm.Meta{
						UpdatedAt: older,
					},
					Spec: swarm.TaskSpec{
						ContainerSpec: &swarm.ContainerSpec{
							Image: "busybox:1.36@sha256:older",
						},
					},
				},
				{
					ID:     "newest-task",
					NodeID: "node-a",
					Status: swarm.TaskStatus{
						State:   swarm.TaskStateRunning,
						Message: "newer",
					},
					Meta: swarm.Meta{
						UpdatedAt: newest,
					},
					Spec: swarm.TaskSpec{
						ContainerSpec: &swarm.ContainerSpec{
							Image: "busybox:1.36@sha256:newer",
						},
					},
				},
				{
					ID:     "fallback-hostname-task",
					NodeID: "node-b",
					Status: swarm.TaskStatus{
						State: swarm.TaskStateRunning,
					},
					Meta: swarm.Meta{
						UpdatedAt: newest.Add(-1 * time.Minute),
					},
					Spec: swarm.TaskSpec{
						ContainerSpec: &swarm.ContainerSpec{
							Image: "busybox:1.36@sha256:fallback",
						},
					},
				},
				{
					ID:     "missing-node-task",
					NodeID: "node-missing",
					Status: swarm.TaskStatus{
						State: swarm.TaskStateRunning,
					},
					Meta: swarm.Meta{
						UpdatedAt: newest.Add(-2 * time.Minute),
					},
					Spec: swarm.TaskSpec{
						ContainerSpec: &swarm.ContainerSpec{
							Image: "busybox:1.36@sha256:missing",
						},
					},
				},
			}
			writeJSON(t, w, tasks)
		},
		"/nodes/node-a": func(w http.ResponseWriter, r *http.Request) {
			nodeInspectHits["node-a"]++
			writeJSON(t, w, swarm.Node{
				ID: "node-a",
				Spec: swarm.NodeSpec{
					Annotations: swarm.Annotations{
						Name: "manager-a",
					},
				},
				Description: swarm.NodeDescription{
					Hostname: "host-a",
				},
			})
		},
		"/nodes/node-b": func(w http.ResponseWriter, r *http.Request) {
			nodeInspectHits["node-b"]++
			writeJSON(t, w, swarm.Node{
				ID: "node-b",
				Description: swarm.NodeDescription{
					Hostname: "worker-b",
				},
			})
		},
		"/nodes/node-missing": func(w http.ResponseWriter, r *http.Request) {
			nodeInspectHits["node-missing"]++
			http.Error(w, "missing", http.StatusNotFound)
		},
	})

	tasks, err := dockerClient.TaskList("backup")

	require.NoError(t, err)
	require.Len(t, tasks, 4)
	require.Equal(t, []string{"newest-task", "fallback-hostname-task", "missing-node-task", "older-task"}, taskIDs(tasks))
	require.Equal(t, "manager-a", tasks[0].NodeName)
	require.Equal(t, "worker-b", tasks[1].NodeName)
	require.Empty(t, tasks[2].NodeName)
	require.Equal(t, "manager-a", tasks[3].NodeName)
	require.Equal(t, "backup", tasks[0].ServiceName)
	require.Equal(t, "busybox:1.36", tasks[0].Image)
	require.Equal(t, map[string]int{
		"node-a":       1,
		"node-b":       1,
		"node-missing": 1,
	}, nodeInspectHits)
}

func TestTaskListEmpty(t *testing.T) {
	t.Parallel()

	dockerClient, _ := newDockerTestClient(t, map[string]http.HandlerFunc{
		"/tasks": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, []swarm.Task{})
		},
	})

	tasks, err := dockerClient.TaskList("backup")

	require.NoError(t, err)
	require.Nil(t, tasks)
}

func taskIDs(tasks []*model.TaskInfo) []string {
	ids := make([]string, 0, len(tasks))
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}
	return ids
}
