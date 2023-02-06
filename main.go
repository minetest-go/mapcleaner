package main

import (
	"errors"
	"os"
	"path"

	"github.com/minetest-go/areasparser"
	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/mtdb"
	"github.com/sirupsen/logrus"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logrus.WithFields(logrus.Fields{"world": wd}).Info("Starting mapcleaner")

	ctx, err := mtdb.New(wd)
	if err != nil {
		panic(err)
	}

	areas_file := path.Join(wd, "areas.json")
	_, err = areasparser.ParseFile(areas_file)
	if err != nil {
		logrus.WithFields(logrus.Fields{"filename": areas_file}).Warn("Areas not found")
	}

	b, err := ctx.Blocks.GetByPos(0, 0, 0)
	if err != nil {
		panic(err)
	}

	if b == nil {
		panic(errors.New("no block found"))
	}

	block, err := mapparser.Parse(b.Data)
	if err != nil {
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"mapping": block.BlockMapping,
	}).Info("Parsed mapblock")
}
