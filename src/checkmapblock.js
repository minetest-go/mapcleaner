
const mapblockparser = require("./mapblockparser");
const executor = require("./executor");


module.exports = function(pos){
	return executor(`
		select data
		from blocks
		where posx = $1 and posy = $2 and posz = $3
	`, [pos.x,pos.y,pos.z], { single_row: true })
	.then(block => {
		console.log(pos, block);
		if (block)
			return mapblockparser.parse(block.data);
		else
			return;
	})
	.then(mapblock => {
		console.log(pos, mapblock);
		//TODO: magic!
		return {
			protected: true,
			generated: true
		};
	});
};
