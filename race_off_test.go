//go:build !race

package dbml

// raceEnabled scales timing assertions under the race detector.
const raceEnabled = false
