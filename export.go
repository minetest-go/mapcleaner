package main

import (
	"fmt"

	"github.com/minetest-go/mtdb/block"
)

func ExportChunk(src, dst block.BlockRepository, x, y, z int) error {

	x1, y1, z1, x2, y2, z2 := GetMapblockBoundsFromChunk(x, y, z)
	for mbx := x1; mbx <= x2; mbx++ {
		for mby := y1; mby <= y2; mby++ {
			for mbz := z1; mbz <= z2; mbz++ {
				block, err := src.GetByPos(mbx, mby, mbz)
				if err != nil {
					return fmt.Errorf("error in src-mapblock %d/%d/%d: %v", mbx, mby, mbz, err)
				}
				if block == nil {
					continue
				}
				err = dst.Update(block)
				if err != nil {
					return fmt.Errorf("error in dst-mapblock %d/%d/%d: %v", mbx, mby, mbz, err)
				}
			}
		}
	}

	return nil
}
