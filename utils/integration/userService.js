var userEndpoint = "http://aston-user-service.eu-gb.mybluemix.net";

var request = require("request");

var authenticate = function(username, password, callback){
  var requestData = {"username": username, "password": password};
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

module.exports = {
  authenticate: authenticate
};
