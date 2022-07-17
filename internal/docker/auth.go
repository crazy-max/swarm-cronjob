package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/distribution/distribution/reference"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
)

// encodeAuthToBase64 serializes the auth configuration as JSON base64 payload
func encodeAuthToBase64(authConfig types.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

// resolveAuthConfig is like registry.ResolveAuthConfig, but if using the
// default index, it uses the default index name for the daemon's platform,
// not the client's platform.
func resolveAuthConfig(_ context.Context, cli command.Cli, index *registrytypes.IndexInfo) types.AuthConfig {
	configKey := index.Name
	if index.Official {
		configKey = registry.IndexServer
	}

	a, _ := cli.ConfigFile().GetAuthConfig(configKey)
	return types.AuthConfig(a)
}

// retrieveAuthTokenFromImage retrieves an encoded auth token given a complete image
func retrieveAuthTokenFromImage(ctx context.Context, cli command.Cli, image string) (string, error) {
	// Retrieve encoded auth token from the image reference
	authConfig, err := resolveAuthConfigFromImage(ctx, cli, image)
	if err != nil {
		return "", err
	}
	encodedAuth, err := encodeAuthToBase64(authConfig)
	if err != nil {
		return "", err
	}
	return encodedAuth, nil
}

// resolveAuthConfigFromImage retrieves that AuthConfig using the image string
func resolveAuthConfigFromImage(ctx context.Context, cli command.Cli, image string) (types.AuthConfig, error) {
	registryRef, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return types.AuthConfig{}, err
	}
	repoInfo, err := registry.ParseRepositoryInfo(registryRef)
	if err != nil {
		return types.AuthConfig{}, err
	}
	return resolveAuthConfig(ctx, cli, repoInfo.Index), nil
}
