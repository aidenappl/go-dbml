package dbml

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/model"
)

// TestJSONRoundTrip checks Database -> JSON -> Database -> JSON is stable
// for every fixture that compiles.
func TestJSONRoundTrip(t *testing.T) {
	var inputs []string
	_ = filepath.WalkDir("testdata", func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(p, ".dbml") && !strings.Contains(p, "/multi/") {
			inputs = append(inputs, p)
		}
		return nil
	})
	n := 0
	for _, p := range inputs {
		src, _ := os.ReadFile(p)
		db, _ := Parse(string(src))
		if db == nil {
			continue
		}
		n++
		first, err := json.Marshal(db)
		if err != nil {
			t.Fatal(err)
		}
		var back model.Database
		if err := json.Unmarshal(first, &back); err != nil {
			t.Fatalf("%s: %v", p, err)
		}
		second, _ := json.Marshal(&back)
		if string(first) != string(second) {
			t.Fatalf("%s: round trip differs\n%s\n%s", p, trunc(string(first)), trunc(string(second)))
		}
	}
	if n < 100 {
		t.Fatalf("only %d inputs compiled", n)
	}
	t.Logf("round-tripped %d databases", n)
}
