var blockchainService = require('./blockchainService');
var config = require('config');

var raiseClaim = function(claim, callback){
  // We can be assured the claim has the relevant fields because of schema validaiton ont he endpoint
  var args = [claim.relatedPolicy, claim.description, claim.incidentDate, claim.type];
  blockchainService.invoke("createClaim", args, callback);
};

var getFullHistory = function(username, callback){
  blockchainService.query("retrieveAllClaims", [], username, callback);
};

module.exports = {
  raiseClaim: raiseClaim,
  getFullHistory: getFullHistory
};
