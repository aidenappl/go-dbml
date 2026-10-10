package importer

import (
	"os"
	"testing"
)

// BenchmarkPostgresDump imports the 8.7k-line Holistics pg_dump fixture.
func BenchmarkPostgresDump(b *testing.B) {
	src, err := os.ReadFile("../testdata/holistics/dbml-core/importer/postgres_importer/input/db_dump.in.sql")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := CoreDBML(string(src), Postgres); err != nil {
			b.Fatal(err)
		}
	}
}
