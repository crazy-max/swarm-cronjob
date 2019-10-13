package model

// Job holds service job details
type Job struct {
	Name        string
	Enable      bool
	Schedule    string
	SkipRunning bool
	Replicas    uint64
}
