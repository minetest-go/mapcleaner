package main

import (
	"github.com/minetest-go/mapparser"
	"github.com/sirupsen/logrus"
)

func Process() error {
	for {
		if state.ChunkX > 400 {
			// move to next z stride
			state.ChunkX = -400
			state.ChunkZ++

			logrus.WithFields(logrus.Fields{
				"chunk_y": state.ChunkY,
				"chunk_z": state.ChunkZ,
			}).Info("Processing next z-stride")

			err := SaveState()
			if err != nil {
				return err
			}
		}
		if state.ChunkZ > 400 {
			// move to next y stride
			state.ChunkX = -400
			state.ChunkZ = -400
			state.ChunkZ++

			logrus.WithFields(logrus.Fields{
				"chunk_y": state.ChunkY,
			}).Info("Processing next y-layer")
		}
		if state.ChunkY > 400 {
			// done
			return nil
		}

		logrus.WithFields(logrus.Fields{
			"chunk_x": state.ChunkX,
			"chunk_y": state.ChunkY,
			"chunk_z": state.ChunkZ,
		}).Debug("Processing")

		b, err := ctx.Blocks.GetByPos(0, 0, 0)
		if err != nil {
			return err
		}

		if b != nil {
			block, err := mapparser.Parse(b.Data)
			if err != nil {
				return err
			}

			logrus.WithFields(logrus.Fields{
				"mapping": block.BlockMapping,
			}).Debug("Parsed mapblock")
		}

		// shift to next chunk
		state.ChunkX++
	}
}
