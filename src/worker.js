const iterator = require("./iterator");
const checkchunkwithneighbours = require("./checkchunkwithneighbours");
const removechunk = require("./removechunk");

const it = iterator();

function worker() {
  it()
  .then(pos => {
    if (pos){
      if (!pos.z){
        // no data in row, skip
        setTimeout(worker, 0);
        return;
      }
      console.log("worker-chunk", pos);
      checkchunkwithneighbours(pos)
      .then(result => {
        console.log("check-chunk", pos, result);

        if (!result.protected && result.generated){
          //not protected and generated, remove
          console.log("removing chunk", pos);
          removechunk(pos)
          .then(() => {
            //proceed with next chunk
            setTimeout(worker, 500);
          });
        } else {
          //proceed with next chunk
          setTimeout(worker, 0);
        }
      });
    } else {
      //done
      console.log("done!");
    }
  })
  .catch(e => {
    console.log(e);
  });
}

module.exports = worker;
