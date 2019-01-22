package internal

import (
	"strconv"

	"github.com/crazy-max/cron"
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

// CrudJob adds, updates or removes cron job service
func CrudJob(serviceName string, dcli *client.Client, c *cron.Cron) (bool, error) {
	// Find existing job
	var jobEntry *cron.Entry
	for _, entry := range c.Entries() {
		if entry.Name == serviceName {
			jobEntry = entry
			break
		}
	}

	// Check service exists
	service, err := Service(dcli, serviceName)
	if err != nil {
		if jobEntry != nil {
			Logger.Debug().Msgf("Remove cronjob for service %s", serviceName)
			return true, c.Remove(jobEntry.Name)
		}
		Logger.Debug().Msgf("Service %s does not exist (removed)", serviceName)
		return false, nil
	}

	// Cronjob
	cronjob := cronjob{
		name:        service.Spec.Name,
		enable:      false,
		skipRunning: false,
	}

	// Find swarm.cronjob labels
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

	// Disabled or non-cron service
	if !cronjob.enable {
		if jobEntry != nil {
			Logger.Debug().Msgf("Disable cronjob for service %s", service.Spec.Name)
			return true, c.Remove(jobEntry.Name)
		}
		Logger.Debug().Msgf("Cronjob disabled for service %s", service.Spec.Name)
		return false, nil
	}

	// Add/Update job
	if jobEntry != nil {
		if err := c.Remove(jobEntry.Name); err != nil {
			return true, err
		}
		Logger.Debug().Msgf("Update cronjob for service %s with schedule %s", service.Spec.Name, cronjob.schedule)
	} else {
		Logger.Info().Msgf("Add cronjob for service %s with schedule %s", service.Spec.Name, cronjob.schedule)
	}

	return true, c.AddJob(cronjob.schedule, &worker{dcli, cronjob}, cronjob.name)
}
