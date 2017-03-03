var blockchainService = require('./blockchainService');
var config = require('config');

var getFullHistory = function(username, callback){
  blockchainService.query("retrieveAllPolicies", [], username, callback);
};

//For now we're getting all the policies and iterating but this is obviously inefficient.
//We should add a query to the chaincode for a specific policy id
var getPolicyWithId = function(policyId, username, callback) {
  getFullHistory(username, function(result) {
    var policies = JSON.parse(result.results);

    for (var i = 0; i < policies.length; i++) {
      if (policies[i].id == policyId) {
        callback(policies[i]);
        return
      }
    }

    //policy not found
    return callback();
  });
}

module.exports = {
  getFullHistory: getFullHistory,
  getPolicyWithId: getPolicyWithId
};
