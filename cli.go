package main

import (
	"fmt"
	"io"

	"github.com/alecthomas/kingpin/v2"
)

const (
	appName        = "tinygo-flash"
	appDescription = ""
)

type cli struct {
	outStream io.Writer
	errStream io.Writer
}

var (
	app    = kingpin.New(appName, appDescription)
	port   = app.Flag("port", "COM port").String()
	target = app.Flag("target", "target").Required().Enum(`pyportal`, `feather-m4`, `trinket-m0`)
	uf2    = app.Arg("uf2", "*.uf2").Required().ExistingFile()
)

// Run ...
func (c *cli) Run(args []string) error {
	app.UsageWriter(c.errStream)

	if VERSION != "" {
		app.Version(fmt.Sprintf("%s version %s build %s", appName, VERSION, BUILDDATE))
	} else {
		app.Version(fmt.Sprintf("%s version - build -", appName))
	}
	app.HelpFlag.Short('h')

	k, err := app.Parse(args[1:])
	if err != nil {
		return err
	}

	switch k {
	default:
		if len(*port) == 0 {
			p, err := getDefaultPort()
			if err != nil {
				return err
			}
			//fmt.Println(p)
			*port = p
		}

		err := flash(*port, *target, *uf2)
		if err != nil {
			return err
		}
	}

	return nil
}
