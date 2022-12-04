package xctrainer

import (
	"strconv"
)

type Line struct {
	line   string
	cursor int
}

func NewLine(line string) *Line {
	l := new(Line)
	l.line = line
	l.cursor = 0
	return l
}

func (l *Line) TakeByte() byte {
	ret := l.line[l.cursor]
	l.cursor++
	return ret
}

func (l *Line) TakeString(len int) string {
	ret := l.line[l.cursor : l.cursor+len]
	l.cursor += len
	return ret
}

func (l *Line) TakeInt(len int) int {
	ret, err := strconv.ParseInt(l.line[l.cursor:l.cursor+len], 10, 64)
	if err != nil {
		return 0
	}
	l.cursor += len
	return int(ret)
}

func (l *Line) TakeHexInt(len int) int {
	ret, err := strconv.ParseInt(l.line[l.cursor:l.cursor+len], 16, 64)
	if err != nil {
		return 0
	}
	l.cursor += len

	return int(ret)
}

func (l *Line) TakeSignedHexInt(len int) int {
	ret := l.TakeHexInt(len)
	mask := 1 << (len*4 - 1)
	if mask&ret != 0 {
		ret = ^ret
		ret++
		ret &= ^((^0) << (len * 4))
		ret = -ret
	}
	return ret
}

func controlCode(line string) byte {
	var cc byte

	cc = 0

	for _, ch := range line {
		cc = cc ^ byte(ch)
	}
	return cc
}

func (l *Line) ControlCode() bool {
	cc, err := strconv.ParseInt(l.line[l.cursor:l.cursor+2], 16, 8)

	if err != nil {
		return false
	}
	return byte(cc) == controlCode(l.line[:l.cursor])
}
