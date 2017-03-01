var blockchainService = require('./blockchainService');
var paymentService = require('./paymentService');
var config = require('config');

var init = function() {

  //Register the claim settled listener
  blockchainService.registerEventListener("ClaimSettled", function(event){
    console.log("event: " + event.toString())
    var payload = JSON.parse(event.payload.toString())
    console.log("Received claim settled event for claimId: " + payload.claimId);
    
    paymentService.payoutClaim(payload.claimId,
      payload.policyId, payload.user);
  });
}

module.exports = {
  init: init
};
