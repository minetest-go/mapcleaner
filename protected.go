package main

var protected_chunks = make(map[string]*bool)

func IsProtected(chunk_x, chunk_y, chunk_z int) (bool, error) {
	key := GetChunkKey(chunk_x, chunk_y, chunk_z)
	p := protected_chunks[key]

	if p == nil {
		x1, y1, z1, x2, y2, z2 := GetMapblockBounds(chunk_x, chunk_y, chunk_z)

		for x := x1; x <= x2; x++ {
			for y := y1; y <= y2; y++ {
				for z := z1; z <= z2; z++ {
					_, err := ctx.Blocks.GetByPos(x, y, z)
					if err != nil {
						return false, err
					}
				}
			}
		}
	}

	return *p, nil
}

func GetMapblockBounds(x, y, z int) (x1, y1, z1, x2, y2, z2 int) {
	x1 = (x * 5) - 2
	y1 = (y * 5) - 2
	z1 = (z * 5) - 2
	x2 = x1 + 4
	y2 = y1 + 4
	z2 = z1 + 4
	return
}
