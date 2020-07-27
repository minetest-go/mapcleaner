
const mapblockparser = require("./mapblockparser");
const executor = require("./executor");
const protected = require("./protected");

module.exports = function(pos){
	return executor(`
		select data
		from blocks
		where posx = $1 and posy = $2 and posz = $3
	`, [pos.x,pos.y,pos.z], { single_row: true })
	.then(block => {
		if (block)
			return mapblockparser.parse(block.data);
		else
			return;
	})
	.then(mapblock => {
		if (mapblock){
			return({
				protected: protected(mapblock),
				generated: true
			});
		} else {
			return {
				protected: false,
				generated: false
			};
		}
	});
};
