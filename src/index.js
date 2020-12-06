const worker = require("./worker");
const areas = require("./areas");
const app = require("./app");
require("./api");

let area_parse_handler = null;

if (process.env.AREAS_FILE) {
  areas.parse(process.env.AREAS_FILE);
  // recheck areas every 5 minutes
  area_parse_handler = setInterval(() => areas.parse(process.env.AREAS_FILE), 5*60*1000);
}

app.listen(8080, err => {
	if (err){
    console.error(err);
	} else {
		console.log('Listening on http://127.0.0.1:8080');
  }
});

worker()
.then(result => {
  console.log(`processed ${result.chunkcount} chunks, removed ${result.removecount} chunks`);
  if (area_parse_handler !== null){
    clearInterval(area_parse_handler);
  }
});
