package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/minetest-go/areasparser"
	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/block"
	"github.com/sirupsen/logrus"
)

// initializeExportDirectory initializes the output directory for exporting.
func initializeExportDirectory(export_dir string) (export_db block.BlockRepository, err error) {
	err = os.MkdirAll(export_dir, 0777)
	if err != nil {
		return nil, fmt.Errorf("could not create export-db: %v", err)
	}

	worldmt_file, err := os.Create(path.Join(export_dir, "world.mt"))
	if err != nil {
		return nil, fmt.Errorf("could not create export-db: %v", err)
	}
	_, err = worldmt_file.WriteString("backend = sqlite3\n")
	if err != nil {
		return nil, fmt.Errorf("could not create export-db: %v", err)
	}
	worldmt_file.Close()

	export_db, err = mtdb.NewBlockDB(export_dir)
	if err != nil {
		return nil, fmt.Errorf("could not create export-db: %v", err)
	}
	return export_db, nil
}

func ProcessExportProtected(areas []*areasparser.Area) error {
	if len(areas) == 0 {
		return fmt.Errorf("no areas found, aborting")
	}

	export_dir := path.Join(wd, "area-export")

	logrus.WithFields(logrus.Fields{
		"export_dir": export_dir,
		"area-count": len(areas),
	}).Info("exporting area-protected chunks")

	export_db, err := initializeExportDirectory(export_dir)
	if err != nil {
		return err
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

func ProccessExportAllProtected() error {
	// syntax suggar
	type f logrus.Fields
	I := func(msg string, f f) { logrus.WithFields(logrus.Fields(f)).Info(msg) }
	D := func(msg string, f f) { logrus.WithFields(logrus.Fields(f)).Debug(msg) }
	W := func(msg string, f f) { logrus.WithFields(logrus.Fields(f)).Warning(msg) }

	// Load protected nodes
	if err := LoadProtectedNodes(); err != nil {
		return err
	}

	export_dir := path.Join(wd, "area-export")
	export_info := f{
		"export_dir":      export_dir,
		"area_count":      len(protected_areas),
		"export_all":      true,
		"protected_nodes": protected_nodenames,
	}
	I("exporting all protected areas; parsing the entire map", export_info)

	export_db, err := initializeExportDirectory(export_dir)
	if err != nil {
		return err
	}
	defer export_db.Close()
	D("export_db initialized", f{"export_dir": export_dir})

	// Check source size
	total_blocks, err := block_repo.Count()
	if err != nil {
		return err
	}

	// Initialize variables
	D("initializing iterator from source", f{"total_blocks": total_blocks})
	it, closer, err := block_repo.Iterator(block.AsBlockPos(-33000, -33000, -33000))
	if err != nil {
		return err
	}
	defer closer.Close()

	// Stats
	exported_chunks := map[string]bool{}
	chunk_count := 0
	block_count := 0
	start := time.Now()
	D("parsing all blocks", f{})
	for b := range it {
		block_count++

		protected, err := IsBlockProtected(b)
		if err != nil {
			W("error detecting if block is protected", f{"block": b})
			return err
		}

		if protected {
			chunkx, chunky, chunkz := GetChunkPosFromMapblock(b.PosX, b.PosY, b.PosZ)

			// Export surounding chunks as protected elements can be inside a construction that is
			// at the edge of a mapblock/mapchunk.
			for x := chunkx - 1; x <= chunkx+1; x++ {
				for y := chunky - 1; y <= chunky+1; y++ {
					for z := chunkz - 1; z <= chunkz+1; z++ {
						key := GetChunkKey(x, y, z)
						if exported_chunks[key] {
							D("chunk already exported, skipping", f{"key": key})
							continue
						}

						I("exporting chunk", f{"x": x, "y": y, "z": z})
						err = ExportChunk(block_repo, export_db, x, y, z)
						if err != nil {
							return fmt.Errorf("export error in chunk %d/%d/%d: %v", x, y, z, err)
						}

						// mark as exported
						exported_chunks[key] = true
						chunk_count++
					}
				}
			}
		}

		// Report progress every 100 blocks
		if block_count%1000 == 0 {
			progress := 100 * float64(block_count) / float64(total_blocks)
			stats := f{
				"expo_chunks": chunk_count,
				"proc_blocks": block_count,
				"progress":    fmt.Sprintf("%.02f%%", progress),
				"elapsed":     time.Since(start).String(),
			}
			I("processing blocks", stats)
		}
	}

	I("export finished", f{"exported_chunks": chunk_count, "processed_blocks": block_count})
	return nil
}
