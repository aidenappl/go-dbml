package dbml

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aidenappl/go-dbml/diag"
)

// TestParseFileContextCancel: a canceled compile returns within
// max(50 ms, 5% of an uncanceled run measured in the same test) of the
// moment cancel() actually ran (timer delays under load are not the
// compiler's), well before the full compile, with only a CANCELED
// diagnostic. Under DBML_PERF_STRICT=1 the latency is also checked against
// the cancel's planned time.
func TestParseFileContextCancel(t *testing.T) {
	var b strings.Builder
	for i := 0; b.Len() < 2<<20; i++ {
		fmt.Fprintf(&b, "Table t%d {\n  id int [pk]\n  p int [ref: > t%d.id, note: 'n']\n}\n", i, i)
	}
	src := b.String()
	full := time.Now()
	if r := ParseFileWith(DefaultFilepath, src, Options{}); r.Database == nil {
		t.Fatal(r.Diagnostics.Errors()[:1])
	}
	whole := time.Since(full)
	allowed := max(50*time.Millisecond, whole*5/100)
	for _, frac := range []int{0, 10, 33, 66} {
		// best of three attempts: a scheduling spike on a loaded machine
		// is not a missing poll (which would show on every attempt)
		var errs []string
		for attempt := 0; attempt < 3; attempt++ {
			errs = cancelOnce(t, src, whole, allowed, frac)
			if len(errs) == 0 {
				break
			}
		}
		for _, e := range errs {
			t.Error(e)
		}
	}
}

// cancelOnce cancels a compile of src at frac% of whole and returns what
// went wrong.
func cancelOnce(t *testing.T, src string, whole, allowed time.Duration, frac int) []string {
	var errs []string
	{
		after := whole * time.Duration(frac) / 100
		ctx, cancel := context.WithCancel(context.Background())
		var firedNs atomic.Int64
		start := time.Now()
		timer := time.AfterFunc(after, func() {
			firedNs.Store(int64(time.Since(start)))
			cancel()
		})
		r := ParseFileContext(ctx, DefaultFilepath, src, Options{})
		el := time.Since(start)
		timer.Stop()
		cancel()
		fired := time.Duration(firedNs.Load())
		if r.Database != nil {
			// finished before the cancel ran: only acceptable for late cancels
			if firedNs.Load() != 0 && fired < el-allowed {
				errs = append(errs, fmt.Sprintf("cancel at %v (ran at %v): compile completed anyway after %v", after, fired, el))
			}
			return errs
		}
		if len(r.Diagnostics) != 1 || r.Diagnostics[0].Code != diag.CANCELED || r.Partial != nil {
			return append(errs, fmt.Sprintf("cancel at %v: %d diags, first %v, partial %v", after, len(r.Diagnostics), r.Diagnostics[0].Code, r.Partial != nil))
		}
		latency := el - fired
		t.Logf("cancel planned at %v ran at %v: returned after %v (latency %v, allowed %v; whole compile %v)", after, fired, el, latency, allowed, whole)
		if latency > allowed {
			errs = append(errs, fmt.Sprintf("cancel ran at %v but the compile returned after %v (latency %v > %v)", fired, el, latency, allowed))
		}
		if perfStrict() && el > after+allowed {
			errs = append(errs, fmt.Sprintf("strict: cancel planned at %v returned after %v", after, el))
		}
	}
	return errs
}
