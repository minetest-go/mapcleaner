const executor = require("./executor");

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
          minx: Math.floor(blockbounds.minx / 16) - 1,
          miny: Math.floor(blockbounds.miny / 16) - 1,
          minz: Math.floor(blockbounds.minz / 16) - 1,
          maxx: Math.floor(blockbounds.maxx / 16) + 1,
          maxy: Math.floor(blockbounds.maxy / 16) + 1,
          maxz: Math.floor(blockbounds.maxz / 16) + 1
        };

        console.log("chunkbounds", bounds);

        pos.x = bounds.minx;
        pos.y = bounds.miny;
        pos.z = bounds.minz;
        return pos;
      });
    }
    return new Promise(resolve => {
      if (pos.z >= bounds.maxz){
        // shift x
        pos.x++;
        pos.z = bounds.minz;
      } else {
        // shift z
        pos.z++;
      }

      if (pos.x >= bounds.maxx){
        // shift y
        pos.y++;
        pos.x = bounds.minx;
      }

      if (pos.y > bounds.maxy){
        resolve(null);
      }

      resolve(pos);
    });
  };
};
