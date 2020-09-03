// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package grafana

import (
	"fmt"
	"strconv"
	"time"

	"github.com/timberio/go-datemath"
)

// Timestamp is the number of milliseconds since the Unix epoch in UTC time zone.
type Timestamp int64

// ParseTimestamp parses the time expression and assigns to the timestamp.
//
// The expression can be one of:
//  - an integer, representing number of milliseconds since Unix epoch in UTC time zone.
//  - a formatted time value, in the form '2006-01-02 15:04:05Z'
//  - a relative time in the form 'now-30m/d' (or those supported by <https://pkg.go.dev/github.com/timberio/go-datemath>)
//    measured in local time zone.
func ParseTimestamp(expr string) (Timestamp, error) {
	if integer, err := strconv.ParseInt(expr, 10, 64); err == nil {
		return Timestamp(integer), nil
	}

	if t, err := time.Parse("2006-01-02 15:04:05Z07:00", expr); err == nil {
		return Timestamp(t.Unix() * 1000), nil
	}

	rel, err := datemath.ParseAndEvaluate(expr, datemath.WithLocation(time.Local))
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp expression '%s':\n%w", expr, err)
	}
	return Timestamp(rel.Unix() * 1000), nil
}
