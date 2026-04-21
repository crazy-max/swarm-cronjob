package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeImage(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected string
	}{
		{
			name:     "strips digest suffix",
			image:    "busybox:latest@sha256:abcdef",
			expected: "busybox:latest",
		},
		{
			name:     "keeps tag without digest",
			image:    "busybox:latest",
			expected: "busybox:latest",
		},
		{
			name:     "keeps digest marker at start untouched",
			image:    "@sha256:abcdef",
			expected: "@sha256:abcdef",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, normalizeImage(test.image))
		})
	}
}
