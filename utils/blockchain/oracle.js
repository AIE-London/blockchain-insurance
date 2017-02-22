var vehicleValuationService = require('../vehicle/valuationService')
var blockchainInvoke = require('./blockchainService')
var NodeCache = require("node-cache");
var config = require('config')

var requestIdCache = new NodeCache( { stdTTL: 600, checkperiod: 60 } );

var getVechicleValuationAndCallbackToChain = function(requestId, styleId, mileage, callbackFunction, callback) {

  if (!requestIdCache.get(requestId)) {
    requestIdCache.set(requestId, true);
    vehicleValuationService.getVehicleValuation(styleId, mileage, function (vehicleValue) {
      console.log("Received value: " + vehicleValue);
      callbackVehicleValuationToChaincode(requestId, callbackFunction, parseInt(vehicleValue), callback);
    });
  } else {
    //Already requested by different peer, so just callback straight away
    callback();
  }
};

var callbackVehicleValuationToChaincode = function(requestId, callbackFunctionName, vehicleValue, callback) {
  callbackToChaincode(callbackFunctionName, [requestId, "" + vehicleValue], callback);
};

var callbackToChaincode = function(callbackFunctionName, args, callback) {
  blockchainInvoke.invoke(callbackFunctionName, args, config.blockchain.oracleUser, callback);
};

module.exports = {
  requestVehicleValuation: function(requestId, styleId, mileage, callbackFunctionName, callback){
    getVechicleValuationAndCallbackToChain(requestId, styleId, mileage, callbackFunctionName, callback);
  }
};

