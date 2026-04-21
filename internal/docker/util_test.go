package docker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDockerTestClient(t *testing.T, routes map[string]http.HandlerFunc) (*DockerClient, *[]string) {
	t.Helper()

	requests := make([]string, 0, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := stripVersionPrefix(r.URL.Path)
		requests = append(requests, path)

		handler, ok := routes[path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)

	apiClient, err := client.New(
		client.WithHost(server.URL),
		client.WithAPIVersion("1.47"),
		client.WithHTTPClient(server.Client()),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, apiClient.Close())
	})

	return &DockerClient{api: apiClient}, &requests
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) bool {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return assert.Fail(t, "failed to encode test response", "error: %v", err)
	}
	return true
}

func stripVersionPrefix(path string) string {
	if len(path) < len("/v1.0/") || path[0] != '/' || path[1] != 'v' {
		return path
	}

	slashIndex := strings.Index(path[2:], "/")
	if slashIndex < 0 {
		return path
	}

	return path[slashIndex+2:]
}

func assertQueryHasFilterValues(t *testing.T, query url.Values, key string, expected []string) bool {
	t.Helper()

	rawFilters := query.Get("filters")
	if !assert.NotEmpty(t, rawFilters) {
		return false
	}

	var filters map[string]map[string]bool
	if err := json.Unmarshal([]byte(rawFilters), &filters); err != nil {
		assert.Fail(t, "failed to decode request filters", "error: %v", err)
		return false
	}

	actualSet, ok := filters[key]
	if !assert.True(t, ok) {
		return false
	}

	actual := make([]string, 0, len(actualSet))
	for value := range actualSet {
		actual = append(actual, value)
	}
	slices.Sort(actual)
	slices.Sort(expected)

	return assert.Equal(t, expected, actual)
}
