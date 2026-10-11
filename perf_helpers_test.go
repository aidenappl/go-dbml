package dbml

import (
	"os"
	"runtime"
	"testing"
	"time"
)

// Performance gates are load-robust by default: they compare the same
// workload at two sizes, measured interleaved in the same test (best of
// five each) in process CPU time (other processes competing for the
// machine do not inflate it), and assert the growth stays near linear. Absolute budgets
// apply only with DBML_PERF_STRICT=1 (`dev perf`, on an idle machine).

func perfStrict() bool { return os.Getenv("DBML_PERF_STRICT") == "1" }

// assertScales times run(small) and run(big) five times each, interleaved,
// and fails when the best big time exceeds slack x (big/small) x the best
// small time (floored at 20 ms, below which timings are noise). With
// slack 2.5 a 4x input may take up to 10x (allocation and GC grow a little
// faster than the input); a quadratic pass (16x) still fails.
func assertScales(t *testing.T, name string, small, big int, slack float64, run func(n int)) {
	t.Helper()
	ts, tb := time.Duration(1<<62), time.Duration(1<<62)
	for i := 0; i < 5; i++ {
		runtime.GC()
		start := cpuNow()
		run(small)
		ts = min(ts, cpuNow()-start)
		runtime.GC()
		start = cpuNow()
		run(big)
		tb = min(tb, cpuNow()-start)
	}
	base := max(ts, 20*time.Millisecond)
	ratio := float64(tb) / float64(base)
	allowed := slack * float64(big) / float64(small)
	t.Logf("%s: n=%d %v, n=%d %v (x%.1f, allowed x%.1f)", name, small, ts, big, tb, ratio, allowed)
	if ratio > allowed {
		t.Errorf("%s grows super-linearly: %d -> %d input took %v -> %v (x%.1f > x%.1f)", name, small, big, ts, tb, ratio, allowed)
	}
}

// strictBudget fails when el exceeds d, only under DBML_PERF_STRICT=1.
func strictBudget(t *testing.T, name string, el, d time.Duration) {
	t.Helper()
	if perfStrict() && el > d {
		t.Errorf("%s took %v (strict budget %v)", name, el, d)
	}
}

// The scaling gate must still catch a quadratic pass: a quadratic loop at
// 4x the size takes ~16x, beyond the 10x a 2.5 slack allows.
func TestScalingGateCatchesQuadratic(t *testing.T) {
	sink := 0
	quadratic := func(n int) {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				sink += i ^ j
			}
		}
	}
	ts, tb := time.Duration(1<<62), time.Duration(1<<62)
	for i := 0; i < 3; i++ {
		start := time.Now()
		quadratic(8000)
		ts = min(ts, time.Since(start))
		start = time.Now()
		quadratic(32000)
		tb = min(tb, time.Since(start))
	}
	ratio := float64(tb) / float64(max(ts, 20*time.Millisecond))
	t.Logf("quadratic 4x input: %v -> %v (x%.1f), sink %d", ts, tb, ratio, sink)
	if ratio <= 2.5*4 {
		t.Errorf("a quadratic workload (x%.1f) would pass the scaling gate (x10)", ratio)
	}
}
