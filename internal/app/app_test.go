package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/crazy-max/swarm-cronjob/internal/worker"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/registry"
	"github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"
	cron "github.com/robfig/cron/v3"
	"github.com/stretchr/testify/require"
)

type appDockerStub struct {
	service        *model.ServiceInfo
	serviceErr     error
	serviceHits    int
	serviceList    []*model.ServiceInfo
	serviceListErr error
	eventMsgs      <-chan events.Message
	eventErrs      <-chan error
}

func (s *appDockerStub) DistributionInspect(context.Context, string, string) (registry.DistributionInspect, error) {
	return registry.DistributionInspect{}, nil
}

func (s *appDockerStub) RetrieveAuthTokenFromImage(context.Context, string) (string, error) {
	return "", nil
}

func (s *appDockerStub) ServiceUpdate(context.Context, string, swarm.Version, swarm.ServiceSpec, client.ServiceUpdateOptions) (client.ServiceUpdateResult, error) {
	return client.ServiceUpdateResult{}, nil
}

func (s *appDockerStub) ServiceInspectWithRaw(context.Context, string, client.ServiceInspectOptions) (swarm.Service, []byte, error) {
	return swarm.Service{}, nil, nil
}

func (s *appDockerStub) Events(context.Context, client.EventsListOptions) (<-chan events.Message, <-chan error) {
	return s.eventMsgs, s.eventErrs
}

func (s *appDockerStub) Service(string) (*model.ServiceInfo, error) {
	s.serviceHits++
	if s.serviceErr != nil {
		return nil, s.serviceErr
	}
	return s.service, nil
}

func (s *appDockerStub) ServiceList(*model.ServiceListArgs) ([]*model.ServiceInfo, error) {
	if s.serviceListErr != nil {
		return nil, s.serviceListErr
	}
	return s.serviceList, nil
}

func (s *appDockerStub) TaskList(string) ([]*model.TaskInfo, error) {
	return nil, nil
}

