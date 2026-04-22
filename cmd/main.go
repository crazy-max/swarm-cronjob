package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	_ "time/tzdata"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/swarm-cronjob/internal/app"
	"github.com/crazy-max/swarm-cronjob/internal/logging"
	"github.com/crazy-max/swarm-cronjob/internal/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func run() error {
	cli := model.Cli{}
	_ = kong.Parse(&cli,
		kong.Name("swarm-cronjob"),
		kong.Description(`Create jobs on a time-based schedule on Swarm. More info: https://github.com/crazy-max/swarm-cronjob`),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	logging.Configure(&cli)
	log.Info().Msgf("Starting swarm-cronjob %s", version)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc, err := app.New()
	if err != nil {
		return errors.Wrap(err, "cannot initialize swarm-cronjob")
	}

	if err := sc.Run(ctx); err != nil {
		return errors.Wrap(err, "cannot run swarm-cronjob")
	}

	if cause := context.Cause(ctx); cause != nil {
		log.Warn().Msg(cause.Error())
	}

	return nil
}
