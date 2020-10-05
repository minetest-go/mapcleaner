const fs = require('fs');
const coordinates = require("./coordinates");

function get_key(pos){
  return `${pos.x}/${pos.y}/${pos.z}`;
}

let areas = [];
let protected_chunks = [];

module.exports.parse = function(filename){
  console.log(`Reading areas from '${filename}'`);
  let rawdata = fs.readFileSync(filename);
  areas = JSON.parse(rawdata);
  console.log(`Registered ${areas.length} areas`);

  let chunk_count = 0;
  areas.forEach(function(area){
    if (!area || !area.pos1 || !area.pos2){
      return;
    }
    const chunkpos1 = coordinates.get_chunkpos_from_pos(area.pos1);
    const chunkpos2 = coordinates.get_chunkpos_from_pos(area.pos2);

    for (let x=chunkpos1.x; x<=chunkpos2.x; x++){
      for (let y=chunkpos1.y; y<=chunkpos2.y; y++){
        for (let z=chunkpos1.z; z<=chunkpos2.z; z++){
          const key = get_key({ x:x, y:y, z:z });
          protected_chunks[key] = true;
          chunk_count++;
        }
      }
    }
  });

  console.log(`Area protected chunks: ${chunk_count}`);
};

module.exports.is_chunk_protected = function(chunkpos){
  return protected_chunks[get_key(chunkpos)];
};
