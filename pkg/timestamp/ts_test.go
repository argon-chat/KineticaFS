package timestamp

import (
	"testing"
	"time"
)

func TestCurrentTimestampAt_01_01_2025(t *testing.T) {
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if got := CurrentTimestampAt(testTime); got != 0 {
		t.Errorf("Expected 0, got %d", got)
	}
}

func TestCurrentTimestampAt_01_01_2024(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if got := CurrentTimestampAt(testTime); got != 0 {
		t.Errorf("Expected 0, got %d", got)
	}
}

func TestCurrentTimestampAt_01_01_2026(t *testing.T) {
	testTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if got := CurrentTimestampAt(testTime); got != 31536000 {
		t.Errorf("Expected 31536000, got %d", got)
	}
}
