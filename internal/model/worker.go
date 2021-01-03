package model

// Job holds service job details
type Job struct {
	Name         string
	Enable       bool
	Schedule     string
	SkipRunning  bool
	RegistryAuth bool
	Replicas     uint64
}
