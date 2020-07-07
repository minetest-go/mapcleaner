
const app = require("./app");
const pool = require("./pool");

const port = process.env.PORT || 8080;

app.listen(+port, () => {
  console.log(`Listening on http://127.0.0.1:${port}`);
});

pool.connect()
.then(client => {
  return client.query(`select count(*) from blocks`)
  .then(sql_res => {
    console.log(sql_res);
    client.release();
  })
  .catch(e => {
    client.release();
    console.error(e.stack);
  });
});
