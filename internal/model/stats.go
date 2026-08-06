package model

import (
	"encoding/json"
	"math"
	"time"
)

// TableStats is size and usage information about a table.
type TableStats struct {
	// EstimatedRows is pg_class.reltuples: an estimate, not a count.
	EstimatedRows int64 `json:"estimated_rows"`

	// TotalBytes includes indexes and TOAST.
	TotalBytes int64 `json:"total_bytes"`

	Usage UsageStats `json:"usage"`
}

// UsageCounters are the raw activity counters from pg_stat_user_tables.
type UsageCounters struct {
	SeqScans int64 `json:"seq_scans"`
	IdxScans int64 `json:"idx_scans"`
	Inserts  int64 `json:"inserts"`
	Updates  int64 `json:"updates"`
	Deletes  int64 `json:"deletes"`
}

// UsageStats binds activity counters to the moment they started counting.
//
// A usage counter without its reset timestamp is not interpretable, and acting
// on one is the most dangerous mistake this tool could make. The counters reset
// on pg_stat_reset() and on pg_upgrade, and they are per node: a table read
// exclusively on a read replica shows zero reads on the primary. "This table
// has not been used in years" is not a conclusion these numbers can support on
// their own.
//
// So the counters are unexported and the only way to obtain them is Counters,
// which hands back the reset context in the same call. This is deliberate
// friction: a caller cannot accidentally read SeqScans without being handed the
// reason it might be lying.
type UsageStats struct {
	counters   UsageCounters
	resetAt    time.Time
	resetKnown bool
}

// NewUsageStats pairs counters with the moment statistics were last reset.
func NewUsageStats(c UsageCounters, resetAt time.Time) UsageStats {
	return UsageStats{counters: c, resetAt: resetAt, resetKnown: true}
}

// NewUsageStatsUnknownReset records counters whose reset moment could not be
// determined. Callers must treat the counters as uninterpretable.
func NewUsageStatsUnknownReset(c UsageCounters) UsageStats {
	return UsageStats{counters: c}
}

// Counters returns the activity counters together with the moment they started
// counting. The third result is false when that moment is unknown, in which
// case the counters carry no meaning and must not drive any finding.
func (u UsageStats) Counters() (UsageCounters, time.Time, bool) {
	return u.counters, u.resetAt, u.resetKnown
}

// Interpretable reports whether the counters can support a conclusion.
func (u UsageStats) Interpretable() bool { return u.resetKnown }

type usageStatsJSON struct {
	Counters UsageCounters `json:"counters"`
	ResetAt  *time.Time    `json:"stats_reset_at"`
}

// MarshalJSON always emits the reset moment alongside the counters, so the
// serialized form carries the same warning the Go API does. A null reset means
// the counters are not interpretable.
func (u UsageStats) MarshalJSON() ([]byte, error) {
	out := usageStatsJSON{Counters: u.counters}
	if u.resetKnown {
		t := u.resetAt
		out.ResetAt = &t
	}
	return json.Marshal(out)
}

// UnmarshalJSON restores counters and reset context.
func (u *UsageStats) UnmarshalJSON(b []byte) error {
	var in usageStatsJSON
	if err := json.Unmarshal(b, &in); err != nil {
		return err
	}
	u.counters = in.Counters
	if in.ResetAt != nil {
		u.resetAt, u.resetKnown = *in.ResetAt, true
	} else {
		u.resetAt, u.resetKnown = time.Time{}, false
	}
	return nil
}

// ColumnStats is planner statistics for one column, used to reject impossible
// candidates before any table I/O happens.
//
// The histogram bounds and most common values in pg_stats ARE user data. They
// live in unexported fields so that encoding/json cannot reach them, they never
// appear in any output, log or error, and they exist only long enough to
// produce a number.
type ColumnStats struct {
	// NullFraction is the estimated proportion of NULLs, 0..1.
	NullFraction float64 `json:"null_fraction"`

	// NDistinct follows the pg_stats convention: a positive value is an
	// absolute count, a negative value is a ratio of the table's row count,
	// and zero means unknown. Use EstimatedDistinct to resolve it.
	NDistinct float64 `json:"n_distinct"`

	// Present is false when the column has no statistics at all, typically
	// because the table was never ANALYZEd. The prefilter must stay silent in
	// that case rather than invent a rejection.
	Present bool `json:"present"`

	hasBounds bool
	lowBound  string
	highBound string
}

// SetBounds records the histogram endpoints. The values are user data: they are
// stored unexported, are never serialized, and must never reach any output.
func (c *ColumnStats) SetBounds(low, high string) {
	c.hasBounds, c.lowBound, c.highBound = true, low, high
}

// HasBounds reports whether histogram endpoints are available. It deliberately
// exposes the availability, never the values: there is no exported accessor for
// the bounds themselves and there must not be one. The typed range comparison
// that consumes them belongs inside this package, so the values never cross the
// boundary.
func (c ColumnStats) HasBounds() bool { return c.hasBounds }

// EstimatedDistinct resolves the pg_stats n_distinct convention against a row
// count. The second result is false when the estimate is unavailable.
func (c ColumnStats) EstimatedDistinct(rows int64) (int64, bool) {
	switch {
	case !c.Present, c.NDistinct == 0:
		return 0, false
	case c.NDistinct > 0:
		return int64(math.Round(c.NDistinct)), true
	default:
		if rows <= 0 {
			return 0, false
		}
		return int64(math.Round(-c.NDistinct * float64(rows))), true
	}
}
