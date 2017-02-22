var blockchainService = require('./blockchainService');
var config = require('config');

var getFullHistory = function(username, callback){
  blockchainService.query("retrieveAllPolicies", [], username, callback);
};

module.exports = {
  getFullHistory: getFullHistory
};
