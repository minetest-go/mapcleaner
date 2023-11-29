package main

import (
	"fmt"
	"math"

	"github.com/minetest-go/areasparser"
)

func GetChunkKey(x, y, z int) string {
	return fmt.Sprintf("%d/%d/%d", x, y, z)
}

func floorOfDivision(n, q int) int {
	return int(math.Floor(float64(n) / float64(q)))
}

func GetMapblockPosFromNode(x, y, z int) (int, int, int) {
	return floorOfDivision(x, 16), floorOfDivision(y, 16), floorOfDivision(z, 16)
}

func GetMapblockBoundsFromChunk(x, y, z int) (x1, y1, z1, x2, y2, z2 int) {
	x1 = (x * 5)
	y1 = (y * 5)
	z1 = (z * 5)
	x2 = x1 + 4
	y2 = y1 + 4
	z2 = z1 + 4
	return
}

func GetChunkPosFromMapblock(x, y, z int) (int, int, int) {
	return floorOfDivision((x), 5), floorOfDivision((y), 5), floorOfDivision((z), 5)
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
