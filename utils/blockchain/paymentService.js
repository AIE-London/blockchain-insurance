var claimService = require('./claimService');
var policyService = require('./policyService');
var emailService = require('../integration/emailService');
var config = require('config');

//Payments would be made off chain and then call back into the chain to confirm payment
var payoutClaim = function(claimId, policyId, linkedClaimId) {
  console.log("'Paying out' claim with id: " + claimId);
  
  //Actual payment logic would start here in a production system
  
  //Payment made

  //Iterate through all insurers and mark as paid if the claim is associated with the insurer
  usersConfig = config.blockchain.users;
  for (var i = 0; i < usersConfig.length; i++) {
    var user = usersConfig[i];
    var role = getRoleForUser(user);

    if (role == "insurer") {
      confirmPaidOutForInsurer(claimId, policyId, user.enrollmentId);

      if (linkedClaimId) {
        confirmPaidOutForInsurer(linkedClaimId, policyId, user.enrollmentId)
      }
    }
  }
}

var getRoleForUser = function(user) {
  for (var i = 0; i < user.attributes.length; i++) {
    var attribute = user.attributes[i];

    if (attribute.name == "role") {
      return attribute.value;
    }
  }
}

var confirmPaidOutForInsurer = function(claimId, policyId, insurerUsername) {

  claimService.getClaimWithId(claimId, insurerUsername, function(claim) {
    if (claim && claim.details.settlement && claim.details.settlement.payments) {
      for (var i = 0; i < claim.details.settlement.payments.length; i++) {
        var payment = claim.details.settlement.payments[i];
        
        if (payment.sender == insurerUsername) {

          //If we're not liable, dont payout until the other insurer has paid us
          if (claim.details.liable == true || hasPaidPaymentFromLiableInsurer(claim, insurerUsername)) {
            claimService.confirmPaidOut(claimId, payment.id, insurerUsername, function (result) {
              if (result.error) {
                console.error("There was a problem when marking the claim as paid in the blockchain: " + result.error);
              } else {
                console.log("confirmPaidOut invoked for claimId: " + claimId);
                if (payment.recipientType == "claimant") {
                  policyService.getPolicyWithId(policyId, insurerUsername, function(policy){
                    if (policy) {
                      sendEmail(claimId, policyId, policy.relations.owner);
                    } else {
                      console.log("Unable to send email because we failed to retrieve the policy");
                    }
                  });
                }
              }
            });
          }
        }
      }
    } else {
      console.log("Cannot get claim with id: " + claimId + " for insurer: " + insurerUsername + ". This may be expected");
    }
  })
};

var hasPaidPaymentFromLiableInsurer = function(claim, insurerUsername) {
  for (var i = 0; i < claim.details.settlement.payments.length; i++) {
    var payment = claim.details.settlement.payments[i]

    if (payment.recipient == insurerUsername && payment.senderType == "insurer" && payment.status == "paid") {
      return true;
    }
  }

  return false
};

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
