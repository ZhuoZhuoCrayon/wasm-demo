package query

import (
	"testing"
	"time"
)

func TestNewTimeRangeIterator(t *testing.T) {
	tests := []struct {
		beginTimeStr string
		endTimeStr   string
		interval     time.Duration
		want         int
	}{
		{
			beginTimeStr: "2024-04-02 14:00:00",
			endTimeStr:   "2024-04-02 14:03:00",
			interval:     1 * time.Minute,
			want:         3,
		},
	}

	for _, tt := range tests {
		tr, err := NewTimeRangeFromStr(tt.beginTimeStr, tt.endTimeStr)
		if err != nil {
			t.Errorf("failed to New TimeRange: %v", err)
		}
		ti := NewTimeRangeIterator(&tr, tt.interval)
		assertTimeRanges(t, ti, tt.want)
	}
}

func TestNewTimeRangeIteratorFromSegments(t *testing.T) {
	tests := []struct {
		beginTimeStr string
		endTimeStr   string
		segments     int
		want         int
	}{
		{
			beginTimeStr: "2024-04-02 14:00:00",
			endTimeStr:   "2024-04-02 14:03:00",
			segments:     3,
			want:         3,
		},
	}
	for _, tt := range tests {
		tr, err := NewTimeRangeFromStr(tt.beginTimeStr, tt.endTimeStr)
		if err != nil {
			t.Errorf("failed to New TimeRange: %v", err)
		}
		ti := NewTimeRangeIteratorFromSegments(&tr, tt.segments)
		assertTimeRanges(t, ti, tt.want)
	}
}

func assertTimeRanges(tb testing.TB, ti *TimeRangeIterator, want int) {
	tb.Helper()
	actual := 0
	for {
		_, end := ti.Next()
		if end {
			break
		}
		actual++
	}
	if actual != want {
		tb.Errorf("timeRange length want %d but got %d", want, actual)
	}
}
