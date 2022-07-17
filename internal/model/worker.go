package model

// Job holds service job details
type Job struct {
	Name          string
	Enable        bool
	Schedule      string
	SkipRunning   bool
	RegistryAuth  bool
	QueryRegistry *bool
	Capabilities  []string
	Replicas      uint64
}
