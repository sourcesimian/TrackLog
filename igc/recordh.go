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
}

func (r *HRecord) Write(w *bufio.Writer) {
	fmt.Fprintf(w, "H%c%s%s: %s%s", r.source, r.code, r.longName, r.value(), CRLF)
}
