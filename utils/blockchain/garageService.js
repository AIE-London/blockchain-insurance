var blockchainInvoke = require('./blockchainService');
var config = require('config');

var addGarageReport = function(garageReport, callback){
  // We can be assured the garageReport has the relevant fields because of schema validaiton ont he endpoint
  var args = [garageReport.claimId, garageReport.garage, garageReport.estimatedCost, garageReport.writeOff];
  blockchainInvoke.invoke("addGarageReport", args, callback);
};

module.exports = {
  addGarageReport: addGarageReport
};
