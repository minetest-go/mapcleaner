const app = require("./app");
const worker = require("./worker");
const areas = require("./areas");

const port = process.env.PORT || 8080;

app.listen(+port, () => {
  console.log(`Listening on http://127.0.0.1:${port}`);
});

if (process.env.AREAS_FILE) {
  areas.parse(process.env.AREAS_FILE);
}

worker();
