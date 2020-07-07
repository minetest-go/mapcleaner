const pool = require("./pool");

module.exports = function(query, values, options){

  const single_row = options && options.single_row;

  return new Promise(function(resolve, reject){
    pool.connect()
    .then(client => {
      return client.query(query, values)
      .then(sql_res => {
        if (!sql_res.rows || sql_res.rows.length == 0){
          // no result
          resolve(single_row ? null : []);
        }
        resolve(single_row ? sql_res.rows[0] : sql_res.rows);
        client.release();
      })
      .catch(e => {
        client.release();
        console.error(e.stack);
        reject();
      });
    });
  });
};
