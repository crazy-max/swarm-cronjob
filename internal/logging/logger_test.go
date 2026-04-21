package logging

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestConfigureJSON(t *testing.T) {
	output := captureConfiguredOutput(t, &model.Cli{
		LogJSON:  true,
		LogLevel: "info",
	}, func() {
		log.Info().Msg("json log")
	})

	lines := nonEmptyLines(output)
	require.Len(t, lines, 1)

	var event map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &event))
	require.Equal(t, "info", event["level"])
	require.Equal(t, "json log", event["message"])
	require.Contains(t, event, "time")
}

func TestConfigureConsole(t *testing.T) {
	output := captureConfiguredOutput(t, &model.Cli{
		LogJSON:  false,
		LogLevel: "info",
	}, func() {
		log.Info().Msg("console log")
	})

	lines := nonEmptyLines(output)
	require.Len(t, lines, 1)
	require.Contains(t, lines[0], "console log")
	require.False(t, json.Valid([]byte(lines[0])))
}

func TestConfigureLevel(t *testing.T) {
	output := captureConfiguredOutput(t, &model.Cli{
		LogJSON:  true,
		LogLevel: "warn",
	}, func() {
		log.Info().Msg("hidden")
		log.Warn().Msg("visible")
	})

	lines := nonEmptyLines(output)
	require.Len(t, lines, 1)

	var event map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &event))
	require.Equal(t, "warn", event["level"])
	require.Equal(t, "visible", event["message"])
}

func captureConfiguredOutput(t *testing.T, cli *model.Cli, emit func()) string {
	t.Helper()

	oldStdout := os.Stdout
	oldLogger := log.Logger
	oldLevel := zerolog.GlobalLevel()

	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = oldStdout
		log.Logger = oldLogger
		zerolog.SetGlobalLevel(oldLevel)
	})

	Configure(cli)
	emit()

	require.NoError(t, writer.Close())

	output, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())

	return string(output)
}

func nonEmptyLines(output string) []string {
	lines := strings.Split(output, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			filtered = append(filtered, line)
		}
	}
	return filtered
}
