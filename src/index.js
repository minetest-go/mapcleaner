const app = require("./app");
const worker = require("./worker");

const port = process.env.PORT || 8080;

app.listen(+port, () => {
  console.log(`Listening on http://127.0.0.1:${port}`);
});

worker();
