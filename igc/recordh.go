package igc

import (
	"bufio"
	"fmt"
)

type HRecord struct {
	source   byte
	code     string
	longName string
	value    func() string
	nocolon  bool
}

func (r *HRecord) Write(w *bufio.Writer) {
	sep := ": "
	if r.nocolon {
		sep = ""
	}
	long := ""
	if len(r.longName) > 0 {
		long = r.longName
	}
	fmt.Fprintf(w, "H%c%s%s%s%s%s", r.source, r.code, long, sep, r.value(), CRLF)
}
