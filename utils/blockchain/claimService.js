var blockchainService = require('./blockchainService');
var config = require('config');

var raiseClaim = function(claim, username, callback){
  // We can be assured the claim has the relevant fields because of schema validaiton ont he endpoint
  var args = [claim.relatedPolicy, claim.description, claim.incidentDate, claim.type];
  blockchainService.invoke("createClaim", args, username, callback);
};

var getFullHistory = function(username, callback){
  blockchainService.query("retrieveAllClaims", [], username, callback);
};

var makeClaimAgreement = function(claimId, agreement, callback){
  blockchainService.invoke("agreePayoutAmount", [claimId, agreement], callback);
};

module.exports = {
  raiseClaim: raiseClaim,
  getFullHistory: getFullHistory,
  makeClaimAgreement: makeClaimAgreement
};
