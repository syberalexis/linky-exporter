package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/syberalexis/linky-exporter/pkg/core"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version          = "dev"
	defaultPort      = 9901
	defaultAddress   = "0.0.0.0"
	defaultDevice    = "/dev/serial0"
	defaultBaudRate  = 1200
	defaultFrameSize = 7
	defaultParity    = "ParityNone"
	defaultStopBits  = "Stop1"
)

// Linky-exporter command main
func main() {
	exporter := &core.LinkyExporter{}

	// Globals
	app := kingpin.New(filepath.Base(os.Args[0]), "")
	app.HelpFlag.Short('h')
	app.Version(version)
	app.Action(func(c *kingpin.ParseContext) error { exporter.Run(); return nil })

	// Flags
	app.Flag("address", "Listen address").Default(fmt.Sprintf("%s", defaultAddress)).Short('a').StringVar(&exporter.Address)
	app.Flag("baud", "Baud rate").Default(fmt.Sprintf("%d", defaultBaudRate)).Short('b').IntVar(&exporter.BaudRate)
	app.Flag("device", "Device to read").Default(fmt.Sprintf("%s", defaultDevice)).Short('d').StringVar(&exporter.Device)
	app.Flag("parity", "Serial parity").Default(fmt.Sprintf("%s", defaultParity)).
		HintOptions("ParityNone", "N", "ParityOdd", "O", "ParityEven", "E", "ParityMark", "M", "ParitySpace", "S").StringVar(&exporter.Parity)
	app.Flag("port", "Listen port").Default(fmt.Sprintf("%d", defaultPort)).Short('p').IntVar(&exporter.Port)
	app.Flag("size", "Serial frame size").Default(fmt.Sprintf("%d", defaultFrameSize)).IntVar(&exporter.FrameSize)
	app.Flag("stopbits", "Serial stopbits").Default(fmt.Sprintf("%s", defaultStopBits)).
		HintOptions("Stop1", "1", "Stop1Half", "15", "Stop2", "2").StringVar(&exporter.StopBits)

	// Parsing
	args, err := app.Parse(os.Args[1:])
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		if err != nil {
			log.Error(err)
		}
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		kingpin.MustParse(args, err)
	}
}
