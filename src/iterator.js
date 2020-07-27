
module.exports = function(){
  const pos = { x: -2000, y: -2000, z: -2000 };

  return function(){
    return new Promise(resolve => {
      if (pos.z >= 2000){
        pos.x++;
        pos.z = -2000;
      } else {
        pos.z++;
      }

      if (pos.x >= 2000){
        pos.y++;
        pos.x = -2000;
      }

      if (pos.y > 2000){
        resolve(null);
      }

      resolve(pos);
    });
  };
};
