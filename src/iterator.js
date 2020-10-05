const executor = require("./executor");

function find_y_bounds(){
  return executor(`
    select
      min(posy) as miny, max(posy) as maxy
    from blocks
  `, [], { single_row: true });
}

function find_z_bounds(posx, posy){
  return executor(`
    select
      min(posz) as minz, max(posz) as maxz
    from blocks
    where posx = $1 and posy = $2
  `, [posx, posy], { single_row: true });
}


function shift_x(ctx){
  return new Promise(resolve => {
    if (ctx.pos.x > 400){
      // end
      ctx.pos.x = -400;
      resolve(true);

    } else {
      // shift
      ctx.pos.x++;
      resolve(false);
    }
  });
}

function shift_y(ctx){
  if (!ctx.maxy){
    // initialize

  } else if (ctx.pos.y >= ctx.maxy) {
    // end
    return true;

  } else {
    // shift
    ctx.pos.y++;

  }
}

function shift_z(ctx){
  if (!ctx.maxz){
    // initialize
    return find_z_bounds(ctx.pos.x * 5, ctx.pos.y * 5)
    .then(bounds => {
      if (!bounds.maxz){
        // end
        return true;

      } else {
        ctx.maxz = bounds.maxz;
        ctx.pos.z = bounds.minz;
        return false;
      }
    });
  } else if (ctx.pos.z > ctx.maxz) {
    // reset
    ctx.maxz = null;
    ctx.pos.z = null;
    return Promise.resolve(true);

  } else {
    // shift
    ctx.pos.z++;
    return Promise.resolve(false);
  }
}

module.exports = function(){
  const ctx = {
    pos: { x: -400, y: null, z: -400 },
    maxy: null
  };

  return function(){
    if (!ctx.maxy){
      return find_y_bounds()
      .then(bounds => {
        ctx.maxy = Math.floor(bounds.maxy / 5) + 1;
        ctx.pos.y = Math.floor(bounds.miny / 5) - 1;

        return find_z_bounds(ctx.pos.x * 5, ctx.pos.y * 5);
      })
      .then(bounds => {
        ctx.pos.z = Math.floor(bounds.minz / 5) - 1;
        ctx.maxz = Math.floor(bounds.maxz / 5) + 1;
        return ctx.pos;
      });
    }

    return shift_z(ctx)
    .then(endz => {

      if (endz){
        return shift_x(ctx)
        .then(endx => {

          if (endx) {
            return shift_y(ctx);
          }
        });
      }

    })
    .then(endy => endy ? null : ctx.pos);

  };
};
