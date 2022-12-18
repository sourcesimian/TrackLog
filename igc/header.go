package igc

import (
	"bufio"
	"time"
)

type Header struct {
	Manufacturer     string
	Date             time.Time
	Pilot            string
	GliderType       string
	GliderID         string
	GPSDatum         string
	CompetitionID    string
	CompetitionClass string
}

func (h *Header) WriteHeader(w *bufio.Writer) {
	var hrs = []HRecord{
		{
			source: 'F',
			code:   "DTE",
			value: func() string {
				return h.Date.Format("020106")
			},
			nocolon: true,
		},
		{
			source:   'F',
			code:     "PLT",
			longName: "PILOT",
			value: func() string {
				return h.Pilot
			},
		},
		{
			source:   'F',
			code:     "GTY",
			longName: "GLIDERTYPE",
			value: func() string {
				return h.GliderType
			},
		},
		{
			source:   'F',
			code:     "GID",
			longName: "GLIDERID",
			value: func() string {
				return h.GliderID
			},
		},
		{
			source:   'F',
			code:     "DTM",
			longName: "100GPSDATUM",
			value: func() string {
				return h.GPSDatum
			},
		},
		{
			source:   'F',
			code:     "CID",
			longName: "COMPETITIONID",
			value: func() string {
				return h.CompetitionID
			},
		},
		{
			source:   'F',
			code:     "CCL",
			longName: "COMPETITION CLASS",
			value: func() string {
				return h.CompetitionClass
			},
		},
		{
			source:   'F',
			code:     "SIT",
			longName: "SITE",
			value: func() string {
				return ""
			},
		},
	}

	for _, rec := range hrs {
		rec.Write(w)
	}
}
