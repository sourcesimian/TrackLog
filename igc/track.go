package igc

// ref: https://xp-soaring.github.io/igc_file_format/igc_format_2008.html

import (
	"bufio"
	"fmt"
)

const (
	CRLF = "\r\n"
)

type Track struct {
	Header
	points []Point
}

func NewTrack() *Track {
	track := &Track{}
	return track
}

func (t *Track) Write(w *bufio.Writer) {
	fmt.Fprintf(w, "A%s%s", t.Manufacturer, CRLF)
	t.WriteHeader(w)

	for _, point := range t.points {
		r := NewBRecord(&point)
		r.Write(w)
	}
	t.WriteFooter(w)
}

func (t *Track) AddPoint(p *Point) {
	t.points = append(t.points, *p)
}

func (t *Track) WriteFooter(w *bufio.Writer) {
	fmt.Fprintf(w, "LXGD https://github.com/sourcesimian/TrackLog")
}
