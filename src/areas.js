const fs = require('fs');

let areas = [];

module.exports.parse = function(filename){
  console.log(`Reading areas from '${filename}'`);
  let rawdata = fs.readFileSync(filename);
  areas = JSON.parse(rawdata);
  console.log(`Registered ${areas.length} areas`);
};
