var vehicleValuationService = require('../vehicle/valuationService')
var NodeCache = require("node-cache");
var hfc = require('hfc');

//TODO
var PEER_ADDRESS = "localhost:7051";
var MEMBERSRVC_ADDRESS = "localhost:7054";
var KEYSTORE_PATH =  "/Users/Craig/keyValStore";

var chain = hfc.newChain("insurance");

// Configure the KeyValStore which is used to store sensitive keys
// as so it is important to secure this storage.
// The FileKeyValStore is a simple file-based KeyValStore, but you
// can easily implement your own to store whereever you want.
// To work correctly in a cluster, the file-based KeyValStore must
// either be on a shared file system shared by all members of the cluster
// or you must implement you own KeyValStore which all members of the
// cluster can share.
chain.setKeyValStore( hfc.newFileKeyValStore(KEYSTORE_PATH) );

// Set the URL for membership services
chain.setMemberServicesUrl("grpc://" + MEMBERSRVC_ADDRESS);

// Add at least one peer's URL.  If you add multiple peers, it will failover
// to the 2nd if the 1st fails, to the 3rd if both the 1st and 2nd fails, etc.
chain.addPeer("grpc://" + PEER_ADDRESS);

chain.setDevMode(true);

var requestIdCache = new NodeCache( { stdTTL: 600, checkperiod: 60 } );


var getVechicleValuationAndCallbackToChain = function(requestId, styleId, mileage, callbackFunction, callback) {

  if (!requestIdCache.get(requestId)) {
    requestIdCache.set(requestId, true);
    vehicleValuationService.getVehicleValuation(styleId, mileage, function (vehicleValue) {
      console.log("Received value: " + vehicleValue);
      callbackVehicleValuationToChaincode(requestId, callbackFunction, parseInt(vehicleValue), callback);
    });
  }
};

var callbackVehicleValuationToChaincode = function(requestId, callbackFunctionName, vehicleValue, callback) {
  callbackToChaincode(callbackFunctionName, [requestId, "" + vehicleValue], callback);
};

var callbackToChaincode = function(callbackFunctionName, args, callback) {

  console.log("Enrolling");
  //chain.enroll("bob", "NOE63pEQbL25", function(err, user) {
  //chain.enroll("WebAppAdmin", "DJY27pEnl16d", function (err, user) {
  chain.enroll("jim", "6avZQLwcUe9b", function (err, user) {
    if (err) {
      console.error(err);
      console.log("Attemping to get user");

      chain.getUser("jim", function (err, userViaGet) {
        if (err) {
          console.error(err);
          return;
        }

        invokeCallback(callbackFunctionName, args, userViaGet, callback);
      });
      return;
    }

    invokeCallback(callbackFunctionName, args, user, callback);

  });
};

function invokeCallback(callbackFunctionName, args, user, callback) {
  var invokeRequest = {
    // Name (hash) required for invoke
    chaincodeID: "insurance",
    // Function to trigger
    fcn: callbackFunctionName,
    // Parameters for the invoke function
    args: args
  };

  console.log("Invoke request: " + JSON.stringify(invokeRequest));
  var tx = user.invoke(invokeRequest);

  // Listen for the 'submitted' event
  tx.on('submitted', function(results) {
    callback();
    console.log("submitted invoke: %j",results);
  });
  // Listen for the 'complete' event.
  tx.on('complete', function(results) {
    console.log("completed invoke: %j",results);
  });
  // Listen for the 'error' event.
  tx.on('error', function(err) {
    console.log("error on invoke: %j",err);
  });
}

module.exports = {
  requestVehicleValuation: function(requestId, styleId, mileage, callbackFunctionName, callback){
    getVechicleValuationAndCallbackToChain(requestId, styleId, mileage, callbackFunctionName, callback);
  }
};

