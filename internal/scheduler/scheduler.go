package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	// An alias/shadow for each remote defined type
	Uid = uuid.UUID
	Job = gocron.Job

	// Scheduler represents a new active Scheduler object
	Scheduler struct {
		scheduler gocron.Scheduler

		// allow to extend schedule expressions by including seconds
		withSeconds bool
	}
)

// NewScheduler instantiates a new Scheduler object
func NewScheduler(withSeconds bool) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize scheduler")
	}
	return &Scheduler{
		scheduler,
		withSeconds,
	}, nil
}

// Start begins scheduling jobs for execution based on each job's definition
// this method is asynchronous
func (s *Scheduler) Start() {
	s.scheduler.Start()
}

// Stop stops the execution of all jobs in the scheduler.
func (s *Scheduler) Stop() {
	s.scheduler.StopJobs()
}

// CronJob defines a new job using the crontab syntax: `* * * * *`.
// It also supports descriptors such as @monthly, @weekly, @daily, etc.
func (s *Scheduler) CronJob(expression string, task func()) (gocron.Job, error) {
	return s.scheduler.NewJob(
		gocron.CronJob(
			expression, s.withSeconds,
		),
		gocron.NewTask(task),
	)
}

// OneTimeJob defines a new kind of job which will going to run only once
// it allows a custom expression builder to run after some X seconds, such as:
//   - 30, 60 (1 min), 300 (5 mins), etc.
func (s *Scheduler) OneTimeJob(expression string, task func()) (gocron.Job, error) {
	var res time.Duration
	if strings.Contains(expression, "*") {
		return nil, errors.New(fmt.Sprintf("invalid expression syntax for defining a One Time Job: %s", expression))
	}
	expr := strings.TrimSpace(strings.Split(expression, " ")[0])
	if expr != "" {
		val, err := strconv.Atoi(expr)
		if err != nil {
			return nil, err
		}
		res = time.Duration(val) * time.Second
	} else {
		res = 60 * time.Second
	}

	return s.scheduler.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(time.Now().Add(res)),
		),
		gocron.NewTask(task),
	)
}

// CountJobs aggregate and count all the jobs currently in the scheduler.
func (s *Scheduler) CountJobs() int {
	return len(s.scheduler.Jobs())
}

// RemoveJob remove the job by UUID
func (s *Scheduler) RemoveJob(jobId uuid.UUID) {
	s.scheduler.RemoveJob(jobId)
}
