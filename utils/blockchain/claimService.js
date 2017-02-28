var blockchainService = require('./blockchainService');
var config = require('config');

var raiseClaim = function(claim, username, callback){
  // We can be assured the claim has the relevant fields because of schema validaiton ont he endpoint

  // Map UI Types to Blockchain Types
  if (claim.type.toLowerCase() === "multiple party"){
    claim.type = "multiple_parties";
  } else if (claim.type.toLowerCase() === "Single Party"){
    claim.type = "single_party";
  }

  var args = [claim.relatedPolicy, claim.description, claim.incidentDate, claim.type];

  // Only required for multi-party; added here rather than within multi-party if so can push to args
  if (claim.otherPartyReg){
    args.push(claim.otherPartyReg, (claim.atFault === true))
  };

  blockchainService.invoke("createClaim", args, username, callback);
};

var getFullHistory = function(username, callback){
  blockchainService.query("retrieveAllClaims", [], username, callback);
};

var makeClaimAgreement = function(claimId, agreement, username, callback){
  blockchainService.invoke("agreePayoutAmount", [claimId, agreement.toString()], username, callback);
};

var confirmPaidOut = function(claimId, username, callback) {
  blockchainService.invoke("confirmPaidOut", [claimId], username, callback);
}

module.exports = {
  raiseClaim: raiseClaim,
  getFullHistory: getFullHistory,
  makeClaimAgreement: makeClaimAgreement,
  confirmPaidOut: confirmPaidOut
};
