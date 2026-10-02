package model

import (
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
)

func TestCliServiceLabelFilters(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  string
		args []string
		want []string
	}{
		{name: "default"},
		{name: "repeated flags", args: []string{"--service-label", "stack=production", "--service-label", "scheduler"}, want: []string{"stack=production", "scheduler"}},
		{name: "environment", env: "stack=production,scheduler", want: []string{"stack=production", "scheduler"}},
		{name: "flags override environment", env: "stack=staging", args: []string{"--service-label", "stack=production"}, want: []string{"stack=production"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("SERVICE_LABELS", tc.env)
			cli := Cli{}
			parser, err := kong.New(&cli)
			require.NoError(t, err)
			_, err = parser.Parse(tc.args)
			require.NoError(t, err)
			require.Equal(t, tc.want, cli.ServiceLabel)
		})
	}
}
