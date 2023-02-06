package main

import "github.com/sirupsen/logrus"

func RemoveChunk(chunk_x, chunk_y, chunk_z int) error {
	x1, y1, z1, x2, y2, z2 := GetMapblockBounds(chunk_x, chunk_y, chunk_z)

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			for z := z1; z <= z2; z++ {
				logrus.WithFields(logrus.Fields{
					"x": x,
					"y": y,
					"z": z,
				}).Debug("Removing mapblock")

				err := ctx.Blocks.Delete(x, y, z)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
