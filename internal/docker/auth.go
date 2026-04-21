package docker

import (
	"context"

	"github.com/distribution/reference"
	"github.com/docker/cli/cli/command"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
)

// resolveAuthConfig is like registry.ResolveAuthConfig, but if using the
// default index, it uses the default index name for the daemon's platform,
// not the client's platform.
func resolveAuthConfig(_ context.Context, cli command.Cli, index *registrytypes.IndexInfo) registrytypes.AuthConfig {
	configKey := index.Name
	if index.Official {
		configKey = registry.IndexServer
	}
	a, _ := cli.ConfigFile().GetAuthConfig(configKey)
	return registrytypes.AuthConfig(a)
}

// retrieveAuthTokenFromImage retrieves an encoded auth token given a complete image
func retrieveAuthTokenFromImage(ctx context.Context, cli command.Cli, image string) (string, error) {
	// Retrieve encoded auth token from the image reference
	authConfig, err := resolveAuthConfigFromImage(ctx, cli, image)
	if err != nil {
		return "", err
	}
	encodedAuth, err := registrytypes.EncodeAuthConfig(authConfig)
	if err != nil {
		return "", err
	}
	return encodedAuth, nil
}

// resolveAuthConfigFromImage retrieves that AuthConfig using the image string
func resolveAuthConfigFromImage(ctx context.Context, cli command.Cli, image string) (registrytypes.AuthConfig, error) {
	registryRef, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return registrytypes.AuthConfig{}, err
	}
	repoInfo, err := registry.ParseRepositoryInfo(registryRef)
	if err != nil {
		return registrytypes.AuthConfig{}, err
	}
	return resolveAuthConfig(ctx, cli, repoInfo.Index), nil
}
