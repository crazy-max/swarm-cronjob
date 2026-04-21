package docker

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	configtypes "github.com/docker/cli/cli/config/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	dockerregistry "github.com/docker/docker/registry"
	"github.com/stretchr/testify/require"
)

func TestAuthConfigFromOfficialImage(t *testing.T) {
	cli := newDockerCLIWithAuthConfigs(t, map[string]configtypes.AuthConfig{
		dockerregistry.IndexServer: {
			Username: "hub-user",
			Password: "hub-password",
		},
	})

	authConfig, err := resolveAuthConfigFromImage(context.Background(), cli, "alpine:latest")

	require.NoError(t, err)
	require.Equal(t, "hub-user", authConfig.Username)
	require.Equal(t, "hub-password", authConfig.Password)
	require.Equal(t, dockerregistry.IndexServer, authConfig.ServerAddress)
}

func TestAuthConfigFromRegistryImage(t *testing.T) {
	cli := newDockerCLIWithAuthConfigs(t, map[string]configtypes.AuthConfig{
		"ghcr.io": {
			Username: "octo-user",
			Password: "octo-password",
		},
	})

	authConfig, err := resolveAuthConfigFromImage(context.Background(), cli, "ghcr.io/acme/app:1.2.3")

	require.NoError(t, err)
	require.Equal(t, "octo-user", authConfig.Username)
	require.Equal(t, "octo-password", authConfig.Password)
	require.Equal(t, "ghcr.io", authConfig.ServerAddress)
}

func TestAuthTokenFromImage(t *testing.T) {
	cli := newDockerCLIWithAuthConfigs(t, map[string]configtypes.AuthConfig{
		"ghcr.io": {
			Username: "octo-user",
			Password: "octo-password",
		},
	})

	encodedAuth, err := retrieveAuthTokenFromImage(context.Background(), cli, "ghcr.io/acme/app:1.2.3")

	require.NoError(t, err)

	authConfig, err := registrytypes.DecodeAuthConfig(encodedAuth)
	require.NoError(t, err)
	require.Equal(t, "octo-user", authConfig.Username)
	require.Equal(t, "octo-password", authConfig.Password)
	require.Equal(t, "ghcr.io", authConfig.ServerAddress)
}

func newDockerCLIWithAuthConfigs(t *testing.T, authConfigs map[string]configtypes.AuthConfig) command.Cli {
	t.Helper()

	configDir := t.TempDir()
	previousConfigDir := config.Dir()
	config.SetDir(configDir)
	t.Cleanup(func() {
		config.SetDir(previousConfigDir)
	})

	dockerConfig := configfile.New(filepath.Join(configDir, config.ConfigFileName))
	dockerConfig.AuthConfigs = authConfigs
	require.NoError(t, dockerConfig.Save())

	cli, err := command.NewDockerCli()
	require.NoError(t, err)

	return cli
}
