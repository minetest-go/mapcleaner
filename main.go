package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/minetest-go/areasparser"
	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/block"
	"github.com/sirupsen/logrus"
)

var Version string
var block_repo block.BlockRepository
var wd string

func main() {

	help := flag.Bool("help", false, "show help")
	debug := flag.Bool("debug", false, "set the loglevel to debug")
	mode := flag.String("mode", "prune_unproteced", "set the working mode [prune_unproteced|export_protected]")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("loglevel set to debug")
	}

	var err error
	wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	if Version == "" {
		Version = "DEV"
	}

	logrus.WithFields(logrus.Fields{
		"world":   wd,
		"version": Version,
	}).Info("Starting mapcleaner")

	// load main block database
	block_repo, err = mtdb.NewBlockDB(wd)
	if err != nil {
		panic(err)
	}

	// load areas file
	areas_file := path.Join(wd, "areas.dat")
	areas, err := areasparser.ParseFile(areas_file)
	if err != nil {
		logrus.WithFields(logrus.Fields{"filename": areas_file}).Warn("Areas not found")
	} else {
		for _, area := range areas {
			if area == nil {
				continue
			}
			PopulateAreaProtection(area)
		}
	}

	// start main process
	switch *mode {
	case "prune_unproteced":
		err = ProcessRemoveUnprotected()
	case "export_protected":
		err = ProcessExportProtected(areas)
	default:
		panic(fmt.Sprintf("mode not implemented: '%s'", *mode))
	}
	if err != nil {
		panic(err)
	}

	logrus.Info("mapcleaner exiting")
}
