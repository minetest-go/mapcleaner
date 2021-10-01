const executor = require("./executor");

module.exports.find_y_bounds = function(){
  return executor(`
    select
      min(posy) as miny, max(posy) as maxy
    from blocks
  `, [], { single_row: true });
};

module.exports.find_z_bounds = function(minx, maxx, miny, maxy){
  return executor(`
    select
      min(posz) as minz, max(posz) as maxz
    from blocks
    where posx >= $1 and posx <= $2
      and posy >= $3 and posy <= $4
  `, [minx, maxx, miny, maxy], { single_row: true });
};
