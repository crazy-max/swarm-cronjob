//go:build !windows
// +build !windows

package main

import (
	"golang.org/x/sys/unix"
)

const (
	SIGTERM = unix.SIGTERM
)
