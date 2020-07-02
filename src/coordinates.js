
function get_mapblock_from_pos(pos){
	return {
		x: Math.floor(pos.x / 16),
		y: Math.floor(pos.y / 16),
		z: Math.floor(pos.z / 16)
	};
}

function get_chunkpos_from_pos(pos){
	const mapblock = get_mapblock_from_pos(pos);
	const aligned_mapblock = {
		x: mapblock.x + 2,
		y: mapblock.y + 2,
		z: mapblock.z + 2
	};

	return {
		x: Math.floor(aligned_mapblock.x / 5),
		y: Math.floor(aligned_mapblock.y / 5),
		z: Math.floor(aligned_mapblock.z / 5)
	};
}

function get_mapblocks_from_chunk(chunkpos){
	const minmapblock = {
		x: (chunkpos.x * 5) - 2,
		y: (chunkpos.y * 5) - 2,
		z: (chunkpos.z * 5) - 2,
	};
	const maxmapblock = {
		x: minmapblock.x + 4,
		y: minmapblock.y + 4,
		z: minmapblock.z + 4
	};

	return {
		min: minmapblock,
		max: maxmapblock
	};
}

function get_blocks_from_mapblock(mapblock){
	const min = {
		x: (mapblock.x) * 16,
		y: (mapblock.y) * 16,
		z: (mapblock.z) * 16
	};
	const max = {
		x: min.x + 15,
		y: min.y + 15,
		z: min.z + 15
	};

	return {
		min: min,
		max: max
	};
}

module.exports.get_mapblock_from_pos = get_mapblock_from_pos;
module.exports.get_chunkpos_from_pos = get_chunkpos_from_pos;
module.exports.get_mapblocks_from_chunk = get_mapblocks_from_chunk;
module.exports.get_blocks_from_mapblock = get_blocks_from_mapblock;
