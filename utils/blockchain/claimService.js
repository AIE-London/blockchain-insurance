var blockchainInvoke = require('./invokeService');
var config = require('config');

var raiseClaim = function(claim, callback){
  // We can be assured the claim has the relevant fields because of schema validaiton ont he endpoint
  var args = [claim.relatedPolicy, claim.description, claim.incidentDate, claim.type];
  blockchainInvoke.invoke("createClaim", args, callback);
};

module.exports = {
  raiseClaim: raiseClaim
};
