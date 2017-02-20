/**
 * Created by dcotton on 04/11/2016.
 */
var unless = function(path, middleware) {
  return function(req, res, next) {
    if(Object.prototype.toString.call(path) === '[object Array]' ){
      console.log("isarray "+JSON.stringify(path));
      for(var i = 0; i<path.length; i++){
        console.log(path[i] +' may equal ' + req.path);
        if (path[i] === req.path || isWildcardMatch(path[i], req.path)) {
          return next();
        }
        console.log(path[i] +' !== ' + req.path);
      }
      return middleware(req, res, next);
    }
    else {
      if (path === req.path || isWildcardMatch(path, req.path)) {
        return next();
      } else {
        return middleware(req, res, next);
      }
    }
  };
};

var isWildcardMatch = function(path, requestPath) {
  if (endsWith(path, "*")) {
    if (requestPath.startsWith(path.replace("*", ""))) {
      return true;
    }
  }

  return false;
}

var endsWith = function(toTest, toMatch) {
  return toTest.indexOf(toMatch, this.length - toMatch.length) !== -1;
}

module.exports = {
  unlessRoute: function(path, middleware){
    return unless(path, middleware);
  }
};
