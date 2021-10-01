
const coordinates = require("./coordinates");
const checkmapblock = require("./checkmapblock");
const areas = require("./areas");

const cache = {};

module.exports = function(pos){
	const str = `${pos.x}/${pos.y}/${pos.z}`;

	if (cache[str]){
		return Promise.resolve(cache[str]);
	}

	if (areas.is_chunk_protected(pos)){
		const result = {
			protected: true,
			generated: true
		};
		cache[str] = result;

		return Promise.resolve(result);
	}

	const mapblocks = coordinates.get_mapblocks_from_chunk(pos);

	const promises = [];
	for (var x=mapblocks.min.x; x<=mapblocks.max.x; x++){
		for (var y=mapblocks.min.y; y<=mapblocks.max.y; y++){
			for (var z=mapblocks.min.z; z<=mapblocks.max.z; z++){
				const coord = { x:x, y:y, z:z };
				promises.push(checkmapblock(coord));
			}
		}
	}

	return Promise.all(promises)
	.then(results => {
		const result = {
			protected: results.some(res => res.protected),
			generated: results.some(res => res.generated)
		};

		cache[str] = result;

		return result;
	});
};
