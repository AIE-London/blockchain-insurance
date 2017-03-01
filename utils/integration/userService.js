var userEndpoint = "http://aston-user-service.eu-gb.mybluemix.net";

var request = require("request");

var authenticate = function(body, callback){
  var requestData = body;
  request({
    url: userEndpoint + "/auth/user/login",
    method: "POST",
    json: requestData
  }, function(error, response, body) {

    console.log(body);

    if (response.statusCode == 200){
      callback(body);
    } else {
      callback(response.statusCode);
    }
  })
};

///user/{id}/push
var getUserPushTokens = function(username, callback){
  request({
    url: userEndpoint + "/user/" + username + "/push",
    method: "GET"
  }, function(error, response, body){
    if (response.statusCode == 200){
      callback(JSON.parse(body));
    } else {
      callback(response.statusCode);
    }
  })
};


module.exports = {
  authenticate: authenticate,
  getUserPushTokens: getUserPushTokens
};
