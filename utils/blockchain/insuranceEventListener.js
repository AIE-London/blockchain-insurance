var blockchainService = require('./blockchainService');
var paymentService = require('./paymentService');
var config = require('config');

var init = function() {

  //Register the listeners
  blockchainService.registerEventListener("ClaimSettled", payoutCallback);
  blockchainService.registerEventListener("InsurerPaymentPaid", payoutCallback);
}

var payoutCallback = function(event){
  var payload = JSON.parse(event.payload.toString())
  console.log("Received claim settled event for claimId: " + payload.claimId);

  paymentService.payoutClaim(payload.claimId, payload.policyId, payload.linkedClaimId);
};

module.exports = {
  init: init
};
