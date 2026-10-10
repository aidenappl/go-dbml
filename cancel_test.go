package dbml

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aidenappl/go-dbml/diag"
)

func TestParseFileContextCancel(t *testing.T) {
	var b strings.Builder
	for i := 0; b.Len() < 4<<20; i++ {
		fmt.Fprintf(&b, "Table t%d {\n  id int [pk]\n  p int [ref: > t%d.id, note: 'n']\n}\n", i, i)
	}
	src := b.String()
	full := time.Now()
	if r := ParseFileWith(DefaultFilepath, src, Options{}); r.Database == nil {
		t.Fatal(r.Diagnostics.Errors()[:1])
	}
	whole := time.Since(full)
	for _, after := range []time.Duration{0, whole / 10, whole / 3, whole * 2 / 3} {
		ctx, cancel := context.WithCancel(context.Background())
		timer := time.AfterFunc(after, cancel)
		start := time.Now()
		r := ParseFileContext(ctx, DefaultFilepath, src, Options{})
		el := time.Since(start)
		timer.Stop()
		cancel()
		if r.Database != nil {
			// finished before the cancel landed: only acceptable for late cancels
			if after < whole/2 {
				t.Errorf("cancel after %v: compile completed in %v", after, el)
			}
			continue
		}
		if len(r.Diagnostics) != 1 || r.Diagnostics[0].Code != diag.CANCELED || r.Partial != nil {
			t.Errorf("cancel after %v: %d diags, first %v, partial %v", after, len(r.Diagnostics), r.Diagnostics[0].Code, r.Partial != nil)
		}
		if el > after+whole/10+50*time.Millisecond && !raceEnabled { // -race distorts latency
			t.Errorf("cancel after %v returned after %v (whole compile %v)", after, el, whole)
		}
		t.Logf("cancel after %v: returned after %v (whole compile %v)", after, el, whole)
	}
}
