package main

import (
	"bufio"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/minetest-go/areasparser"
	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/mtdb/block"
	"github.com/sirupsen/logrus"
)

var protected_nodenames = make(map[string]bool)
var protected_areas = make(map[string]bool)

// caches
var protected_chunks = expirable.NewLRU[string, bool](1000, nil, time.Minute*10)
var emerged_chunks = expirable.NewLRU[string, bool](1000, nil, time.Minute*10)

func PopulateAreaProtection(area *areasparser.Area) {
	logrus.WithFields(logrus.Fields{
		"pos1":  area.Pos1,
		"pos2":  area.Pos2,
		"name":  area.Name,
		"owner": area.Owner,
	}).Info("Adding area protection")

	x1, y1, z1 := GetChunkPosFromNode(area.Pos1.X, area.Pos1.Y, area.Pos1.Z)
	x2, y2, z2 := GetChunkPosFromNode(area.Pos2.X, area.Pos2.Y, area.Pos2.Z)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			for z := z1; z <= z2; z++ {
				key := GetChunkKey(x, y, z)
				protected_areas[key] = true
			}
		}
	}
}

func LoadProtectedNodes() error {
	file, err := os.Open(path.Join(wd, "mapcleaner_protect.txt"))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			// comment or empty line
			continue
		}

		protected_nodenames[line] = true
		logrus.WithFields(logrus.Fields{
			"nodename": line,
		}).Info("Adding nodename to protected list")
	}

	return nil
}

// check all 8 corners of the chunk for existing mapblocks
func IsEmerged(chunk_x, chunk_y, chunk_z int) (bool, error) {
	key := GetChunkKey(chunk_x, chunk_y, chunk_z)
	e, ok := emerged_chunks.Get(key)
	if ok {
		return e, nil
	}

	// check if first mapblock exists
	x1, y1, z1, x2, y2, z2 := GetMapblockBoundsFromChunk(chunk_x, chunk_y, chunk_z)

	emerged := false
	for _, x := range []int{x1, x2} {
		for _, y := range []int{y1, y2} {
			for _, z := range []int{z1, z2} {
				data, err := block_repo.GetByPos(x, y, z)
				if err != nil {
					return false, err
				}

				if data != nil {
					// emerged
					emerged = true
					break
				}
			}
			if emerged {
				break
			}
		}
		if emerged {
			break
		}
	}

	// cache for next time
	emerged_chunks.Add(key, emerged)
	return emerged, nil
}

func IsProtected(chunk_x, chunk_y, chunk_z int) (bool, error) {
	key := GetChunkKey(chunk_x, chunk_y, chunk_z)

	// check area protection first
	if protected_areas[key] {
		return true, nil
	}

	p, ok := protected_chunks.Get(key)
	if ok {
		return p, nil
	}

	x1, y1, z1, x2, y2, z2 := GetMapblockBoundsFromChunk(chunk_x, chunk_y, chunk_z)

	protected := false
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			for z := z1; z <= z2; z++ {
				logrus.WithFields(logrus.Fields{
					"x": x,
					"y": y,
					"z": z,
				}).Debug("Checking mapblock")

				block, err := block_repo.GetByPos(x, y, z)
				if err != nil {
					return false, err
				}

				if block == nil {
					// no block here
					continue
				}

				b, err := mapparser.Parse(block.Data)
				if err != nil {
					return false, err
				}

				for _, name := range b.BlockMapping {
					if protected_nodenames[name] {
						// protected block here
						protected = true
						break
					}
				}

				if protected {
					break
				}
			}
			if protected {
				break
			}
		}
		if protected {
			break
		}
	}

	protected_chunks.Add(key, protected)
	return protected, nil
}

// check chunk with surroundings
func IsProtectedWithNeighbors(chunk_x, chunk_y, chunk_z int) (bool, error) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				p, err := IsProtected(chunk_x+x, chunk_y+y, chunk_z+z)
				if err != nil {
					return false, err
				}
				if p {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// IsBlockProtected checks if the block is protected either by areas mod or by the
// user provided nodes.
func IsBlockProtected(b *block.Block) (protected bool, err error) {
	// Is this block protected by areas mod?
	key := GetChunkKey(GetChunkPosFromMapblock(b.PosX, b.PosY, b.PosZ))
	protected = protected_areas[key]
	if protected {
		return true, nil
	}

	// Is this block protected by the user selected nodes?
	parsedBlock, err := mapparser.Parse(b.Data)
	if err != nil {
		return false, err
	}
	for _, name := range parsedBlock.BlockMapping {
		if protected_nodenames[name] {
			logrus.WithFields(logrus.Fields{"block_mapping": parsedBlock.BlockMapping, "found": name}).Debug("block is protected")
			// protected block here
			protected = true
			break
		}
	}
	return protected, nil
}
