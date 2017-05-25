var hfc = require('hfc');
var config = require('config')

let VCAP_SERVICES = process.env.VCAP_SERVICES;
if (!VCAP_SERVICES) {
  VCAP_SERVICES = require('../../vcap-services.json');
} else {
  VCAP_SERVICES = JSON.parse(VCAP_SERVICES);
}
let BLOCKCHAIN_MEMBER_SVC = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.ca;
let MEMBER_SVC = BLOCKCHAIN_MEMBER_SVC[Object.keys(BLOCKCHAIN_MEMBER_SVC)[0]];

let BLOCKCHAIN_PEER = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.peers[0];
let PEER_ADDRESS = BLOCKCHAIN_PEER.api_host + BLOCKCHAIN_PEER.api_port;
let MEMBERSRVC_ADDRESS = MEMBER_SVC.api_host + MEMBER_SVC.api_port;
let KEYSTORE_PATH = __dirname + config.blockchain.keystorePath;
let EVENTS_ADDRESS = BLOCKCHAIN_PEER.event_host + BLOCKCHAIN_PEER.event_port;
let CHAINCODE_ID  = process.env.CHAINCODE_ID;

var ATTRS = ['username', 'role'];

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

chain.setDevMode(config.blockchain.devMode);

// Configure the event hub
chain.eventHubConnect("grpc://" + EVENTS_ADDRESS);
var eventHub = chain.getEventHub();
var registeredEventListenerIds = [];

// disconnect when done listening for events
process.on('exit', function() {
  for (var i = 0; i < registeredEventListenerIds.length; i++) {
    eventHub.unregisterChaincodeEvent(registeredEventListenerIds[i]);
  }
  chain.eventHubDisconnect();
});

var login = function(name, secret, callback){
  console.log("Enrolling: " + name);
  chain.enroll(name, secret, function(err, user){
    if (err){
      console.error(err);
      console.log("Failed to enroll, so getting user..");
      chain.getUser(name, function(err, userViaGet){
        if (err) {
          console.error(err);
          callback({error: err});
        }
        console.log(userViaGet);
        callback(userViaGet);
      });
    }else {
      callback(user);
    }
  });
};

var loginAndInvoke = function(functionName, args, username, callback) {

  console.log("Enrolling " + username);
  chain.enroll(username, "123456789123", function (err, user) { // Enroll secret only used in dev mode
    if (err) {
      console.error(err);
      console.log("Attemping to get user");

      chain.getUser(username, function (err, userViaGet) {
        if (err) {
          console.error(err);
          callback(err);
          return;
        }
        invoke(functionName, args, userViaGet, callback);
      });
      return;
    }
    invoke(functionName, args, user, callback);

  });
};

var invoke = function(functionName, args, user, callback) {
  var invokeRequest = {
    // Name (hash) required for invoke
    chaincodeID: CHAINCODE_ID,
    // Function to trigger
    fcn: functionName,
    // Parameters for the invoke function
    args: args,
    attrs: ATTRS
  };

  console.log("Invoke request: " + JSON.stringify(invokeRequest));
  var tx = user.invoke(invokeRequest);

  // Listen for the 'submitted' event
  tx.on('submitted', function(results) {
    //callback(); -- Removed as we expect 'complete' to be triggered also, where we'd like the response
    console.log("submitted invoke: %j",results);
  });
  // Listen for the 'complete' event.
  tx.on('complete', function(results) {
    console.log("completed invoke: %j",results);
    callback({"results": results});
  });
  // Listen for the 'error' event.
  tx.on('error', function(err) {
    callback(err);
    console.log("error on invoke: %j",err);
  });
};

var loginAndQuery = function(funcionName, args, username, callback){
  login(username, "123456789123", function(user){ // Enroll secret only used in dev mode
    console.log("--- USER ---");
    console.log(user);
    if (user.error){
      callback(user.error);
    } else {
      query(funcionName, args, user, callback);
    }
  })
};

var query = function(functionName, args, user, callback){
  var queryRequest = {
    chaincodeID: CHAINCODE_ID,
    fcn: functionName,
    args: args,
    attrs: ATTRS
  };

  console.log("QUERY REQUEST: " + JSON.stringify(queryRequest));

  var tx = user.query(queryRequest);

  tx.on('submitted', function(results) {
    //callback(); -- Removed as we expect 'complete' to be triggered also, where we'd like the response
    console.log("submitted query: %j",results);
  });

  tx.on('complete', function(results){
    console.log("completed query: %j", results);
    callback({"results": new Buffer(results.result,'hex').toString()});
  });

  tx.on('error', function(err){
    console.log("error on query: %j", err);
    callback(err);
  })
};

var registerEventListener = function(eventName, callback){
  console.log("Registering listener for event name: " + eventName);
  var registrationId = eventHub.registerChaincodeEvent(CHAINCODE_ID, eventName, callback);

  console.log("Registered listener for event name: " + eventName + " with id: " + JSON.stringify(registrationId));
  registeredEventListenerIds.push(registrationId);
};

module.exports = {
  invoke: function(functionName, args, username, callback){
    loginAndInvoke(functionName, args, username, callback);
  },
  query: function(functionName, args, username, callback){
    loginAndQuery(functionName, args, username, callback);
  },
  registerEventListener : registerEventListener
};
