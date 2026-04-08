package main

import (
	"database/sql"
	"fmt"
	"path"

	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/mtdb/worldconfig"
	"github.com/sirupsen/logrus"
)

var batchDB *sql.DB
var batchDBType string

// InitBatchDB opens a direct DB connection used for bulk Y-layer prefetch queries.
// Must be called before LoadYLayer.
func InitBatchDB() error {
	wc, err := worldconfig.Parse(path.Join(wd, "world.mt"))
	if err != nil {
		return err
	}

	batchDBType = wc[worldconfig.CONFIG_MAP_BACKEND]

	switch batchDBType {
	case worldconfig.BACKEND_POSTGRES:
		batchDB, err = sql.Open("postgres", wc[worldconfig.CONFIG_PSQL_MAP_CONNECTION])
		if err != nil {
			return err
		}
		logrus.Info("Batch DB initialized (PostgreSQL)")

	case worldconfig.BACKEND_SQLITE3:
		batchDB, err = sql.Open("sqlite3", path.Join(wd, "map.sqlite"))
		if err != nil {
			return err
		}
		hasPosCol, err := sqliteHasPosColumn(batchDB)
		if err != nil {
			batchDB.Close()
			batchDB = nil
			return err
		}
		if hasPosCol {
			return fmt.Errorf("SQLite legacy pos-column format is not supported by batched mode; use prune_unprotected instead")
		}
		logrus.Info("Batch DB initialized (SQLite)")

	default:
		return fmt.Errorf("unsupported backend '%s' for batched mode", batchDBType)
	}

	return nil
}

func sqliteHasPosColumn(db *sql.DB) (bool, error) {
	rows, err := db.Query("pragma table_info(blocks)")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}

	nameIdx := -1
	for i, c := range cols {
		if c == "name" {
			nameIdx = i
			break
		}
	}
	if nameIdx == -1 {
		return false, nil
	}

	vals := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return false, err
		}
		switch v := vals[nameIdx].(type) {
		case []byte:
			if string(v) == "pos" {
				return true, nil
			}
		case string:
			if v == "pos" {
				return true, nil
			}
		}
	}

	return false, rows.Err()
}

// LoadYLayer fetches all mapblocks for chunk_y ±1 (3 chunk-Y layers = 15 mapblock rows)
// across the full configured X/Z scan area plus ±1 chunk neighbors, in a single query.
// The returned map is keyed by "posX/posY/posZ" in mapblock coordinates.
func LoadYLayer(chunk_y int) (map[string]*block.Block, error) {
	minY := (chunk_y - 1) * 5
	maxY := (chunk_y+1)*5 + 4

	minX := (state.FromX - 1) * 5
	maxX := (state.ToX+1)*5 + 4
	minZ := (state.FromZ - 1) * 5
	maxZ := (state.ToZ+1)*5 + 4

	logrus.WithFields(logrus.Fields{
		"chunk_y":          chunk_y,
		"mapblock_y_range": fmt.Sprintf("[%d, %d]", minY, maxY),
	}).Info("Prefetching Y-layer block cache")

	var rows *sql.Rows
	var err error

	switch batchDBType {
	case worldconfig.BACKEND_POSTGRES:
		rows, err = batchDB.Query(
			`SELECT posX, posY, posZ, data FROM blocks
			 WHERE posX >= $1 AND posX <= $2
			   AND posY >= $3 AND posY <= $4
			   AND posZ >= $5 AND posZ <= $6`,
			minX, maxX, minY, maxY, minZ, maxZ)
	case worldconfig.BACKEND_SQLITE3:
		rows, err = batchDB.Query(
			`SELECT x, y, z, data FROM blocks
			 WHERE x >= ? AND x <= ?
			   AND y >= ? AND y <= ?
			   AND z >= ? AND z <= ?`,
			minX, maxX, minY, maxY, minZ, maxZ)
	default:
		return nil, fmt.Errorf("unsupported backend '%s'", batchDBType)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	layer := make(map[string]*block.Block)
	for rows.Next() {
		b := &block.Block{}
		if err := rows.Scan(&b.PosX, &b.PosY, &b.PosZ, &b.Data); err != nil {
			return nil, err
		}
		layer[fmt.Sprintf("%d/%d/%d", b.PosX, b.PosY, b.PosZ)] = b
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"chunk_y":       chunk_y,
		"blocks_loaded": len(layer),
	}).Info("Y-layer block cache loaded")

	return layer, nil
}

// buildProtectedChunkSet does a single pass over the loaded layer, parsing each
// mapblock once, and returns a set of chunk keys that contain protected nodes or
// are covered by areas protection. The result is valid for the lifetime of the
// current Y-layer and should be rebuilt whenever LoadYLayer is called.
func buildProtectedChunkSet(layer map[string]*block.Block) (map[string]bool, error) {
	protected := make(map[string]bool)

	// seed with areas-protected chunks first (no parsing needed)
	for key, v := range protected_areas {
		if v {
			protected[key] = true
		}
	}

	for _, mb := range layer {
		chunk_x, chunk_y, chunk_z := GetChunkPosFromMapblock(mb.PosX, mb.PosY, mb.PosZ)
		key := GetChunkKey(chunk_x, chunk_y, chunk_z)
		if protected[key] {
			// already known protected, skip parsing
			continue
		}

		b, err := mapparser.Parse(mb.Data)
		if err != nil {
			return nil, err
		}
		for _, name := range b.BlockMapping {
			if protected_nodenames[name] {
				protected[key] = true
				break
			}
		}
	}

	return protected, nil
}

// batchRemoveChunk deletes all 125 mapblocks of a chunk in a single range query.
func batchRemoveChunk(chunk_x, chunk_y, chunk_z int) error {
	x1, y1, z1, x2, y2, z2 := GetMapblockBoundsFromChunk(chunk_x, chunk_y, chunk_z)

	logrus.WithFields(logrus.Fields{
		"chunk_x": chunk_x,
		"chunk_y": chunk_y,
		"chunk_z": chunk_z,
	}).Debug("Bulk removing chunk mapblocks")

	switch batchDBType {
	case worldconfig.BACKEND_POSTGRES:
		_, err := batchDB.Exec(
			`DELETE FROM blocks
			 WHERE posX >= $1 AND posX <= $2
			   AND posY >= $3 AND posY <= $4
			   AND posZ >= $5 AND posZ <= $6`,
			x1, x2, y1, y2, z1, z2)
		return err
	case worldconfig.BACKEND_SQLITE3:
		_, err := batchDB.Exec(
			`DELETE FROM blocks
			 WHERE x >= ? AND x <= ?
			   AND y >= ? AND y <= ?
			   AND z >= ? AND z <= ?`,
			x1, x2, y1, y2, z1, z2)
		return err
	default:
		return fmt.Errorf("unsupported backend '%s'", batchDBType)
	}
}
