package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/swarm-cronjob/internal/app"
	"github.com/crazy-max/swarm-cronjob/internal/logging"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/rs/zerolog/log"
)

var (
	sc      *app.SwarmCronjob
	cli     model.Cli
	version = "dev"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse command line
	_ = kong.Parse(&cli,
		kong.Name("swarm-cronjob"),
		kong.Description(`Create jobs on a time-based schedule on Swarm. More info: https://github.com/crazy-max/swarm-cronjob`),
		kong.UsageOnError(),
		kong.Vars{
			"version": fmt.Sprintf("%s", version),
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	// Load timezone location
	location, err := time.LoadLocation(cli.Timezone)
	if err != nil {
		log.Panic().Err(err).Msgf("Cannot load timezone %s", cli.Timezone)
	}

	// Init
	logging.Configure(&cli, location)
	log.Info().Msgf("Starting swarm-cronjob %s", version)

	// Handle os signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-channel
		if sc != nil {
			sc.Close()
		}
		log.Warn().Msgf("Caught signal %v", sig)
		os.Exit(1)
	}()

	// Init
	sc, err = app.New(location)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize swarm-cronjob")
	}

	// Run
	if err := sc.Run(); err != nil {
		log.Panic().Err(err).Msg("")
	}
}
