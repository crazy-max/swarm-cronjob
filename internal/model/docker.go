package model

import (
	"time"

	"github.com/docker/docker/api/types/swarm"
)

// ServiceEvent represents attributes of a service event
type ServiceEvent struct {
	Service     string `mapstructure:"name"`
	UpdateState struct {
		Old string `mapstructure:"updatestate.old"`
		New string `mapstructure:"updatestate.new"`
	} `mapstructure:",squash"`
}

// ServiceListArgs are options to list services
type ServiceListArgs struct {
	Name   string
	Labels []string
}

// ServiceInfo represents attributes of a service
type ServiceInfo struct {
	Raw          swarm.Service
	ID           string
	Name         string
	Image        string
	Mode         ServiceMode
	Labels       map[string]string
	Actives      uint64
	Replicas     uint64
	Rollback     bool
	UpdatedAt    time.Time
	UpdateStatus string
}

// ServiceMode is a service mode
type ServiceMode string

// Services modes available
const (
	ServiceModeReplicated = ServiceMode("replicated")
	ServiceModeGlobal     = ServiceMode("global")
)

// TaskInfo represents attributes of a task
type TaskInfo struct {
	swarm.Task
	NodeName    string
	ServiceName string
	Image       string
}
