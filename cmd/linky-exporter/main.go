package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/linky-exporter/pkg/core"
	"github.com/syberalexis/linky-exporter/pkg/prom"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Default variables
	version          = "dev"
	defaultPort      = 9901
	defaultAddress   = "0.0.0.0"
	defaultBaudRate  = 1200
	defaultFrameSize = 7
	defaultParity    = "ParityNone"
	defaultStopBits  = "Stop1"

	app        = kingpin.New(filepath.Base(os.Args[0]), "")
	appVersion = app.Version(version)
	help       = app.HelpFlag.Short('h')
	debug      = app.Flag("debug", "Enable debug mode.").Bool()

	address = app.Flag("address", "Listen address").Default(fmt.Sprintf("%s", defaultAddress)).Short('a').String()
	port    = app.Flag("port", "Listen port").Default(fmt.Sprintf("%d", defaultPort)).Short('p').Int()

	auto       = app.Flag("auto", "Automatique mode").Bool()
	historical = app.Flag("historical", "Historical mode").Bool()
	standard   = app.Flag("standard", "Standard mode").Bool()
	device     = app.Flag("device", "Device to read").Required().Short('d').String()

	baudrate = app.Flag("baud", "Baud rate").Short('b').Int()
	size     = app.Flag("size", "Serial frame size").Int()
	parity   = app.Flag("parity", "Serial parity").HintOptions("ParityNone", "N", "ParityOdd", "O", "ParityEven", "E", "ParityMark", "M", "ParitySpace", "S").String()
	stopBits = app.Flag("stopbits", "Serial stopbits").HintOptions("Stop1", "1", "Stop1Half", "15", "Stop2", "2").String()
)

// Linky-exporter command main
func main() {
	// Main action
	app.Action(func(c *kingpin.ParseContext) error { run(); return nil })

	// Parsing
	args, err := app.Parse(os.Args[1:])

	if err != nil {
		log.Error(errors.Wrapf(err, "Error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		kingpin.MustParse(args, err)
	}
}

// Main run function
func run() {
	if debug != nil && *debug {
		log.SetLevel(log.DebugLevel)
		log.Info("Debug mode enabled !")
	}

	// Checks before running
	_, error := os.Stat(*device)
	if error != nil {
		log.Fatal(error)
	}

	// Parse parameters
	connector := core.LinkyConnector{Device: *device}
	detect := auto != nil && *auto
	if !detect {
		if standard != nil && *standard {
			connector.Mode = core.Standard
			connector.BaudRate = core.Standard.BaudRate
			connector.FrameSize = core.Standard.FrameSize
			connector.Parity = core.Standard.Parity
			connector.StopBits = core.Standard.StopBits
		} else if historical != nil && *historical {
			connector.Mode = core.Historical
			connector.BaudRate = core.Historical.BaudRate
			connector.FrameSize = core.Historical.FrameSize
			connector.Parity = core.Historical.Parity
			connector.StopBits = core.Historical.StopBits
		} else {
			detect = true
		}
		if baudrate != nil && *baudrate != 0 {
			connector.BaudRate = *baudrate
		}
		if size != nil && *size != 0 {
			connector.FrameSize = *size
		}
		if parity != nil && *parity != "" {
			log.Debug("Parse parity ", *parity)
			connector.Parity = core.ParseParity(*parity)
		}
		if stopBits != nil && *stopBits != "" {
			log.Debug("Parse Stop Bits ", *stopBits)
			connector.StopBits = core.ParseStopBits(*stopBits)
		}
	}

	// Auto detection mode
	if detect {
		err := connector.Detect()
		log.Debug("device:", connector.Device, " mode:", connector.Mode, " baudrate:", connector.BaudRate, " framesize:", connector.FrameSize, " parity:", connector.Parity, " stopbits:", connector.StopBits)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Run exporter
	exporter := prom.LinkyExporter{Address: *address, Port: *port}
	exporter.Run(connector)
}
