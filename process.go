package main

import (
	"github.com/sirupsen/logrus"
)

func Process() error {
	for {
		if state.ChunkX > state.ToX {
			// move to next z stride
			state.ChunkX = state.FromX
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
		if state.ChunkZ > state.ToZ {
			// move to next y stride
			state.ChunkX = state.FromX
			state.ChunkY++
			state.ChunkZ = state.FromZ

			// purge cache after each layer
			ClearCache()

			logrus.WithFields(logrus.Fields{
				"chunk_y": state.ChunkY,
			}).Info("Processing next y-layer")
		}
		if state.ChunkY > state.ToY {
			// done
			return SaveState()
		}

		logrus.WithFields(logrus.Fields{
			"chunk_x": state.ChunkX,
			"chunk_y": state.ChunkY,
			"chunk_z": state.ChunkZ,
		}).Debug("Processing")

		emerged, err := IsEmerged(state.ChunkX, state.ChunkY, state.ChunkZ)
		if err != nil {
			return err
		}
		if emerged {
			protected, err := IsProtectedWithNeighbors(state.ChunkX, state.ChunkY, state.ChunkZ)
			if err != nil {
				return err
			}

			if !protected {
				logrus.WithFields(logrus.Fields{
					"chunk_x": state.ChunkX,
					"chunk_y": state.ChunkY,
					"chunk_z": state.ChunkZ,
				}).Info("Removing chunk")

				err = RemoveChunk(state.ChunkX, state.ChunkY, state.ChunkZ)
				if err != nil {
					return err
				}

				state.RemovedChunks++
			} else {
				state.RetainedChunks++
			}
		}

		state.ProcessedChunks++

		// shift to next chunk
		state.ChunkX++
	}
}
