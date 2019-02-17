package model

type Job struct {
	Name        string
	Enable      bool
	Schedule    string
	SkipRunning bool
}
