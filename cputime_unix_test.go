//go:build unix

package dbml

import (
	"syscall"
	"time"
)

// cpuNow is the process's user+system CPU time: unlike wall-clock time it
// does not grow when other processes compete for the machine.
func cpuNow() time.Duration {
	var ru syscall.Rusage
	if syscall.Getrusage(syscall.RUSAGE_SELF, &ru) != nil {
		return time.Duration(time.Now().UnixNano())
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}
