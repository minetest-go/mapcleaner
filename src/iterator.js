const executor = require("./executor");
const coordinates = require("./coordinates");

function find_bounds(){
  return executor(`
    select
      min(posx) as minx, max(posx) as maxx,
      min(posy) as miny, max(posy) as maxy,
      min(posz) as minz, max(posz) as maxz
    from blocks
  `, [], { single_row: true });
}

module.exports = function(){
  let bounds;
  const pos = { x: -2000, y: -2000, z: -2000 };

  return function(){
    if (!bounds){
      return find_bounds()
      .then(blockbounds => {
        console.log("blockbounds", blockbounds);
        bounds = {
          min: coordinates.get_chunkpos_from_mapblock({
            x: blockbounds.minx,
            y: blockbounds.miny,
            z: blockbounds.minz
          }),
          max: coordinates.get_chunkpos_from_mapblock({
            x: blockbounds.maxx,
            y: blockbounds.maxy,
            z: blockbounds.maxz
          })
        };

        console.log("chunkbounds", bounds);

        pos.x = bounds.min.x - 1;
        pos.y = bounds.min.y - 1;
        pos.z = bounds.min.z - 1;

        bounds.max.x += 2;
        bounds.max.y += 2;
        bounds.max.z += 2;

        return pos;
      });
    }
    return new Promise(resolve => {
      if (pos.z >= bounds.max.z){
        // shift x
        pos.x++;
        pos.z = bounds.min.z;
      } else {
        // shift z
        pos.z++;
      }

      if (pos.x >= bounds.max.x){
        // shift y
        pos.y++;
        pos.x = bounds.min.x;
      }

      if (pos.y > bounds.max.y){
        resolve(null);
      }

      resolve(pos);
    });
  };
};
