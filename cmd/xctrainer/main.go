package main

import (
	"bufio"
	"flag"
	"fmt"
	"path"
	"time"

	"log"
	"os"

	"github.com/tarm/serial"

	"github.com/sourcesimian/TrackLog/igc"
	"github.com/sourcesimian/TrackLog/xctrainer/mxp"
)

func main() {
	os.Exit(xctrainer())
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func xctrainer() int {
	pilotName := flag.String("pilot", "", "Pilot name")
	gliderType := flag.String("gliderType", "", "Glider type")
	mxpInput := flag.String("mxp", "", "MXP input")
	igcFile := flag.String("igc", "", "IGC output file")
	igcDir := flag.String("igcDir", "", "IGC output directory")
	timeout := flag.Int("timeout", 5, "Timeout")
	baudRate := flag.Int("baud", 57600, "Baud rate")
	many := flag.Bool("many", false, "Many track logs")
	out := true

	flag.Parse()
	if len(flag.Args()) > 0 {
		fmt.Println("Bad args")
		usage()
		return 1
	}

	if *mxpInput == "" {
		usage()
		return 1
	}

	// * Input
	inputStat, err := os.Stat(*mxpInput)
	if err != nil {
		log.Fatal(err)
	}
	var reader *bufio.Reader
	if inputStat.Mode().IsRegular() {
		file, err := os.Open(*mxpInput)
		if err != nil {
			log.Fatal(err)
			return 1
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	} else {
		com := &serial.Config{
			Name:        *mxpInput,
			Baud:        *baudRate,
			ReadTimeout: time.Second * time.Duration(*timeout),
		}
		port, err := serial.OpenPort(com)
		if err != nil {
			log.Fatal(err)
		}
		defer port.Close()
		reader = bufio.NewReader(port)
	}

	for {
		if out {
			fmt.Print("Wait ...")
		}

		// * Load
		mxp := mxp.NewTrack()
		if !mxp.Load(reader, func(line int) {
			if out {
				fmt.Printf("\rLoading %d", line)
			}
		}) {
			return 1
		}

		// * Output
		var writer *bufio.Writer
		if *igcDir != "" {
			*igcFile = path.Join(*igcDir, mxp.Start().Format("20060102T150405Z.igc"))
		}
		if *igcFile != "" {
			file, err := os.OpenFile(*igcFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return 1
			}
			defer file.Close()
			writer = bufio.NewWriter(file)
			if out {
				fmt.Printf("\rWriting: %s\n", *igcFile)
			}
		} else {
			writer = bufio.NewWriter(os.Stdout)
			defer os.Stdout.Close()
		}

		// * Write
		igc := igc.NewTrack()
		mxp.ToIGC(igc)
		igc.Pilot = *pilotName
		igc.GliderType = *gliderType

		igc.Write(writer)
		writer.Flush()

		if !*many {
			break
		}
	}
	return 0
}
