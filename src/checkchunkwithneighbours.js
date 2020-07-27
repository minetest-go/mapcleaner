const checkchunk = require("./checkchunk");

module.exports = function(pos){
	const promises = [];
	for (var x=pos.x-1; x<=pos.x+1; x++){
		for (var y=pos.y-1; y<=pos.y+1; y++){
			for (var z=pos.z-1; z<=pos.z+1; z++){
				const coord = { x:x, y:y, z:z };
				promises.push(checkchunk(coord));
			}
		}
	}

	return Promise.all(promises)
	.then(results => {
		return {
			protected: results.some(res => res.protected),
			generated: results.some(res => res.generated)
		};
	});
};
