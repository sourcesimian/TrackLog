package mxp

import (
	"time"

	"github.com/sourcesimian/TrackLog/xctrainer"
	"github.com/sourcesimian/TrackLog/gps"
)

type Point struct {
	t       *Track
	line    int
	elapsed int
	flags   byte
	lat     int
	lon     int
	AltGNSS int
	AltBaro int
}

func NewPoint(t *Track, line string) *Point {
	// ref: http://www.penguin.cz/~ondrap/aircotec.php
	//
	// ppppttttFFllllllooooooaaaaCC
	//
	// All fields are in hexadecimal format.
	// pppp: Line number, starts from 0, is incremented with each line.
	// tttt: Number of seconds elapsed from the beginning of tracklog. 'TTTT'+'tttt' is the time of the point.
	// FF: Flags. Bit-wise OR of following values:
	// For XC: 0x01 - the GPS had sattelite signal, 0x02 - 3D fix, 0x80 - Mark
	// For TN: 0x01 - GPSspeed/course valid, 0x02 - always 0, 0x80 - Mark
	// llllll: Latitude. Divide the number by 24000.0 to get degrees. (2's complement format for negative numbers).
	// oooooo: Longitude, same format as latitude.
	// aaaa: Altitude in meters. (2's complement format for negative numbers).
	// CC: Control code. XOR of all characters before the Control code.
	// When both GPS and baro height are sent, the aaaa becomes aaaabbbb instead with GPS height first and baro height second.
	p := new(Point)
	p.t = t

	l := xctrainer.NewLine(line)
	p.line = l.TakeHexInt(4)
	p.elapsed = l.TakeHexInt(4)

	p.flags = byte(l.TakeHexInt(2))

	p.lat = l.TakeSignedHexInt(6)
	p.lon = l.TakeSignedHexInt(6)
	p.AltGNSS = l.TakeHexInt(4)
	if t.height == 2 {
		p.AltBaro = l.TakeHexInt(4)
	}

	if !l.ControlCode() {
		return nil
	}

	return p
}

func (p *Point) Time() time.Time {
	return p.t.Start().Add(time.Second * time.Duration(p.elapsed))
}


func DtoDMm(d int) gps.DMm {
	sign := 1
	if d < 0 {
		sign = -1
		d = -d
	}
	factor := 24000
	deg := int(d/factor)
	d -= (deg*factor)
	d *= 60
	min := int(d/factor)
	d -= (min*factor)
	d *= 1000
	dmin3 := int(d/factor)

	return gps.DMm {
		Deg: deg * sign,
		Min: min,
		Dmin3: dmin3,
	}
}