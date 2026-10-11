//go:build !unix

package dbml

import "time"

// cpuNow falls back to wall-clock time where getrusage is unavailable.
func cpuNow() time.Duration { return time.Duration(time.Now().UnixNano()) }
