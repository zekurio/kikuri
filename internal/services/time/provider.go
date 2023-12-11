package time

import "time"

// Provider is an interface for providing the current time,
type Provider interface {

	// Now returns the current time.
	Now() time.Time

	// NowUnix returns the current time as unix timestamp.
	NowUnix() int64

	//
}
