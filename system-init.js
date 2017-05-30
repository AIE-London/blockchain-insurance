
var hfc = require('hfc');
let fs = require('fs');
let blockchain = require('./utils/blockchain/blockchain-helpers.js');
let deploy = require('./deploy.js');

const config = require('./config.json');

// Setup config

let VCAP_SERVICES = process.env.VCAP_SERVICES;
if (!VCAP_SERVICES) {
  VCAP_SERVICES = require('./vcap-services.json');
} else {
  VCAP_SERVICES = JSON.parse(VCAP_SERVICES);
}
let BLOCKCHAIN_MEMBER_SVC = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.ca;
let MEMBER_SVC = BLOCKCHAIN_MEMBER_SVC[Object.keys(BLOCKCHAIN_MEMBER_SVC)[0]];

let BLOCKCHAIN_PEER = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.peers[0];
var PEER_ADDRESS = BLOCKCHAIN_PEER.api_host + ":" + BLOCKCHAIN_PEER.api_port;
var MEMBERSRVC_ADDRESS = MEMBER_SVC.api_host + ":" + MEMBER_SVC.api_port;
let USERS = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.users;
var KEYSTORE_PATH = __dirname + "/keyValStore-4a780ee6";

// Create a client chain.
// The name can be anything as it is only used internally.
var chain = hfc.newChain("blockchainInsurance");

var certFile = 'us.blockchain.ibm.com.cert';

// Configure the KeyValStore which is used to store sensitive keys
// as so it is important to secure this storage.
// The FileKeyValStore is a simple file-based KeyValStore, but you
// can easily implement your own to store whereever you want.
chain.setKeyValStore( hfc.newFileKeyValStore('./keyValStore-4a780ee6') );

var cert = fs.readFileSync(certFile);

// Set the URL for member services
console.log(MEMBERSRVC_ADDRESS);
chain.setECDSAModeForGRPC(true);
chain.addPeer("grpcs://" + PEER_ADDRESS, {
  pem: cert
});
chain.setMemberServicesUrl("grpcs://" + MEMBERSRVC_ADDRESS, {
  pem: cert
});

// Enroll "WebAppAdmin" which is already registered because it is
// listed in fabric/membersrvc/membersrvc.yaml with its one time password.
// If "WebAppAdmin" has already been registered, this will still succeed
// because it stores the state in the KeyValStore
// (i.e. in '/tmp/keyValStore' in this sample).
chain.enroll(USERS[0].enrollId, USERS[0].enrollSecret, function(err, admin) {
   if (err) {
      console.error("ERROR: failed to enroll %s", err)
      // if it's failing to enroll - chances are they're already enrolled. just move on.
   }
   console.log("SUCCESS: Enrolled: %s", USERS[0].enrollId);
   // Successfully enrolled WebAppAdmin during initialization.
   // Set this user as the chain's registrar which is authorized to register other users.
   chain.setRegistrar(admin);
   // Now begin listening for web app requests
   let registrationRequests = config.users.map(user => {
    return new Promise((resolve, reject) => {
      console.log("[Registration] Submitting request for user: " + user.enrollmentId);
      chain.register({
        "enrollmentID": user.enrollmentId,
        "affiliation": user.affiliation,
        "attributes": user.attributes
      }, function(err, createdUser) {
        console.log("[Registration] SUCCESSFULLY registered user: " + user.enrollmentId);
        if (err) {
          console.error("ERROR: failed to register %s", err);
          //throw err;
        }
        // Issue an invoke request
        blockchain.enroll(PEER_ADDRESS, {
          enrollId: user.enrollmentId,
          enrollSecret: createdUser
        })
        .then(() => {
          console.log("Enrolled: " + user.enrollmentId )
          resolve();
        })
        .catch(err => {
          console.log("ERROR: failed to enroll %s",err)
          // if it's failing to enroll - chances are they're already enrolled. just mark as done and move on.
          resolve();
        });
      });
    });
   });
   Promise.all(registrationRequests).then(() => {
      deploy.deploy(PEER_ADDRESS, config.users[1]);
   })
});