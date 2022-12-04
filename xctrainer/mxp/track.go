package mxp

import (
	"bufio"
	"log"
	"math"
	"strings"
	"time"

	"github.com/sourcesimian/TrackLog/igc"
	"github.com/sourcesimian/TrackLog/xctrainer"
)

type Track struct {
	device        string
	version       string
	serial        string
	flight        int
	time          time.Time
	pointCount    int
	pointInterval int
	height        int

	points []*Point
}

func NewTrack() *Track {
	t := new(Track)
	return t
}

func (t *Track) Load(r *bufio.Reader, progress func(int)) bool {
	points := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal(err)
			}
			return false
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}

		if len(t.device) == 0 {
			if t.ParseHeader(line) == false {
				return false
			}
		} else {
			if t.addPoint(line) == false {
				return false
			}
			points++
			progress(points)
		}
		if points == t.pointCount {
			break
		}
	}
	return points > 0
}

func (t *Track) ToIGC(trk *igc.Track) {
	trk.Manufacturer = "XGD"
	trk.Date = t.time
	trk.GPSDatum = "WGS-84"

	for _, point := range t.points {
		lat := DtoDMm(point.lat)
		lon := DtoDMm(point.lon)

		p := &igc.Point{
			Time:    point.Time(),
			Lat:     lat,
			Lon:     lon,
			AltGNSS: point.AltGNSS,
			AltBaro: point.AltBaro,
		}

		trk.AddPoint(p)
	}
}

func (t *Track) ParseHeader(line string) bool {
	// ref: http://www.penguin.cz/~ondrap/aircotec.php
	//
	// @TTVVVVSSSSFFYYMMDDhhmmssGGPPPPIIB**CC
	//
	// TT (char): Type of GPS. 'xc' for XC trainer, 'tn' for Top-navigator.
	// VVVV (hexadecimal): Version. In this case: 2.4-05 for XC trainer, v.2405 for TN.
	// SSSS (hexadecimal): Serial number.
	// FF (hexadecimal): Flight number mod 256.
	// YYMMDDhhmmss (decimal): Time of the first point in tracklog (UTC).
	// GG: Geodetic datum, always 1 means WGS-84.
	// PPPP (hexadecimal): Count of points that will be sent following header.
	// II (hexadecimal): Interval in seconds between points. It was chosen by the user when sending the tracklog on the device.
	// CC (hexadecimal): Control code. XOR of all characters before the control code.
	// B (*|0|1|2): Type of height that will be sent:
	// * (older version) or 0 - baro height
	// 1 - GPS height
	// 2 - both baro and GPS height

	l := xctrainer.NewLine(line)

	if l.TakeByte() != '@' {
		return false
	}
	t.device = l.TakeString(2)
	t.version = l.TakeString(4)
	t.serial = l.TakeString(4)
	t.flight = l.TakeHexInt(2)

	t.time = time.Date(2000+l.TakeInt(2), time.Month(l.TakeInt(2)), l.TakeInt(2), l.TakeInt(2), l.TakeInt(2), l.TakeInt(2), 0, time.UTC)

	gg := l.TakeString(2)
	if gg != "01" {
		return false
	}

	t.pointCount = l.TakeHexInt(4)
	t.pointInterval = l.TakeHexInt(2)

	b := l.TakeByte()
	switch b {
	case '*' | '0':
		t.height = 0
	case '1':
		t.height = 1
	case '2':
		t.height = 2
	}

	l.TakeString(2) // **

	return l.ControlCode()
}

func (t *Track) Start() time.Time {
	return t.time
}

func (t *Track) addPoint(line string) bool {
	r := NewPoint(t, line)
	if r == nil {
		return false
	}
	t.points = append(t.points, r)
	return true
}

func (t *Track) Duration() time.Duration {
	return time.Duration(t.points[len(t.points)-1].elapsed) * time.Second
}

func (t *Track) AltRange() (int, int) {
	min := int(math.MaxInt32)
	max := 0

	for _, p := range t.points {
		if p.AltGNSS > max {
			max = p.AltGNSS
		}
		if p.AltGNSS < min {
			min = p.AltGNSS
		}
	}
	return min, max
}
