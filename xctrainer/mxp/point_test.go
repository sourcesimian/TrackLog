package mxp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sourcesimian/TrackLog/xctrainer/mxp"
	"time"
)

type test1_t struct {
	line    string
	time    time.Time
	lat     float64
	lon     float64
	altGnss int
	altBaro int
}

func TestPoint(t *testing.T) {
	tm := func(s string) time.Time {
		r, _ := time.Parse(time.RFC3339, s)
		return r
	}

	const epsilon = 0.00001
	const header = "@xc351100A8032211130744270128EE012**1F"
	var tests = []test1_t{
		{
			line:    "0001000103F37EDB06EEAF01D302097E",
			time:    tm("2022-11-13T07:44:28Z"),
			lat:     -34.1454,
			lon:     18.9299,
			altGnss: 467,
			altBaro: 521,
		},
	}

	for _, test := range tests {

		tr := mxp.NewTrack()
		tr.ParseHeader(header)

		p := mxp.NewPoint(tr, test.line)

		assert.NotNil(t, p, "Bad Line")

		assert.Equal(t, test.time, p.Time(), "Time")
		assert.InEpsilon(t, test.lat, p.Lat, epsilon, "Lat")
		assert.InEpsilon(t, test.lon, p.Lon, epsilon, "Long")
		assert.Equal(t, test.altGnss, p.AltGNSS, "AltGNSS")
		assert.Equal(t, test.altBaro, p.AltBaro, "AltBaro")
	}
}
