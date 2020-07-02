const zlib = require("zlib");

module.exports.parse = function(data){
	const buffer = Buffer.from(data);

	let offset = 0;

	// version stuff
	const version = buffer.readUInt8(offset);

	if (version < 25 || version > 28) {
		throw new Error("mapblock version not supported: " + version);
	}

	if (version >= 27) {
		offset = 4;
	} else {
		//u16 lighting_complete not present
		offset = 2;
	}

	// content

	const content_width = buffer.readUInt8(offset);
	const params_width = buffer.readUInt8(offset+1);

	if (content_width != 2 || params_width != 2){
		throw new Error("content/param width mismatch!")
	}

	//mapdata (blocks)

	if (version >= 27) {
		offset = 6;
	} else {
		offset = 4;
	}

	const metadata_buffer = buffer.subarray(offset);

	const inflate = zlib.createInflate();
	inflate.on("data", function(e){
		console.log("event", e);
		console.log(inflate.bytesWritten);
	});
	inflate.write(metadata_buffer);

	console.log(inflate);

	const metadata = zlib.inflateSync(metadata_buffer);

	return {
		version: version
	};
};
