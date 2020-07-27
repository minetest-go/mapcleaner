const iterator = require("./iterator");
const checkchunk = require("./checkchunk");

const it = iterator();

function worker() {
  it()
  .then(pos => {
    if (pos){
      checkchunk(pos)
      .then(result => {
        //TODO: check if protected/generated
        setTimeout(worker, 10);
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
