package main

import "fmt"

func GetChunkKey(x, y, z int) string {
	return fmt.Sprintf("%d/%d/%d", x, y, z)
}

func GetMapblockPosFromNode(x, y, z int) (int, int, int) {
	return int(x / 16), int(y / 16), int(z / 16)
}

func GetChunkPosFromMapblock(x, y, z int) (int, int, int) {
	return int((x + 2) / 5), int((y + 2) / 5), int((z + 2) / 5)
}

func GetChunkPosFromNode(x, y, z int) (int, int, int) {
	m_x, m_y, m_z := GetMapblockPosFromNode(x, y, z)
	return GetChunkPosFromMapblock(m_x, m_y, m_z)
}
