package docker

import (
	"context"
	"strings"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// Client for Swarm
type Client interface {
	DistributionInspect(ctx context.Context, image, encodedAuth string) (registry.DistributionInspect, error)
	RetrieveAuthTokenFromImage(ctx context.Context, image string) (string, error)
	ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error)
	ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error)
	Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error)

	ServiceList(args *model.ServiceListArgs) ([]*model.ServiceInfo, error)
	Service(name string) (*model.ServiceInfo, error)
	TaskList(service string) ([]*model.TaskInfo, error)
}

type DockerClient struct {
	api *client.Client
	cli command.Cli
}

// NewEnvClient initializes a new Docker API client based on environment variables
func NewEnvClient() (*DockerClient, error) {
	dockerAPICli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.24"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Docker API client")
	}

	_, err = dockerAPICli.ServerVersion(context.Background())
	if err != nil {
		return nil, err
	}

	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Docker cli")
	}

	return &DockerClient{
		api: dockerAPICli,
		cli: dockerCli,
	}, err
}

// DistributionInspect returns the image digest with full Manifest
func (c *DockerClient) DistributionInspect(ctx context.Context, image, encodedAuth string) (registry.DistributionInspect, error) {
	return c.api.DistributionInspect(ctx, image, encodedAuth)
}

// RetrieveAuthTokenFromImage retrieves an encoded auth token given a complete image
func (c *DockerClient) RetrieveAuthTokenFromImage(ctx context.Context, image string) (string, error) {
	return retrieveAuthTokenFromImage(ctx, c.cli, image)
}

// ServiceUpdate updates a Service. The version number is required to avoid conflicting writes.
// It should be the value as set *before* the update. You can find this value in the Meta field
// of swarm.Service, which can be found using ServiceInspectWithRaw.
func (c *DockerClient) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options types.ServiceUpdateOptions) (types.ServiceUpdateResponse, error) {
	return c.api.ServiceUpdate(ctx, serviceID, version, service, options)
}

// ServiceInspectWithRaw returns the service information and the raw data.
func (c *DockerClient) ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error) {
	return c.api.ServiceInspectWithRaw(ctx, serviceID, opts)
}

// Events returns a stream of events in the daemon. It's up to the caller to close the stream
// by cancelling the context. Once the stream has been completely read an io.EOF error will
// be sent over the error channel. If an error is sent all processing will be stopped. It's up
// to the caller to reopen the stream in the event of an error by reinvoking this method.
func (c *DockerClient) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	return c.api.Events(ctx, options)
}

func normalizeImage(image string) string {
	if i := strings.Index(image, "@sha256:"); i > 0 {
		image = image[:i]
	}
	return image
}
