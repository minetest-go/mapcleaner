const executor = require("./executor");
const coordinates = require("./coordinates");

module.exports = function(pos){
	const mapblocks = coordinates.get_mapblocks_from_chunk(pos);
	console.log("removechunk", pos, mapblocks);

	return executor(`
		delete
		from blocks
		where
			posx >= $1 and posy >= $2 and posz >= $3 and
			posx <= $4 and posy <= $5 and posz <= $6
	`, [
		mapblocks.min.x, mapblocks.min.y, mapblocks.min.z,
		mapblocks.max.x, mapblocks.max.y, mapblocks.max.z
	], { single_row: true });
};
