package export_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/aidenappl/go-dbml/export"
	"github.com/aidenappl/go-dbml/model"
)

// A 20k-table model with 20k refs exports in linear time (buildCore uses
// map lookups and map-based ref de-duplication).
func TestExport20kTables(t *testing.T) {
	const n = 20000
	db := &model.Database{}
	tru := true
	for i := 0; i < n; i++ {
		name := "t" + strconv.Itoa(i)
		db.Tables = append(db.Tables, &model.Table{Name: name, Fields: []*model.Column{
			{Name: "id", Type: model.ColumnType{TypeName: "int"}, PK: &tru},
			{Name: "p_id", Type: model.ColumnType{TypeName: "int"}},
		}})
		db.Refs = append(db.Refs, &model.Ref{Endpoints: [2]model.RefEndpoint{
			{TableName: "t" + strconv.Itoa((i+1)%n), FieldNames: []string{"id"}, Relation: "1"},
			{TableName: name, FieldNames: []string{"p_id"}, Relation: "*"},
		}})
	}
	limit := 1500 * time.Millisecond
	if raceEnabled {
		limit *= 15
	}
	for _, d := range export.Dialects {
		start := time.Now()
		if _, err := export.SQL(db, d); err != nil {
			t.Fatal(err)
		}
		el := time.Since(start)
		t.Logf("%s: %v", d, el)
		if el > limit {
			t.Errorf("%s export took %v", d, el)
		}
	}
}
