const assert = require('assert');
const mapblockparser = require("../src/mapblockparser");
const fs = require("fs");

describe('mapblockparser', function() {
  it('deserializes node names', function() {
		const data = fs.readFileSync("./test/testdata/0.0.0");
		const mapblock = mapblockparser.parse(data);

		assert.equal(mapblock.version, 28);
	});
});
