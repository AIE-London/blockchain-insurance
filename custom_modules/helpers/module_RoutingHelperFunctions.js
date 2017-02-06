/**
 * Created by dcotton on 04/11/2016.
 */
var unless = function(path, middleware) {
  return function(req, res, next) {
    if(Object.prototype.toString.call(path) === '[object Array]' ){
      console.log("isarray "+JSON.stringify(path));
      for(var i = 0; i<path.length; i++){
        console.log(path[i] +' may equal ' + req.path);
        if (path[i] === req.path) {
          return next();
        }
        console.log(path[i] +' !== ' + req.path);
      }
      return middleware(req, res, next);
    }
    else {
      if (path === req.path) {
        return next();
      } else {
        return middleware(req, res, next);
      }
    }
  };
};
module.exports = {
  unlessRoute: function(path, middleware){
    return unless(path, middleware);
  }
};
