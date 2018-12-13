package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crazy-max/cron"
	. "github.com/crazy-max/swarm-cronjob/app"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: AppName,
	Run: cronRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Logger.Error().Err(err).Msg("rootCmd failed")
		os.Exit(1)
	}
}

func cronRun(cmd *cobra.Command, args []string) {
	Logger.Info().Msgf("Starting %s v%s", AppName, AppVersion)

	dcli, err := DockerEnvClient()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Cannot create Docker client")
	}

	services, err := ScheduledServices(dcli)
	if err != nil {
		Logger.Error().Err(err).Msg("Cannot retrieve scheduled services")
	}

	loc, err := time.LoadLocation(GetEnv("TZ", "UTC"))
	if err != nil {
		Logger.Fatal().Err(err).Msgf("Failed to load time zone %s", GetEnv("TZ", "UTC"))
	}
	c := cron.NewWithLocation(loc)
	for _, service := range services {
		if err := UpdateJob(service, dcli, c); err != nil {
			Logger.Error().Err(err).Msgf("Cannot update job for service %s", service.Spec.Name)
		}
	}
	c.Start()

	// Handle os signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		c.Stop()
		os.Exit(1)
	}()

	// Listen Docker events
	filter := filters.NewArgs()
	filter.Add("type", "service")

	msgs, errs := dcli.Events(context.Background(), types.EventsOptions{
		Filters: filter,
	})

	var event ServiceEvent
	for {
		select {
		case err := <-errs:
			Logger.Fatal().Err(err).Msg("Event channel failed")
		case msg := <-msgs:
			err := mapstructure.Decode(msg.Actor.Attributes, &event)
			if err != nil {
				Logger.Warn().Msgf("Cannot decode event, %v", err)
				continue
			}
			Logger.Debug().Msgf("Event triggered for %s (newstate='%s' oldstate='%s')", event.Service, event.UpdateState.New, event.UpdateState.Old)
			service, err := Service(dcli, event.Service)
			if err != nil {
				Logger.Error().Err(err).Msgf("Cannot find service %s", event.Service)
				continue
			}
			if err := UpdateJob(service, dcli, c); err != nil {
				Logger.Error().Err(err).Msgf("Cannot update job for service %s", event.Service)
				continue
			}
			Logger.Debug().Msgf("Number of cronjob tasks : %d", len(c.Entries()))
		}
	}
}
