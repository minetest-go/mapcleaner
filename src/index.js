
const app = require("./app");
const executor = require("./executor");
const mapblockparser = require("./mapblockparser");

const port = process.env.PORT || 8080;

app.listen(+port, () => {
  console.log(`Listening on http://127.0.0.1:${port}`);
});

executor(`
  select
    count(*) as count,
    min(posx) as minx,
    min(posy) as miny,
    min(posz) as minz,
    max(posx) as maxx,
    max(posy) as maxy,
    max(posz) as maxz
  from blocks`, [], { single_row: true })
.then(res => {
  console.log(res);

  let { minx, miny, minz, maxx, maxy, maxz } = res;
  const pos = { x: minx, y: miny, z: minz };

  function shift(pos){
    pos.x++;

    if (pos.x > maxx){
      pos.x = minx;
      pos.z++;
    }

    if (pos.z > maxz){
      pos.z = minz;
      pos.y++;
    }

    if (pos.y > maxy){
      // done
      return true;
    }

    return false;
  }

  function worker(pos){
    executor(`
      select data
      from blocks
      where posx = $1 and posy = $2 and posz = $3
    `, [pos.x,pos.y,pos.z], { single_row: true })
    .then(block => {
      console.log(pos, block);
      if (block)
        return mapblockparser.parse(block.data);
      else
        return;
    })
    .then(mapblock => {
      console.log(pos, mapblock);
      const done = shift(pos);

      if (!done){
        setTimeout(function(){
          worker(pos);
        }, 50);
      }
    });
  }

  worker(pos);
});
