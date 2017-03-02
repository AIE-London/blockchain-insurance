var blockchainService = require('./blockchainService');
var paymentService = require('./paymentService');
var config = require('config');

var init = function() {

  //Register the listeners
  blockchainService.registerEventListener("ClaimSettled", payoutCallback);
  blockchainService.registerEventListener("InsurerPaymentPaid", payoutCallback);
  blockchainService.registerEventListener("InsurerPaymentAdded", payoutCallback);
}

var payoutCallback = function(event){
  var payload = JSON.parse(event.payload.toString())
  console.log("Received claim settled event for claimId: " + payload.claimId);

  paymentService.payoutClaim(payload.claimId, payload.policyId);
};

module.exports = {
  init: init
};
