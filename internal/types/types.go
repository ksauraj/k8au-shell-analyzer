// internal/types/types.go
package types

import "time"

// TimelineEntry represents a single entry in the timeline
type TimelineEntry struct {
	Timestamp time.Time
	Command   string
	Shell     string
}
