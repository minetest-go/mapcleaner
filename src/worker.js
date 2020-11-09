const checkchunkwithneighbours = require("./checkchunkwithneighbours");
const removechunk = require("./removechunk");
const bounds = require("./bounds");
const snooze = ms => new Promise(resolve => setTimeout(resolve, ms));

async function worker() {

  const ybounds = await bounds.find_y_bounds();
  let minychunk = Math.floor((ybounds.miny + 2) / 5);
  const maxychunk = Math.floor((ybounds.maxy + 2) / 5);
  let chunkcount = 0;
  let removecount = 0;

  if (process.env.STARTYCHUNK) {
    minychunk = +process.env.STARTYCHUNK;
  }

  for (let chunky = minychunk; chunky <= maxychunk; chunky++){
    const miny = (chunky * 5) - 2;
    const maxy = miny + 4;
    const progress = Math.floor( (chunky - minychunk) / (maxychunk - minychunk) * 100 );

    console.log(`chunky: ${chunky} y-mablocks: ${miny} to ${maxy} progress: ${progress}%`);

    for (let chunkx = -400; chunkx <= 400; chunkx++){
      const minx = (chunkx * 5) - 2;
      const maxx = minx + 4;

      const zbounds = await bounds.find_z_bounds(minx, maxx, miny, maxy);
      if (zbounds.minz === null)
        // no data in this stride
        continue;

      //console.log(" zbounds", zbounds);

      const minzchunk = Math.floor((zbounds.minz + 2) / 5);
      const maxzchunk = Math.floor((zbounds.maxz + 2) / 5);

      for (let chunkz = minzchunk; chunkz <= maxzchunk; chunkz++){
        const chunkpos = {
          x: chunkx,
          y: chunky,
          z: chunkz
        };

        const result = await checkchunkwithneighbours(chunkpos);
        if (!result.protected && result.generated){
          //not protected and generated, remove
          console.log("removing chunk", chunkpos);
          await removechunk(chunkpos);
          await snooze(1000);
          removecount++;
        }

        chunkcount++;
      }

    }
  }

  return {
    chunkcount: chunkcount,
    removecount: removecount
  };
}

module.exports = worker;
