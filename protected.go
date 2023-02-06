package main

import (
	"bufio"
	"os"
	"path"
	"strings"

	"github.com/minetest-go/areasparser"
	"github.com/minetest-go/mapparser"
	"github.com/sirupsen/logrus"
)

var protected_nodenames = make(map[string]bool)
var protected_areas = make(map[string]bool)
var protected_chunks = make(map[string]*bool)

func ClearProtectionCache() {
	protected_chunks = make(map[string]*bool)
}

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

func IsEmerged(chunk_x, chunk_y, chunk_z int) (bool, error) {
	// check if first mapblock exists
	x1, y1, z1, _, _, _ := GetMapblockBounds(chunk_x, chunk_y, chunk_z)
	data, err := ctx.Blocks.GetByPos(x1, y1, z1)

	// mark chunk as unprotected in case of neighbor check
	protected := false
	protected_chunks[GetChunkKey(chunk_x, chunk_y, chunk_z)] = &protected
	return data != nil, err
}

func IsProtected(chunk_x, chunk_y, chunk_z int) (bool, error) {
	key := GetChunkKey(chunk_x, chunk_y, chunk_z)

	// check area protection first
	if protected_areas[key] {
		return true, nil
	}

	p := protected_chunks[key]

	if p == nil {
		x1, y1, z1, x2, y2, z2 := GetMapblockBounds(chunk_x, chunk_y, chunk_z)

		protected := false
		for x := x1; x <= x2; x++ {
			for y := y1; y <= y2; y++ {
				for z := z1; z <= z2; z++ {
					logrus.WithFields(logrus.Fields{
						"x": x,
						"y": y,
						"z": z,
					}).Debug("Checking mapblock")

					block, err := ctx.Blocks.GetByPos(x, y, z)
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

		p = &protected
		protected_chunks[key] = p
	}

	return *p, nil
}

func GetMapblockBounds(x, y, z int) (x1, y1, z1, x2, y2, z2 int) {
	x1 = (x * 5) - 2
	y1 = (y * 5) - 2
	z1 = (z * 5) - 2
	x2 = x1 + 4
	y2 = y1 + 4
	z2 = z1 + 4
	return
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
