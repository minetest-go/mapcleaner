package main

import (
	"fmt"
	"time"

	"github.com/minetest-go/mtdb/block"
	"github.com/sirupsen/logrus"
)

func ProcessRemoveUnprotectedBatched() error {
	logrus.Warn("pruning unprotected chunks from the database (batched mode) [EXPERIMENTAL]")

	err := LoadProtectedNodes()
	if err != nil {
		return fmt.Errorf("can't load 'mapcleaner_protect.txt' because of '%v' (i'm refusing to work without that file!)", err)
	}

	err = LoadState()
	if err != nil {
		return err
	}

	if err = InitBatchDB(); err != nil {
		return err
	}

	if err = beginBatchTx(); err != nil {
		return err
	}

	layer, err := LoadYLayer(state.ChunkY)
	if err != nil {
		return err
	}

	protectedChunks, err := buildProtectedChunkSet(layer)
	if err != nil {
		return err
	}

	for {
		if state.ChunkX > state.ToX {
			// commit deletions for this z-stride before saving state
			if err = commitBatchTx(); err != nil {
				return err
			}

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

			if err = beginBatchTx(); err != nil {
				return err
			}
		}
		if state.ChunkZ > state.ToZ {
			// commit before loading the new Y-layer (long I/O operation)
			if err = commitBatchTx(); err != nil {
				return err
			}

			// move to next y stride
			state.ChunkX = state.FromX
			state.ChunkY++
			state.ChunkZ = state.FromZ

			layer, err = LoadYLayer(state.ChunkY)
			if err != nil {
				return err
			}

			protectedChunks, err = buildProtectedChunkSet(layer)
			if err != nil {
				return err
			}

			if err = beginBatchTx(); err != nil {
				return err
			}

			logrus.WithFields(logrus.Fields{
				"chunk_y": state.ChunkY,
			}).Info("Processing next y-layer")
		}
		if state.ChunkY > state.ToY {
			// commit any remaining deletions then save final state
			if err = commitBatchTx(); err != nil {
				return err
			}
			return SaveState()
		}

		logrus.WithFields(logrus.Fields{
			"chunk_x": state.ChunkX,
			"chunk_y": state.ChunkY,
			"chunk_z": state.ChunkZ,
		}).Debug("Processing")

		emerged := batchIsEmerged(layer, state.ChunkX, state.ChunkY, state.ChunkZ)
		if emerged {
			protected := batchIsProtectedWithNeighbors(protectedChunks, state.ChunkX, state.ChunkY, state.ChunkZ)

			if !protected {
				logrus.WithFields(logrus.Fields{
					"chunk_x": state.ChunkX,
					"chunk_y": state.ChunkY,
					"chunk_z": state.ChunkZ,
				}).Info("Removing chunk")

				err = batchRemoveChunk(state.ChunkX, state.ChunkY, state.ChunkZ)
				if err != nil {
					return err
				}

				state.RemovedChunks++
			} else {
				logrus.WithFields(logrus.Fields{
					"chunk_x": state.ChunkX,
					"chunk_y": state.ChunkY,
					"chunk_z": state.ChunkZ,
				}).Info("Retaining chunk")

				state.RetainedChunks++
			}
		}

		state.ProcessedChunks++

		// shift to next chunk
		state.ChunkX++
		time.Sleep(time.Millisecond * time.Duration(state.Delay))
	}
}

func batchBlockKey(x, y, z int) string {
	return fmt.Sprintf("%d/%d/%d", x, y, z)
}

func batchGetBlock(layer map[string]*block.Block, x, y, z int) *block.Block {
	return layer[batchBlockKey(x, y, z)]
}

// batchIsEmerged checks the 8 corners of the chunk using the layer cache.
func batchIsEmerged(layer map[string]*block.Block, chunk_x, chunk_y, chunk_z int) bool {
	x1, y1, z1, x2, y2, z2 := GetMapblockBoundsFromChunk(chunk_x, chunk_y, chunk_z)
	for _, x := range []int{x1, x2} {
		for _, y := range []int{y1, y2} {
			for _, z := range []int{z1, z2} {
				if batchGetBlock(layer, x, y, z) != nil {
					return true
				}
			}
		}
	}
	return false
}

// batchIsProtected checks if the chunk is in the pre-built protected set.
func batchIsProtected(protectedChunks map[string]bool, chunk_x, chunk_y, chunk_z int) bool {
	return protectedChunks[GetChunkKey(chunk_x, chunk_y, chunk_z)]
}

// batchIsProtectedWithNeighbors checks the chunk and all 26 surrounding chunks.
func batchIsProtectedWithNeighbors(protectedChunks map[string]bool, chunk_x, chunk_y, chunk_z int) bool {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				if batchIsProtected(protectedChunks, chunk_x+x, chunk_y+y, chunk_z+z) {
					return true
				}
			}
		}
	}
	return false
}
