package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/crazy-max/swarm-cronjob/internal/app"
	"github.com/crazy-max/swarm-cronjob/internal/logging"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	sc      *app.SwarmCronjob
	flags   model.Flags
	version = "dev"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse command line
	kingpin.Flag("timezone", "Timezone assigned to the scheduler.").Envar("TZ").Default("UTC").StringVar(&flags.Timezone)
	kingpin.Flag("log-level", "Set log level.").Envar("LOG_LEVEL").Default("info").StringVar(&flags.LogLevel)
	kingpin.Flag("log-json", "Enable JSON logging output.").Envar("LOG_JSON").Default("false").BoolVar(&flags.LogJson)
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(version).Author("CrazyMax")
	kingpin.CommandLine.Name = "swarm-cronjob"
	kingpin.CommandLine.Help = `Create jobs on a time-based schedule on Swarm. More info on https://github.com/crazy-max/swarm-cronjob`
	kingpin.Parse()

	// Load timezone location
	location, err := time.LoadLocation(flags.Timezone)
	if err != nil {
		log.Panic().Err(err).Msgf("Cannot load timezone %s", flags.Timezone)
	}

	// Init
	logging.Configure(&flags, location)
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
