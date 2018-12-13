package app

import (
	"strconv"

	"github.com/crazy-max/cron"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type cronjob struct {
	name        string
	enable      bool
	schedule    string
	skipRunning bool
}

type worker struct {
	dcli    *client.Client
	cronjob cronjob
}

func (w *worker) Run() {
	RunService(w.dcli, w.cronjob.name, w.cronjob.skipRunning)
}

// UpdateJob adds or updates cron job service
func UpdateJob(service swarm.Service, dcli *client.Client, c *cron.Cron) error {
	cronjob := cronjob{
		name:        service.Spec.Name,
		enable:      false,
		skipRunning: false,
	}

	for labelKey, labelValue := range service.Spec.Labels {
		switch labelKey {
		case "swarm.cronjob.enable":
			enable, err := strconv.ParseBool(labelValue)
			if err != nil {
				Logger.Error().Err(err).Msgf("Cannot parse %s value for service %s", labelKey, service.Spec.Name)
			}
			cronjob.enable = enable
		case "swarm.cronjob.schedule":
			cronjob.schedule = labelValue
		case "swarm.cronjob.skip-running":
			skipRunning, err := strconv.ParseBool(labelValue)
			if err != nil {
				Logger.Error().Err(err).Msgf("Cannot parse %s value for service %s", labelKey, service.Spec.Name)
			}
			cronjob.skipRunning = skipRunning
		}
	}

	// Find existing job
	var jobEntry *cron.Entry
	for _, entry := range c.Entries() {
		if entry.Name == service.Spec.Name {
			jobEntry = entry
			break
		}
	}

	// Remove disabled job
	if jobEntry != nil && !cronjob.enable {
		Logger.Debug().Msgf("Disable cronjob for service %s", service.Spec.Name)
		return c.Remove(jobEntry.Name)
	}

	// Add/Update job
	if jobEntry != nil {
		if err := c.Remove(jobEntry.Name); err != nil {
			return err
		}
		Logger.Debug().Msgf("Update cronjob for service %s with schedule %s", service.Spec.Name, cronjob.schedule)
	} else {
		Logger.Info().Msgf("Add cronjob for service %s with schedule %s", service.Spec.Name, cronjob.schedule)
	}

	return c.AddJob(cronjob.schedule, &worker{dcli, cronjob}, cronjob.name)
}
