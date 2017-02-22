/**
 * Created by dcotton on 04/11/2016.
 */
var authEndpoint = "http://aston-jwt-auth.eu-gb.mybluemix.net/app/";

var request = require('request');

var config =  require('config');

module.exports = {
  signJWT: function (data, options, callback) {
    var requestData = {};
    requestData.token = data;
    requestData.secret = options.secret;
    console.log(authEndpoint+options.name+"/token/sign");
    console.log(requestData);
    request({
      url: authEndpoint+options.name+"/token/sign",
      method: "POST",
      json: requestData
    }, function (error, response, body) {
      if(response.statusCode == 200){
        console.log(body);
        callback(body.token);
      }
      else{
        callback(false);
      }
    });
  },
  verifyJWT: function (data, options, callback) {
    var requestData = {};
    requestData.token = data;
    requestData.secret = options.secret;
    request({
      url: authEndpoint+options.name+"/token/verify",
      method: "POST",
      json: requestData
    }, function (error, response, body) {
      if(response.statusCode == 200){
        callback(body.username);
      }
      else{
        callback(false);
      }
    });
  },
  middleware: function (req, response, next) {
    var header = "";
    var origin = req.headers.origin;
    if((req.method.toUpperCase() === 'OPTIONS') || (origin === 'http://aston-swagger-ui.eu-gb.mybluemix.net') || (origin === 'https://aston-swagger-ui.eu-gb.mybluemix.net')){
      next();
    }
    else{
      if(req.get("Authorization") && req.get("Authorization").split(" ").length > 1 && req.get("Authorization").split(" ")[1]){
        header = req.get("Authorization").split(" ")[1];
      }
      console.log("HEADER: " + header);

      var auth = config.web.auth;
      auth.name = req.params.username;

      //verify JWT with Auth header in syntax "Bearer 12321k3123jlkj"
      module.exports.verifyJWT(header, auth , function(ok){
        if(ok){
          req.user = ok;
          next();
        }
        else{
          console.error("Invalid JWT");
          response.setHeader('Content-Type', 'application/json');
          response.statusCode = 401;
          response.write(JSON.stringify({reason: "Unauthorized"}));
          response.end();
          return;
        }
      });
    }
  },
  checkAuthorized: function (req, response, next) {
    var origin = req.headers.origin;
    if (req.user === req.params.username || (req.method.toUpperCase() === 'OPTIONS') || (origin === 'http://aston-swagger-ui.eu-gb.mybluemix.net') || (origin === 'https://aston-swagger-ui.eu-gb.mybluemix.net')) {
      next();
    } else {
      response.setHeader('Content-Type', 'application/json');
      response.statusCode = 403;
      response.write(JSON.stringify({reason: "You are not authorized to access that resource"}));
      response.end();
      return;
    }
  },
  allowOriginsMiddleware: function(req, res, next) {
    var allowedOrigins = config.allowedOrigins;
    var origin = req.headers.origin;
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');
    res.setHeader('Access-Control-Allow-Headers', 'X-Requested-With,content-type,authorization');
    res.setHeader('Access-Control-Expose-Headers', 'Token');

    if(allowedOrigins.indexOf(origin) > -1){
      res.setHeader('Access-Control-Allow-Origin', origin);
    }
    return next();
  }
};
