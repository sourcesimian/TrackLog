package igc

import (
	"bufio"
	"fmt"
	"math"
)

type BRecord struct {
	point *Point
	fix   int
}

func NewBRecord(p *Point) *BRecord {
	r := &BRecord{
		point: p,
		fix:   3,
	}
	return r
}

func (r *BRecord) Write(w *bufio.Writer) {
	// B HHMMSS DDMMmmm? DDDMMmmm? F PPPPP GGGGG CRLF
	fmt.Fprintf(w, "B%s", r.point.Time.Format("150405"))

	r.point.Lat.Write(w, true)
	r.point.Lon.Write(w, false)
	switch r.fix {
	case 2:
		fmt.Fprintf(w, "V") // 2D fix (no GPS ALt) or no GPS data
	case 3:
		fmt.Fprintf(w, "A") //  3D fix
	default:
		fmt.Fprintf(w, "?")
	}

	fmt.Fprintf(w, "%05d", r.point.AltBaro)
	fmt.Fprintf(w, "%05d", r.point.AltGNSS)
	fmt.Fprintf(w, "%s", CRLF)
}

func DD2DMm(deg float64) (int, int, int) {
	sign := 1
	if deg < 0 {
		deg = -deg
		sign = -1
	}

	d, frac := math.Modf(deg)
	m, frac := math.Modf(frac * 60.0)
	dm, _ := math.Modf(frac * 1000.0)

	return sign * int(d), int(m), int(dm)
}
