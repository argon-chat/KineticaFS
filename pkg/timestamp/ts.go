package timestamp

import "time"

func CurrentTimestampAt(now time.Time) uint32 {
	ref := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	seconds := now.UTC().Sub(ref).Seconds()
	if seconds < 0 {
		return 0
	}
	return uint32(seconds)
}

func CurrentTimestamp() uint32 {
	return CurrentTimestampAt(time.Now())
}
