package igc

import (
	"time"

	"github.com/sourcesimian/TrackLog/gps"
)

type Point struct {
	Time    time.Time
	Lat     gps.DMm
	Lon     gps.DMm
	AltGNSS int
	AltBaro int
}
