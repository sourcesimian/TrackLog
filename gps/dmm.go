package gps

import (
	"bufio"
	"fmt"
	"math"
)

type DMm struct {
	Deg   int
	Min   int
	Dmin3 int
}

func (d *DMm) Write(w *bufio.Writer, lat bool) {
	var h byte
	var width int
	deg := d.Deg
	if lat {
		width = 2
		if deg < 0 {
			h = 'S'
			deg = -deg
		} else {
			h = 'N'
		}
	} else {
		width = 3
		if deg < 0 {
			h = 'W'
			deg = -deg
		} else {
			h = 'E'
		}
	}
	fmt.Fprintf(w, "%0*d%02d%03d%c", width, deg, d.Min, d.Dmin3, h)
}

func (d *DMm) FromDD(dd float64) {
	sign := 1
	if dd < 0 {
		dd = -dd
		sign = -1
	}

	deg, frac := math.Modf(dd)
	d.Deg = int(deg) * sign
	min, frac := math.Modf(frac * 60.0)
	d.Min = int(min)
	dmin3, _ := math.Modf(frac * 1000.0)
	d.Dmin3 = int(dmin3)
}
