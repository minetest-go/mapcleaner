const worker = require("./worker");
const areas = require("./areas");

let area_parse_handler = null;

if (process.env.AREAS_FILE) {
  areas.parse(process.env.AREAS_FILE);
  // recheck areas every 5 minutes
  area_parse_handler = setInterval(() => areas.parse(process.env.AREAS_FILE), 5*60*1000);
}

worker()
.then(result => {
  console.log(`processed ${result.chunkcount} chunks, removed ${result.removecount} chunks`);
  if (area_parse_handler !== null){
    clearInterval(area_parse_handler);
  }
});
