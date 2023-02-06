package main

import (
	"fmt"
	"os"
	"path"

	"github.com/minetest-go/areasparser"
	"github.com/minetest-go/mtdb"
	"github.com/sirupsen/logrus"
)

var ctx *mtdb.Context
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

	logrus.WithFields(logrus.Fields{"world": wd}).Info("Starting mapcleaner")

	ctx, err = mtdb.New(wd)
	if err != nil {
		panic(err)
	}

	areas_file := path.Join(wd, "areas.json")
	_, err = areasparser.ParseFile(areas_file)
	if err != nil {
		logrus.WithFields(logrus.Fields{"filename": areas_file}).Warn("Areas not found")
	}
	//TODO: area protection

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
