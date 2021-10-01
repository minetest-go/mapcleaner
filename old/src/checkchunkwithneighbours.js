const checkchunk = require("./checkchunk");

module.exports = function(pos){
	const promises = [];
	let centralchunkindex;

	for (var x=pos.x-1; x<=pos.x+1; x++){
		for (var y=pos.y-1; y<=pos.y+1; y++){
			for (var z=pos.z-1; z<=pos.z+1; z++){
				const coord = { x:x, y:y, z:z };
				promises.push(checkchunk(coord));

				if (x == pos.x && y == pos.y && z == pos.z){
					// central chunk, remember for later
					centralchunkindex = promises.length - 1;
				}
			}
		}
	}

	return Promise.all(promises)
	.then(results => {
		return {
			// mark as protected if *any* surrounding chunk is protected
			protected: results.some(res => res.protected),
			// return generated property of central chunk
			generated: results[centralchunkindex].generated
		};
	});
};
