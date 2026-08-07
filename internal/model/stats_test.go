package model_test

import (
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

// TestEstimatedRowCountSeparatesUnknownFromEmpty guards a bug found by running
// against a real database: since PostgreSQL 14, reltuples is -1 for a table
// that was never ANALYZEd. Read as a count it makes the table look empty, which
// in scoring turns every unanalyzed table into a small domain table.
func TestEstimatedRowCountSeparatesUnknownFromEmpty(t *testing.T) {
	tests := []struct {
		name      string
		reltuples int64
		wantRows  int64
		wantKnown bool
	}{
		{"never analyzed", model.RowsUnknown, 0, false},
		{"genuinely empty", 0, 0, true},
		{"populated", 4_000_000, 4_000_000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := model.TableStats{EstimatedRows: tt.reltuples}

			rows, known := stats.EstimatedRowCount()
			if rows != tt.wantRows || known != tt.wantKnown {
				t.Errorf("EstimatedRowCount() = %d, %v; want %d, %v",
					rows, known, tt.wantRows, tt.wantKnown)
			}
		})
	}
}
