package model

import "github.com/alecthomas/kong"

// Cli holds command line args, flags and cmds
type Cli struct {
	Version      kong.VersionFlag
	LogLevel     string   `kong:"name='log-level',env='LOG_LEVEL',default='info',help='Set log level.'"`
	LogJSON      bool     `kong:"name='log-json',env='LOG_JSON',default='false',help='Enable JSON logging output.'"`
	ServiceLabel []string `kong:"name='service-label',env='SERVICE_LABELS',help='Additional service label filter. Repeat for multiple labels.'"`
}
