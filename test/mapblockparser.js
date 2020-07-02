const assert = require('assert');
const mapblockparser = require("../src/mapblockparser");
const fs = require("fs");

describe('mapblockparser', function() {
  it('deserializes node names', function() {
		const data = fs.readFileSync("./test/testdata/0.0.0");
		return mapblockparser.parse(data)
    .then(function(mapblock){
      assert.equal(mapblock.version, 28);
      assert.equal(mapblock.static_objects_count, 1);
      assert.equal(mapblock.static_objects_version, 0);
      assert.equal(mapblock.node_names.some(name => name == "air"), true);
    });
	});
});
