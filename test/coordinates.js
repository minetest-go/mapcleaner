const assert = require('assert');
const coordinates = require("../src/coordinates");

describe('coordinates', function() {
  it('returns the proper mapblock for the position', function() {
		let mapblock = coordinates.get_mapblock_from_pos({ x: 1, y: 2, z: 3});
		assert.equal(mapblock.x, 0);
		assert.equal(mapblock.y, 0);
		assert.equal(mapblock.z, 0);

		mapblock = coordinates.get_mapblock_from_pos({ x: 16, y: 2, z: 3});
		assert.equal(mapblock.x, 1);
		assert.equal(mapblock.y, 0);
		assert.equal(mapblock.z, 0);

		mapblock = coordinates.get_mapblock_from_pos({ x: 16, y: 16, z: 3});
		assert.equal(mapblock.x, 1);
		assert.equal(mapblock.y, 1);
		assert.equal(mapblock.z, 0);

		mapblock = coordinates.get_mapblock_from_pos({ x: 16, y: 2, z: 16});
		assert.equal(mapblock.x, 1);
		assert.equal(mapblock.y, 0);
		assert.equal(mapblock.z, 1);

		mapblock = coordinates.get_mapblock_from_pos({ x: -10, y: 2, z: 3});
		assert.equal(mapblock.x, -1);
		assert.equal(mapblock.y, 0);
		assert.equal(mapblock.z, 0);
	});

	it('returns the proper chunk for the position', function() {
		let chunk = coordinates.get_chunkpos_from_pos({ x: 1, y: 2, z: 3});
		assert.equal(chunk.x, 0);
		assert.equal(chunk.y, 0);
		assert.equal(chunk.z, 0);

		chunk = coordinates.get_chunkpos_from_pos({ x: 48, y: 2, z: 3});
		assert.equal(chunk.x, 1);
		assert.equal(chunk.y, 0);
		assert.equal(chunk.z, 0);

		chunk = coordinates.get_chunkpos_from_pos({ x: -32, y: 2, z: 3});
		assert.equal(chunk.x, 0);
		assert.equal(chunk.y, 0);
		assert.equal(chunk.z, 0);

		chunk = coordinates.get_chunkpos_from_pos({ x: -48, y: 2, z: 3});
		assert.equal(chunk.x, -1);
		assert.equal(chunk.y, 0);
		assert.equal(chunk.z, 0);
	});

	it('returns the proper mapblocks for the chunk', function() {
		let mapblocks = coordinates.get_mapblocks_from_chunk({ x: 0, y: 0, z: 0});
		assert.equal(mapblocks.min.x, -2);
		assert.equal(mapblocks.min.y, -2);
		assert.equal(mapblocks.min.z, -2);
		assert.equal(mapblocks.max.x, 2);
		assert.equal(mapblocks.max.y, 2);
		assert.equal(mapblocks.max.z, 2);

		mapblocks = coordinates.get_mapblocks_from_chunk({ x: 0, y: 1, z: 0});
		assert.equal(mapblocks.min.x, -2);
		assert.equal(mapblocks.min.y, 3);
		assert.equal(mapblocks.min.z, -2);
		assert.equal(mapblocks.max.x, 2);
		assert.equal(mapblocks.max.y, 7);
		assert.equal(mapblocks.max.z, 2);

		mapblocks = coordinates.get_mapblocks_from_chunk({ x: -1, y: 1, z: 0});
		assert.equal(mapblocks.min.x, -7);
		assert.equal(mapblocks.min.y, 3);
		assert.equal(mapblocks.min.z, -2);
		assert.equal(mapblocks.max.x, -3);
		assert.equal(mapblocks.max.y, 7);
		assert.equal(mapblocks.max.z, 2);
	});

	it('returns the proper blocks for the mapblock', function() {
		let blocks = coordinates.get_blocks_from_mapblock({ x: 0, y: 0, z: 0});
		assert.equal(blocks.min.x, 0);
		assert.equal(blocks.min.y, 0);
		assert.equal(blocks.min.z, 0);
		assert.equal(blocks.max.x, 15);
		assert.equal(blocks.max.y, 15);
		assert.equal(blocks.max.z, 15);

		blocks = coordinates.get_blocks_from_mapblock({ x: 0, y: 0, z: 1});
		assert.equal(blocks.min.x, 0);
		assert.equal(blocks.min.y, 0);
		assert.equal(blocks.min.z, 16);
		assert.equal(blocks.max.x, 15);
		assert.equal(blocks.max.y, 15);
		assert.equal(blocks.max.z, 31);

		blocks = coordinates.get_blocks_from_mapblock({ x: 0, y: 0, z: -1});
		assert.equal(blocks.min.x, 0);
		assert.equal(blocks.min.y, 0);
		assert.equal(blocks.min.z, -16);
		assert.equal(blocks.max.x, 15);
		assert.equal(blocks.max.y, 15);
		assert.equal(blocks.max.z, -1);
	});
});
