package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
)

// Client represents an active docker object
type Client struct {
	Cli *client.Client
}

// NewEnvClient initializes a new Docker API client based on environment variables
func NewEnvClient() (*Client, error) {
	d, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.12"))
	if err != nil {
		return nil, err
	}

	_, err = d.ServerVersion(context.Background())
	if err != nil {
		return nil, err
	}

	return &Client{Cli: d}, err
}

func normalizeImage(image string) string {
	if i := strings.Index(image, "@sha256:"); i > 0 {
		image = image[:i]
	}
	return image
}
