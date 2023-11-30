package main

import (
	"fmt"
	"os"
	"path"

	"github.com/minetest-go/areasparser"
	"github.com/minetest-go/mtdb"
	"github.com/sirupsen/logrus"
)

func ProcessExportProtected(areas []*areasparser.Area) error {
	export_dir := path.Join(wd, "area-export")
	logrus.WithFields(logrus.Fields{
		"export_dir": export_dir,
		"area-count": len(areas),
	}).Info("exporting area-protected chunks")

	if len(areas) == 0 {
		return fmt.Errorf("no areas found, aborting")
	}

	err := os.MkdirAll(export_dir, 0777)
	if err != nil {
		return fmt.Errorf("could not create directory '%s': %v", export_dir, err)
	}

	worldmt_file, err := os.Create(path.Join(export_dir, "world.mt"))
	if err != nil {
		return fmt.Errorf("could not open 'world.mt': %v", err)
	}
	_, err = worldmt_file.WriteString("backend = sqlite3\n")
	if err != nil {
		return fmt.Errorf("could not write to 'world.mt': %v", err)
	}
	worldmt_file.Close()

	export_db, err := mtdb.NewBlockDB(export_dir)
	if err != nil {
		return fmt.Errorf("could not create export-db: %v", err)
	}
	defer export_db.Close()

	exported_chunks := map[string]bool{}
	chunk_count := 0

	for _, area := range areas {
		if area == nil {
			continue
		}

		p1, p2 := SortPos(area.Pos1, area.Pos2)

		chunk1_x, chunk1_y, chunk1_z := GetChunkPosFromNode(p1.X, p1.Y, p1.Z)
		chunk2_x, chunk2_y, chunk2_z := GetChunkPosFromNode(p2.X, p2.Y, p2.Z)

		logrus.WithFields(logrus.Fields{
			"chunk1_x":  chunk1_x,
			"chunk1_y":  chunk1_y,
			"chunk1_z":  chunk1_z,
			"chunk2_x":  chunk2_x,
			"chunk2_y":  chunk2_y,
			"chunk2_z":  chunk2_z,
			"p1":        p1,
			"p2":        p2,
			"area.Pos1": area.Pos1,
			"area.Pos2": area.Pos2,
		}).Info("exported area info")

		for x := chunk1_x; x <= chunk2_x; x++ {
			for y := chunk1_y; y <= chunk2_y; y++ {
				for z := chunk1_z; z <= chunk2_z; z++ {
					// check if already exported
					key := fmt.Sprintf("%d/%d/%d", x, y, z)
					if exported_chunks[key] {
						logrus.WithFields(logrus.Fields{"key": key}).Info("chunk already exported, skipping")
						continue
					}
					// mark as exported
					exported_chunks[key] = true

					logrus.WithFields(logrus.Fields{
						"x": x,
						"y": y,
						"z": z,
					}).Info("exporting chunk")

					chunk_count++
					err = ExportChunk(block_repo, export_db, x, y, z)
					if err != nil {
						return fmt.Errorf("export error in chunk %d/%d/%d: %v", x, y, z, err)
					}
				}
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"chunk_count": chunk_count,
	}).Info("done exporting chunks")

	return nil
}