func TestCrudJobAddsCronEntryWithParsedWorkerSettings(t *testing.T) {
	expectedQueryRegistry := new(bool)

	dockerStub := &appDockerStub{
		service: &model.ServiceInfo{
			Name: "backup",
			Labels: map[string]string{
				"swarm.cronjob.enable":         "true",
				"swarm.cronjob.schedule":       "0 * * * *",
				"swarm.cronjob.skip-running":   "true",
				"swarm.cronjob.replicas":       "3",
				"swarm.cronjob.registry-auth":  "true",
				"swarm.cronjob.query-registry": "false",
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	processed, err := sc.crudJob("backup")

	require.NoError(t, err)
	require.True(t, processed)
	require.Len(t, sc.jobs, 1)

	entryID := sc.jobs["backup"]
	entry := sc.cron.Entry(entryID)
	require.True(t, entry.Valid())

	workerClient, ok := entry.Job.(*worker.Client)
	require.True(t, ok)
	require.Same(t, dockerStub, workerClient.Docker)
	require.Equal(t, model.Job{
		Name:          "backup",
		Enable:        true,
		Schedule:      "0 * * * *",
		SkipRunning:   true,
		RegistryAuth:  true,
		QueryRegistry: expectedQueryRegistry,
		Replicas:      3,
	}, workerClient.Job)
}

func TestCrudJobDisablesExistingEntryWhenServiceIsDisabled(t *testing.T) {
	dockerStub := &appDockerStub{
		service: &model.ServiceInfo{
			Name: "backup",
			Labels: map[string]string{
				"swarm.cronjob.enable":   "true",
				"swarm.cronjob.schedule": "0 * * * *",
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	processed, err := sc.crudJob("backup")
	require.NoError(t, err)
	require.True(t, processed)
	require.Len(t, sc.jobs, 1)

	existingID := sc.jobs["backup"]
	require.True(t, sc.cron.Entry(existingID).Valid())

	dockerStub.service = &model.ServiceInfo{
		Name: "backup",
		Labels: map[string]string{
			"swarm.cronjob.enable":   "false",
			"swarm.cronjob.schedule": "0 * * * *",
		},
	}

	processed, err = sc.crudJob("backup")

	require.NoError(t, err)
	require.True(t, processed)
	require.Empty(t, sc.jobs)
	require.False(t, sc.cron.Entry(existingID).Valid())
}

func TestCrudJobRemovesExistingEntryWhenServiceNoLongerExists(t *testing.T) {
	dockerStub := &appDockerStub{
		service: &model.ServiceInfo{
			Name: "backup",
			Labels: map[string]string{
				"swarm.cronjob.enable":   "true",
				"swarm.cronjob.schedule": "0 * * * *",
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	processed, err := sc.crudJob("backup")
	require.NoError(t, err)
	require.True(t, processed)

	existingID := sc.jobs["backup"]
	dockerStub.service = nil
	dockerStub.serviceErr = errors.New("service not found")

	processed, err = sc.crudJob("backup")

	require.NoError(t, err)
	require.True(t, processed)
	require.Empty(t, sc.jobs)
	require.False(t, sc.cron.Entry(existingID).Valid())
}

func TestCrudJobSkipsScaledownServicesWithoutReplacingCurrentEntry(t *testing.T) {
	dockerStub := &appDockerStub{
		service: &model.ServiceInfo{
			Name: "backup",
			Labels: map[string]string{
				"swarm.cronjob.enable":   "true",
				"swarm.cronjob.schedule": "0 * * * *",
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	processed, err := sc.crudJob("backup")
	require.NoError(t, err)
	require.True(t, processed)

	existingID := sc.jobs["backup"]
	existingEntry := sc.cron.Entry(existingID)
	require.True(t, existingEntry.Valid())

	dockerStub.service = &model.ServiceInfo{
		Name: "backup",
		Labels: map[string]string{
			"swarm.cronjob.enable":    "true",
			"swarm.cronjob.schedule":  "*/5 * * * *",
			"swarm.cronjob.scaledown": "true",
		},
	}

	processed, err = sc.crudJob("backup")

	require.NoError(t, err)
	require.False(t, processed)
	require.Len(t, sc.jobs, 1)
	require.Equal(t, existingID, sc.jobs["backup"])
	require.True(t, sc.cron.Entry(existingID).Valid())
	require.Equal(t, existingEntry.ID, sc.cron.Entry(existingID).ID)
}

func TestCrudJobTreatsUnchangedExistingEntryAsProcessed(t *testing.T) {
	dockerStub := &appDockerStub{
		service: &model.ServiceInfo{
			Name: "backup",
			Labels: map[string]string{
				"swarm.cronjob.enable":   "true",
				"swarm.cronjob.schedule": "0 * * * *",
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	processed, err := sc.crudJob("backup")
	require.NoError(t, err)
	require.True(t, processed)

	existingID := sc.jobs["backup"]

	processed, err = sc.crudJob("backup")
	require.NoError(t, err)
	require.True(t, processed)
	require.Equal(t, existingID, sc.jobs["backup"])
	require.True(t, sc.cron.Entry(existingID).Valid())
}

func TestReconcileJobsUpdatesExistingEntryWhenScheduleChanges(t *testing.T) {
	dockerStub := &appDockerStub{
		serviceList: []*model.ServiceInfo{
			{
				Name: "backup",
				Labels: map[string]string{
					"swarm.cronjob.enable":   "true",
					"swarm.cronjob.schedule": "0 * * * *",
				},
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	require.NoError(t, sc.reconcileJobs())

	existingID := sc.jobs["backup"]
	require.True(t, sc.cron.Entry(existingID).Valid())

	dockerStub.serviceList = []*model.ServiceInfo{
		{
			Name: "backup",
			Labels: map[string]string{
				"swarm.cronjob.enable":   "true",
				"swarm.cronjob.schedule": "*/5 * * * *",
			},
		},
	}

	require.NoError(t, sc.reconcileJobs())

	updatedID := sc.jobs["backup"]
	require.NotEqual(t, existingID, updatedID)
	require.False(t, sc.cron.Entry(existingID).Valid())

	entry := sc.cron.Entry(updatedID)
	require.True(t, entry.Valid())

	workerClient, ok := entry.Job.(*worker.Client)
	require.True(t, ok)
	require.Equal(t, "*/5 * * * *", workerClient.Job.Schedule)
}

func TestReconcileJobsKeepsExistingEntryWhenJobIsUnchanged(t *testing.T) {
	dockerStub := &appDockerStub{
		serviceList: []*model.ServiceInfo{
			{
				Name: "backup",
				Labels: map[string]string{
					"swarm.cronjob.enable":         "true",
					"swarm.cronjob.schedule":       "0 * * * *",
					"swarm.cronjob.skip-running":   "true",
					"swarm.cronjob.replicas":       "3",
					"swarm.cronjob.registry-auth":  "true",
					"swarm.cronjob.query-registry": "false",
				},
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	require.NoError(t, sc.reconcileJobs())

	existingID := sc.jobs["backup"]
	entry := sc.cron.Entry(existingID)
	require.True(t, entry.Valid())

	require.NoError(t, sc.reconcileJobs())

	require.Equal(t, existingID, sc.jobs["backup"])
	require.Equal(t, entry.ID, sc.cron.Entry(existingID).ID)
}

func TestReconcileJobsRemovesExistingEntryWhenServiceLosesCronLabels(t *testing.T) {
	dockerStub := &appDockerStub{
		service: &model.ServiceInfo{
			Name:   "backup",
			Labels: map[string]string{},
		},
		serviceList: []*model.ServiceInfo{
			{
				Name: "backup",
				Labels: map[string]string{
					"swarm.cronjob.enable":   "true",
					"swarm.cronjob.schedule": "0 * * * *",
				},
			},
		},
	}

	sc := newTestSwarmCronjob(dockerStub)

	require.NoError(t, sc.reconcileJobs())

	existingID := sc.jobs["backup"]
	require.True(t, sc.cron.Entry(existingID).Valid())

	dockerStub.serviceList = nil

	require.NoError(t, sc.reconcileJobs())

	require.Empty(t, sc.jobs)
	require.False(t, sc.cron.Entry(existingID).Valid())
}

func TestRunReturnsWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	sc := newTestSwarmCronjob(&appDockerStub{
		serviceList: []*model.ServiceInfo{},
		eventMsgs:   make(chan events.Message),
		eventErrs:   make(chan error),
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- sc.Run(ctx)
	}()

	select {
	case err := <-errCh:
		t.Fatalf("Run returned before cancellation: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	cancel(nil)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Run to return after cancellation")
	}
}

func TestRunReturnsEventChannelError(t *testing.T) {
	eventErrs := make(chan error, 1)
	eventErrs <- errors.New("boom")

	sc := newTestSwarmCronjob(&appDockerStub{
		serviceList: []*model.ServiceInfo{},
		eventMsgs:   make(chan events.Message),
		eventErrs:   eventErrs,
	})

	err := sc.Run(context.Background())
	require.EqualError(t, err, "event channel failed: boom")
}

func newTestSwarmCronjob(dockerClient *appDockerStub) *SwarmCronjob {
	return &SwarmCronjob{
		docker: dockerClient,
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		)),
		jobs: make(map[string]cron.EntryID),
	}
}
