package worker

import (
	"context"
	"errors"
	"testing"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/registry"
	"github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/require"
)

type workerDockerStub struct {
	serviceResponses []*model.ServiceInfo
	taskList         []*model.TaskInfo
	authToken        string
	authErr          error
	updateErr        error
	updateCalls      []serviceUpdateCall
	serviceCalls     int
	authCalls        int
}

type serviceUpdateCall struct {
	serviceID string
	version   swarm.Version
	replicas  *uint64
	labels    map[string]string
	force     uint64
	options   client.ServiceUpdateOptions
}

func (s *workerDockerStub) DistributionInspect(context.Context, string, string) (registry.DistributionInspect, error) {
	return registry.DistributionInspect{}, nil
}

func (s *workerDockerStub) RetrieveAuthTokenFromImage(context.Context, string) (string, error) {
	s.authCalls++
	return s.authToken, s.authErr
}

func (s *workerDockerStub) ServiceUpdate(_ context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options client.ServiceUpdateOptions) (client.ServiceUpdateResult, error) {
	call := serviceUpdateCall{
		serviceID: serviceID,
		version:   version,
		labels:    cloneLabels(service.Labels),
		force:     service.TaskTemplate.ForceUpdate,
		options:   options,
	}
	if service.Mode.Replicated != nil && service.Mode.Replicated.Replicas != nil {
		call.replicas = new(uint64)
		*call.replicas = *service.Mode.Replicated.Replicas
	}
	s.updateCalls = append(s.updateCalls, call)
	return client.ServiceUpdateResult{}, s.updateErr
}

func (s *workerDockerStub) ServiceInspectWithRaw(context.Context, string, client.ServiceInspectOptions) (swarm.Service, []byte, error) {
	return swarm.Service{}, nil, nil
}

func (s *workerDockerStub) Events(context.Context, client.EventsListOptions) (<-chan events.Message, <-chan error) {
	return nil, nil
}

func (s *workerDockerStub) Service(string) (*model.ServiceInfo, error) {
	if s.serviceCalls >= len(s.serviceResponses) {
		return nil, errors.New("unexpected service call")
	}
	service := s.serviceResponses[s.serviceCalls]
	s.serviceCalls++
	return service, nil
}

func (s *workerDockerStub) ServiceList(*model.ServiceListArgs) ([]*model.ServiceInfo, error) {
	return nil, nil
}

func (s *workerDockerStub) TaskList(string) ([]*model.TaskInfo, error) {
	return s.taskList, nil
}

func TestRunSkipsUpdateWhenJobIsAlreadyActive(t *testing.T) {
	stub := &workerDockerStub{
		serviceResponses: []*model.ServiceInfo{
			{
				Name:    "backup",
				Actives: 2,
				Raw:     replicatedService("svc-1", "backup", 7, 1, "busybox:latest"),
			},
		},
	}

	client := Client{
		Docker: stub,
		Job: model.Job{
			Name:        "backup",
			SkipRunning: true,
			Replicas:    1,
		},
	}

	client.Run()

	require.Empty(t, stub.updateCalls)
	require.Zero(t, stub.authCalls)
}

func TestRunScalesDownReplicatedServiceBeforeUpdating(t *testing.T) {
	stub := &workerDockerStub{
		serviceResponses: []*model.ServiceInfo{
			{
				Name: "backup",
				Mode: model.ServiceModeReplicated,
				Raw:  replicatedService("svc-1", "backup", 7, 3, "busybox:latest"),
			},
			{
				Name: "backup",
				Mode: model.ServiceModeReplicated,
				Raw:  replicatedService("svc-1", "backup", 8, 0, "busybox:latest"),
			},
		},
	}

	client := Client{
		Docker: stub,
		Job: model.Job{
			Name:     "backup",
			Replicas: 2,
		},
	}

	client.Run()

	require.Len(t, stub.updateCalls, 2)

	scaleDown := stub.updateCalls[0]
	require.NotNil(t, scaleDown.replicas)
	require.Zero(t, *scaleDown.replicas)
	require.Equal(t, "true", scaleDown.labels["swarm.cronjob.scaledown"])
	require.EqualValues(t, 7, scaleDown.force)

	finalUpdate := stub.updateCalls[1]
	require.NotNil(t, finalUpdate.replicas)
	require.EqualValues(t, 2, *finalUpdate.replicas)
	require.NotContains(t, finalUpdate.labels, "swarm.cronjob.scaledown")
	require.EqualValues(t, 8, finalUpdate.force)
	require.Equal(t, swarm.RegistryAuthFromSpec, finalUpdate.options.RegistryAuthFrom)
}

func TestRunUsesRegistryAuthAndQueryRegistryFlags(t *testing.T) {
	queryRegistry := new(bool)
	*queryRegistry = true
	stub := &workerDockerStub{
		serviceResponses: []*model.ServiceInfo{
			{
				Name: "backup",
				Mode: model.ServiceModeReplicated,
				Raw:  replicatedService("svc-1", "backup", 9, 1, "busybox:latest"),
			},
		},
		authToken: "encoded-auth",
	}

	client := Client{
		Docker: stub,
		Job: model.Job{
			Name:          "backup",
			RegistryAuth:  true,
			QueryRegistry: queryRegistry,
			Replicas:      1,
		},
	}

	client.Run()

	require.Equal(t, 1, stub.authCalls)
	require.Len(t, stub.updateCalls, 1)

	update := stub.updateCalls[0]
	require.Equal(t, "encoded-auth", update.options.EncodedRegistryAuth)
	require.Empty(t, update.options.RegistryAuthFrom)
	require.True(t, update.options.QueryRegistry)
	require.EqualValues(t, 9, update.force)
}

func TestScaleDownRemovesScaledownLabelFromReturnedService(t *testing.T) {
	stub := &workerDockerStub{
		serviceResponses: []*model.ServiceInfo{
			{
				Name: "backup",
				Raw:  replicatedService("svc-1", "backup", 11, 0, "busybox:latest"),
			},
		},
	}

	client := Client{
		Docker: stub,
		Job:    model.Job{Name: "backup"},
	}

	service := replicatedService("svc-1", "backup", 10, 3, "busybox:latest")
	res, err := client.scaleDown(service)
	require.NoError(t, err)

	require.Len(t, stub.updateCalls, 1)
	update := stub.updateCalls[0]
	require.NotNil(t, update.replicas)
	require.Zero(t, *update.replicas)
	require.Equal(t, "true", update.labels["swarm.cronjob.scaledown"])
	require.NotContains(t, res.Spec.Labels, "swarm.cronjob.scaledown")
}

func replicatedService(id, name string, versionIndex, replicas uint64, image string) swarm.Service {
	replicasCopy := new(uint64)
	*replicasCopy = replicas
	return swarm.Service{
		ID: id,
		Meta: swarm.Meta{
			Version: swarm.Version{Index: versionIndex},
		},
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name:   name,
				Labels: map[string]string{},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: image,
				},
			},
			Mode: swarm.ServiceMode{
				Replicated: &swarm.ReplicatedService{
					Replicas: replicasCopy,
				},
			},
		},
	}
}

func cloneLabels(labels map[string]string) map[string]string {
	if labels == nil {
		return nil
	}
	cloned := make(map[string]string, len(labels))
	for key, value := range labels {
		cloned[key] = value
	}
	return cloned
}
