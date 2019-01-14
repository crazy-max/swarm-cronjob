package main

import (
	"runtime"

	"github.com/crazy-max/swarm-cronjob/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
