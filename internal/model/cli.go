package model

import "github.com/alecthomas/kong"

// Cli holds command line args, flags and cmds
type Cli struct {
	Version  kong.VersionFlag
	Timezone string `kong:"name='timezone',env='TZ',default='UTC',help='Timezone assigned to swarm-cronjob.'"`
	LogLevel string `kong:"name='log-level',env='LOG_LEVEL',default='debug',help='Set log level.'"`
	LogJSON  bool   `kong:"name='log-json',env='LOG_JSON',default='false',help='Enable JSON logging output.'"`
}
