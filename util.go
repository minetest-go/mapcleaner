package main

import "fmt"

func GetChunkKey(x, y, z int) string {
	return fmt.Sprintf("%d/%d/%d", x, y, z)
}
