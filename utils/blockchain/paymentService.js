var claimService = require('./claimService');
var emailService = require('../integration/emailService');
var config = require('config');

//Payments would be made off chain and then call back into the chain to confirm payment
var payoutClaim = function(claimId, policyId, policyOwner, invokerUsername) {
  console.log("'Paying out' claim with id: " + claimId);
  
  //Actual payment logic would start here in a production system
  
  //Payment made

  //Set claim as paid in the blockchain
  claimService.confirmPaidOut(claimId, invokerUsername, function(result) {
    if (result.error) {
      console.error("There was a problem when marking the claim as paid in the blockchain: " + result.error);
    } else {
      console.log("confirmPaidOut invoked for claimId: " + claimId);

      sendEmail(claimId, policyId, policyOwner);
    }
  });
}

var sendEmail = function(claimId, policyId, policyOwner) {
  var emailAddress = getEmailAddress(policyOwner);

  if (!emailAddress) {
    console.error("Cannot get email address for user: " + policyOwner);
    return;
  }

  emailService.sendEmail("Claim Paid", "Your claim with id " + claimId + " has been paid out", [emailAddress], "Blockchain Insurance");
};

//TODO Email addresses should probably be configured elsewhere
var getEmailAddress = function(username) {
  var users = config.blockchain.users;
  for (var i = 0; i < users.length; i++) {
    var user = users[i];
    if (user.enrollmentId == username) {
      return user.emailAddress
    }
  }

  return;
}

module.exports = {
  payoutClaim: payoutClaim
};
