package main

import (
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
	var err error
	wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	err = LoadProtectedNodes()
	if err != nil {
		panic(fmt.Errorf("can't load 'mapcleaner_protect.txt' because of '%v' (i'm refusing to work without that file!)", err))
	}

	if Version == "" {
		Version = "DEV"
	}

	logrus.WithFields(logrus.Fields{
		"world":   wd,
		"version": Version,
	}).Info("Starting mapcleaner")

	block_repo, err = mtdb.NewBlockDB(wd)
	if err != nil {
		panic(err)
	}

	areas_file := path.Join(wd, "areas.dat")
	areas, err := areasparser.ParseFile(areas_file)
	if err != nil {
		logrus.WithFields(logrus.Fields{"filename": areas_file}).Warn("Areas not found")
	} else {
		for _, area := range areas {
			PopulateAreaProtection(area)
		}
	}

	err = LoadState()
	if err != nil {
		panic(err)
	}

	err = Process()
	if err != nil {
		panic(err)
	}

	logrus.Info("Finished mapcleaner run")
}
