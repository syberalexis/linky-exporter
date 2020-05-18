package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/syberalexis/linky-exporter/pkg/core"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

var (
	defaultPort = 9901
)

func main() {
	exporter := &core.LinkyExporter{
		Port: defaultPort,
	}

	// Globals
	app := kingpin.New(filepath.Base(os.Args[0]), "")
	app.HelpFlag.Short('h')
	app.Version("0.0.1")
	app.Action(func(c *kingpin.ParseContext) error { exporter.Run(); return nil })

	// Flags
	app.Flag("port", "Listen port").Default(fmt.Sprintf("%d", defaultPort)).Short('p').IntVar(&exporter.Port)


	// Parsing
	args, err := app.Parse(os.Args[1:])
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		kingpin.MustParse(args, err)
	}
}
