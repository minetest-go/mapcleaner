const worker = require("./worker");
const areas = require("./areas");

if (process.env.AREAS_FILE) {
  areas.parse(process.env.AREAS_FILE);
}

worker()
.then(result => {
  console.log(`processed ${result.chunkcount} chunks, removed ${result.removecount} chunks`);
});
