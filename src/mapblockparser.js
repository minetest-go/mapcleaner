const zlib = require("zlib");

module.exports.parse = data => new Promise(function(resolve, reject) {
	const buffer = Buffer.from(data);

	let offset = 0;

	// version stuff
	const version = buffer.readUInt8(offset);

	if (version < 25 || version > 28) {
		return reject("mapblock version not supported: " + version);
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
		return reject("content/param width mismatch!");
	}

	//mapdata (blocks)

	if (version >= 27) {
		offset = 6;
	} else {
		offset = 4;
	}

	const mapdata_buffer = buffer.subarray(offset);

	let inflate = zlib.createInflate();
	inflate.on("data", function(mapdata){
		console.log("mapdata", mapdata);
		console.log(inflate.bytesWritten);

		offset += inflate.bytesWritten;
		const metadata_buffer = buffer.subarray(offset);
		inflate = zlib.createInflate();

		inflate.on("data", function(metadata){
			console.log("metadata", metadata);
			console.log(inflate.bytesWritten);
			offset += inflate.bytesWritten;

			//static objects version
			const static_objects_version = buffer.readUInt8(offset);
			offset++;

			const static_objects_count = buffer.readUInt16BE(offset);
			offset += 2;

			for (let i=0; i < static_objects_count; i++) {
				offset += 13;
				const dataSize = buffer.readUInt16BE(offset);
				console.log("dataSize", dataSize);
				offset += dataSize + 2;
			}

			//timestamp
			offset += 4;

			//mapping version
			offset++;

			const numMappings = buffer.readUInt16BE(offset);
			console.log("numMappings", numMappings);

			const node_names = [];

			offset += 2;
			for (let i=0; i < numMappings; i++) {
				//const nodeId = buffer.readUInt16BE(offset);
				offset += 2;

				const nameLen = buffer.readUInt16BE(offset);
				offset += 2;

				const blockName = buffer.subarray(offset, offset+nameLen).toString();
				offset += nameLen;

				console.log("blockName", blockName);
				node_names.push(blockName);
			}

			resolve({
				version: version,
				static_objects_count: static_objects_count,
				static_objects_version: static_objects_version,
				node_names: node_names
			});
		});

		inflate.write(metadata_buffer);
	});

	inflate.write(mapdata_buffer);

});
