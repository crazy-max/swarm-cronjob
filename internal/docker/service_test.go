package docker

import (
	"net/http"
	"testing"
	"time"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/moby/moby/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceList(t *testing.T) {
	t.Parallel()

	replicas := uint64(3)
	updatedAt := time.Date(2026, time.April, 22, 14, 30, 0, 0, time.UTC)

	dockerClient, apiRequests := newDockerTestClient(t, map[string]http.HandlerFunc{
		"/services": func(w http.ResponseWriter, r *http.Request) {
			if !assert.Equal(t, http.MethodGet, r.Method) {
				return
			}
			if !assertQueryHasFilterValues(t, r.URL.Query(), "name", []string{"db"}) {
				return
			}
			if !assertQueryHasFilterValues(t, r.URL.Query(), "label", []string{"swarm.cronjob.enable", "swarm.cronjob.schedule"}) {
				return
			}

			services := []swarm.Service{
				{
					ID: "service-global",
					Spec: swarm.ServiceSpec{
						Annotations: swarm.Annotations{
							Name:   "zzz-global",
							Labels: map[string]string{"scope": "global"},
						},
						TaskTemplate: swarm.TaskSpec{
							ContainerSpec: &swarm.ContainerSpec{
								Image: "busybox:1.36@sha256:global",
							},
						},
						Mode: swarm.ServiceMode{
							Global: &swarm.GlobalService{},
						},
					},
					Meta: swarm.Meta{
						UpdatedAt: updatedAt.Add(2 * time.Minute),
					},
				},
				{
					ID: "service-replicated",
					Spec: swarm.ServiceSpec{
						Annotations: swarm.Annotations{
							Name:   "aaa-replicated",
							Labels: map[string]string{"scope": "replicated"},
						},
						TaskTemplate: swarm.TaskSpec{
							ContainerSpec: &swarm.ContainerSpec{
								Image: "nginx:1.27@sha256:replicated",
							},
						},
						Mode: swarm.ServiceMode{
							Replicated: &swarm.ReplicatedService{
								Replicas: &replicas,
							},
						},
					},
					PreviousSpec: &swarm.ServiceSpec{},
					UpdateStatus: &swarm.UpdateStatus{
						State: swarm.UpdateStateCompleted,
					},
					Meta: swarm.Meta{
						UpdatedAt: updatedAt,
					},
				},
			}
			writeJSON(t, w, services)
		},
		"/nodes": func(w http.ResponseWriter, r *http.Request) {
			if !assert.Equal(t, http.MethodGet, r.Method) {
				return
			}

			nodes := []swarm.Node{
				{
					ID: "node-active",
					Status: swarm.NodeStatus{
						State: swarm.NodeStateReady,
					},
				},
				{
					ID: "node-down",
					Status: swarm.NodeStatus{
						State: swarm.NodeStateDown,
					},
				},
			}
			writeJSON(t, w, nodes)
		},
		"/tasks": func(w http.ResponseWriter, r *http.Request) {
			if !assert.Equal(t, http.MethodGet, r.Method) {
				return
			}
			if !assertQueryHasFilterValues(t, r.URL.Query(), "service", []string{"service-global", "service-replicated"}) {
				return
			}

			tasks := []swarm.Task{
				{
					ID:           "task-running-active",
					ServiceID:    "service-replicated",
					NodeID:       "node-active",
					DesiredState: swarm.TaskStateRunning,
					Status: swarm.TaskStatus{
						State: swarm.TaskStateRunning,
					},
				},
				{
					ID:           "task-running-down-node",
					ServiceID:    "service-replicated",
					NodeID:       "node-down",
					DesiredState: swarm.TaskStateRunning,
					Status: swarm.TaskStatus{
						State: swarm.TaskStateRunning,
					},
				},
				{
					ID:           "task-shutdown-global",
					ServiceID:    "service-global",
					NodeID:       "node-active",
					DesiredState: swarm.TaskStateShutdown,
					Status: swarm.TaskStatus{
						State: swarm.TaskStateShutdown,
					},
				},
				{
					ID:           "task-running-global",
					ServiceID:    "service-global",
					NodeID:       "node-active",
					DesiredState: swarm.TaskStateRunning,
					Status: swarm.TaskStatus{
						State: swarm.TaskStateRunning,
					},
				},
			}
			writeJSON(t, w, tasks)
		},
	})

	services, err := dockerClient.ServiceList(&model.ServiceListArgs{
		Name: "db",
		Labels: []string{
			"swarm.cronjob.enable",
			"swarm.cronjob.schedule",
		},
	})

	require.NoError(t, err)
	require.Len(t, services, 2)
	require.Equal(t, []string{"/services", "/nodes", "/tasks"}, *apiRequests)

	replicatedService := services[0]
	require.Equal(t, "aaa-replicated", replicatedService.Name)
	require.Equal(t, "nginx:1.27", replicatedService.Image)
	require.Equal(t, model.ServiceModeReplicated, replicatedService.Mode)
	require.EqualValues(t, 3, replicatedService.Replicas)
	require.EqualValues(t, 1, replicatedService.Actives)
	require.EqualValues(t, 2, replicatedService.Busy)
	require.True(t, replicatedService.Rollback)
	require.Equal(t, string(swarm.UpdateStateCompleted), replicatedService.UpdateStatus)
	require.Equal(t, updatedAt.Local(), replicatedService.UpdatedAt)

	globalService := services[1]
	require.Equal(t, "zzz-global", globalService.Name)
	require.Equal(t, "busybox:1.36", globalService.Image)
	require.Equal(t, model.ServiceModeGlobal, globalService.Mode)
	require.EqualValues(t, 1, globalService.Replicas)
	require.EqualValues(t, 1, globalService.Actives)
	require.EqualValues(t, 1, globalService.Busy)
	require.False(t, globalService.Rollback)
	require.Empty(t, globalService.UpdateStatus)
}

func TestServiceNotFound(t *testing.T) {
	t.Parallel()

	dockerClient, _ := newDockerTestClient(t, map[string]http.HandlerFunc{
		"/services": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, []swarm.Service{})
		},
		"/nodes": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, []swarm.Node{})
		},
		"/tasks": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, []swarm.Task{})
		},
	})

	service, err := dockerClient.Service("missing")

	require.Nil(t, service)
	require.EqualError(t, err, "missing service not found")
}
