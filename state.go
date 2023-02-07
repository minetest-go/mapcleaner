package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type State struct {
	ChunkX          int `json:"chunk_x"`
	ChunkY          int `json:"chunk_y"`
	ChunkZ          int `json:"chunk_z"`
	RemovedChunks   int `json:"removed_chunks"`
	RetainedChunks  int `json:"retained_chunks"`
	ProcessedChunks int `json:"processed_chunks"`
	FromX           int `json:"from_x"`
	FromZ           int `json:"from_z"`
	ToX             int `json:"to_x"`
	ToY             int `json:"to_y"`
	ToZ             int `json:"to_z"`
	Delay           int `json:"delay"`
}

const filename = "mapcleaner.json"

// current state
var state = &State{
	ChunkX: -400,
	ChunkY: -400,
	ChunkZ: -400,
	FromX:  -400,
	FromZ:  -400,
	ToX:    400,
	ToY:    400,
	ToZ:    400,
}

func LoadState() error {
	file_path := path.Join(wd, filename)
	info, err := os.Stat(file_path)
	if info == nil || !info.Mode().IsRegular() {
		logrus.Warn("Creating new state config (this is normal on the first run)")
		return SaveState()
	}

	if err != nil {
		return err
	}

	data, err := os.ReadFile(file_path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, state)
}

func SaveState() error {
	file_path := path.Join(wd, filename)

	f, err := os.OpenFile(file_path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(state)
}
