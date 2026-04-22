package docker

import (
	"context"
	"strings"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/docker/cli/cli/command"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/api/types/registry"
	"github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
)

// Client for Swarm
type Client interface {
	DistributionInspect(ctx context.Context, image, encodedAuth string) (registry.DistributionInspect, error)
	RetrieveAuthTokenFromImage(ctx context.Context, image string) (string, error)
	ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options client.ServiceUpdateOptions) (client.ServiceUpdateResult, error)
	ServiceInspectWithRaw(ctx context.Context, serviceID string, opts client.ServiceInspectOptions) (swarm.Service, []byte, error)
	Events(ctx context.Context, options client.EventsListOptions) (<-chan events.Message, <-chan error)

	ServiceList(args *model.ServiceListArgs) ([]*model.ServiceInfo, error)
	Service(name string) (*model.ServiceInfo, error)
	TaskList(service string) ([]*model.TaskInfo, error)
}

type DockerClient struct {
	api *client.Client
	cli command.Cli
}

// NewEnvClient initializes a new Docker API client based on environment
// variables
func NewEnvClient() (*DockerClient, error) {
	dockerAPICli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Docker API client")
	}

	_, err = dockerAPICli.ServerVersion(context.Background(), client.ServerVersionOptions{})
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
	}, nil
}

// DistributionInspect returns the image digest with full Manifest
func (c *DockerClient) DistributionInspect(ctx context.Context, image, encodedAuth string) (registry.DistributionInspect, error) {
	result, err := c.api.DistributionInspect(ctx, image, client.DistributionInspectOptions{
		EncodedRegistryAuth: encodedAuth,
	})
	if err != nil {
		return registry.DistributionInspect{}, err
	}
	return result.DistributionInspect, nil
}

// RetrieveAuthTokenFromImage retrieves an encoded auth token given a complete
// image.
func (c *DockerClient) RetrieveAuthTokenFromImage(_ context.Context, image string) (string, error) {
	return command.RetrieveAuthTokenFromImage(c.cli.ConfigFile(), image)
}

// ServiceUpdate updates a Service. The version number is required to avoid
// conflicting writes. It should be the value as set *before* the update. You
// can find this value in the Meta field of swarm.Service, which can be found
// using ServiceInspectWithRaw.
func (c *DockerClient) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, service swarm.ServiceSpec, options client.ServiceUpdateOptions) (client.ServiceUpdateResult, error) {
	options.Version = version
	options.Spec = service
	return c.api.ServiceUpdate(ctx, serviceID, options)
}

// ServiceInspectWithRaw returns the service information and the raw data.
func (c *DockerClient) ServiceInspectWithRaw(ctx context.Context, serviceID string, opts client.ServiceInspectOptions) (swarm.Service, []byte, error) {
	result, err := c.api.ServiceInspect(ctx, serviceID, opts)
	if err != nil {
		return swarm.Service{}, nil, err
	}
	return result.Service, result.Raw, nil
}

// Events returns a stream of events in the daemon. It's up to the caller to
// close the stream by cancelling the context. Once the stream has been
// completely read an io.EOF error will be sent over the error channel. If an
// error is sent all processing will be stopped. It's up to the caller to
// reopen the stream in the event of an error by reinvoking this method.
func (c *DockerClient) Events(ctx context.Context, options client.EventsListOptions) (<-chan events.Message, <-chan error) {
	result := c.api.Events(ctx, options)
	return result.Messages, result.Err
}

func normalizeImage(image string) string {
	if i := strings.Index(image, "@sha256:"); i > 0 {
		image = image[:i]
	}
	return image
}
