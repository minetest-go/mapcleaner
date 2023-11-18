package main

import (
	"fmt"

	"github.com/minetest-go/areasparser"
)

func GetChunkKey(x, y, z int) string {
	return fmt.Sprintf("%d/%d/%d", x, y, z)
}

func GetMapblockPosFromNode(x, y, z int) (int, int, int) {
	return int(x / 16), int(y / 16), int(z / 16)
}

func GetMapblockBoundsFromChunk(x, y, z int) (x1, y1, z1, x2, y2, z2 int) {
	x1 = (x * 5) - 2
	y1 = (y * 5) - 2
	z1 = (z * 5) - 2
	x2 = x1 + 4
	y2 = y1 + 4
	z2 = z1 + 4
	return
}

func GetChunkPosFromMapblock(x, y, z int) (int, int, int) {
	return int((x + 2) / 5), int((y + 2) / 5), int((z + 2) / 5)
}

func GetChunkPosFromNode(x, y, z int) (int, int, int) {
	m_x, m_y, m_z := GetMapblockPosFromNode(x, y, z)
	return GetChunkPosFromMapblock(m_x, m_y, m_z)
}

func SortPos(p1, p2 *areasparser.GenericPos) (*areasparser.GenericPos, *areasparser.GenericPos) {
	lo := &areasparser.GenericPos{}
	hi := &areasparser.GenericPos{}

	if p1.X > p2.X {
		hi.X = p1.X
		lo.X = p2.X
	} else {
		hi.X = p2.X
		lo.X = p1.X
	}

	if p1.Y > p2.Y {
		hi.Y = p1.Y
		lo.Y = p2.Y
	} else {
		hi.Y = p2.Y
		lo.Y = p1.Y
	}

	if p1.Z > p2.Z {
		hi.Z = p1.Z
		lo.Z = p2.Z
	} else {
		hi.Z = p2.Z
		lo.Z = p1.Z
	}

	return lo, hi
}
