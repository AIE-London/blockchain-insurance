var path = require('path');
var util = require('util');
var appDir = path.dirname(require.main.filename);
process.env.GOPATH = appDir;

var fs = require('fs');
var config = require('config');
var hfc = require('hfc');

/** CHAINCODE **/

/**
 * Id of the deployed chaincode
 * @type {string}
 */
var chaincodeId;

/**
 * Path of the deployed chaincode
 * @type {string}
 */
var chaincodeIdPath = __dirname + "/chaincodeId";

var chain;

/** NETWORK **/
var network = {};
var deployUser;
let VCAP_SERVICES = process.env.VCAP_SERVICES;
if (!VCAP_SERVICES) {
  VCAP_SERVICES = require('../../vcap-services.json');
} else {
  VCAP_SERVICES = JSON.parse(VCAP_SERVICES);
}
let BLOCKCHAIN_MEMBER_SVC = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.ca;
let MEMBER_SVC = BLOCKCHAIN_MEMBER_SVC[Object.keys(BLOCKCHAIN_MEMBER_SVC)[0]];

let BLOCKCHAIN_PEER = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.peers[0];
var PEER_ADDRESS = BLOCKCHAIN_PEER.api_host + BLOCKCHAIN_PEER.api_port;
var MEMBERSRVC_ADDRESS = MEMBER_SVC.api_host + MEMBER_SVC.api_port;
let USERS = VCAP_SERVICES['ibm-blockchain-5-prod'][0].credentials.users;
var KEYSTORE_PATH = __dirname + config.blockchain.keystorePath;

/** FUNCTIONS **/

/**
 * @sets {chain, network}
 */
var configureNetworkChain = function(){

  if (fileExists(chaincodeIdPath)) {
    chaincodeId = fs.readFileSync(chaincodeIdPath, 'utf8');
    console.log("LOG: Set ChaincodeId: " + chaincodeId);
  }

  console.log("LOG: Configuring Network Chain");
  console.log("Peer address: " + PEER_ADDRESS);
  console.log("Mmberssvc address: " + MEMBERSRVC_ADDRESS);

  // Set Chain
  chain = hfc.newChain('insurance-setup');
  chain.setDevMode(config.blockchain.devMode);

  // Set the URL for membership services
  chain.setMemberServicesUrl("grpc://" + MEMBERSRVC_ADDRESS);

// Add at least one peer's URL.  If you add multiple peers, it will failover
// to the 2nd if the 1st fails, to the 3rd if both the 1st and 2nd fails, etc.
  chain.addPeer("grpc://" + PEER_ADDRESS);
};

/**
 * @sets {caURL, certFile, chain.keyValStore, chain.MemberServicesUrl}
 * @configures {keyValueStore, memberServicesUrl}
 */
var setupCertificates = function(){

  console.log("LOG: Setting Up Certificates");

  /**
   * Configure the KeyValStore which is used to store sensitive keys.
   * This data needs to be located or accessible any time the users enrollmentID
   * perform any functions on the blockchain.  The users are not usable without
   * This data.
   **/

  console.log("KeyValStoreDestination: " + KEYSTORE_PATH);

  chain.setKeyValStore(hfc.newFileKeyValStore(KEYSTORE_PATH));
};


var enrollDeployUser = function() {
  return new Promise(function(resolve, reject){
    var enrollmentId = config.blockchain.setup.deployUser.enrollmentId;
    var enrollSecret = config.blockchain.setup.deployUser.enrollSecret;
    console.log("LOG: Enrolling Deploy User: " + enrollmentId + " - " + enrollSecret);
    chain.enroll(enrollmentId, enrollSecret, function(err, user){
      if (err) {
        console.log("Error enrolling deploy user...maybe already enrolled?");
        console.error(err);
        chain.getUser(enrollmentId, function (err, userFromGet) {
          if (err) {
            console.log("ERROR");
            console.error(err);
            throw Error ("\nError: failed to enroll deploy user " + enrollmentId + ": " + err);
          } else {
            deployUser = userFromGet;
            resolve();
          }
        })
      } else {
        console.log("LOG: Enrolled New User: " + enrollmentId);
        deployUser = user;
        resolve();
      }
    })
  });
}
/**
 * Function to enroll the configured registrar
 * @param user
 */
var enrollRegistrarUser = function(user){
  return new Promise(function(resolve, reject){
    console.log("LOG: Enrolling Registrar: " + user.enrollId + " - " + user.enrollSecret);
    chain.enroll(user.enrollId, user.enrollSecret, function(err, registrarUser){
      if (err) {
        console.log("Error enrolling registrar...maybe already enrolled?");
        console.error(err);
        chain.getUser(user.enrollId, function (err, registrarUserFromGet) {
          if (err) {
            console.log("ERROR");
            console.error(err);
            throw Error ("\nError: failed to enroll registrar user " + user.enrollId + ": " + err);
          }
          resolve(user);
        })
      } else {
        console.log("LOG: Enrolled New User: " + user.enrollId);
        resolve(registrarUser);
      }
    })
  });
};

var enrollNewUsers = function() {
  var promises = [];

  console.log("Enrolling new users");
  usersConfig = config.blockchain.users;
  for (var i = 0; i < usersConfig.length; i++) {
    console.log("*** " + usersConfig[i].enrollmentId);
    var user = {
      "username": usersConfig[i].enrollmentId,
      "affiliation": usersConfig[i].affiliation,
      "attributes": usersConfig[i].attributes
    };

    promises.push(_enrollNewUser(user));
  }

  return Promise.all(promises);
};

/**
 * Enrolls and Sets a Given Network User as the Registrar
 * which is authorised to register other users
 * Generally the Admin user (user[0])
 * @param registeredUser
 */
var configureRegistrar = function(registrarUser){
  console.log("LOG: Configuring Registrar User: " + registrarUser.name);
  chain.setRegistrar(registrarUser);
};

var setupNetwork = function(){

  configureNetworkChain();
  setupCertificates();
  console.log("Completed setupCertificates");
  if (config.blockchain.setup.shouldSetupUsers) {
    return enrollRegistrarUser(USERS[0])
      .then(configureRegistrar)
      .then(function() {
        console.log("Registrar Configured");
      })
      .then(enrollNewUsers)
      .then(function() {
        console.log("New Users Enrolled");
      })
      .then(enrollDeployUser)
      //.then(deployChaincode)
      .then(function() {
        console.log("Setup completed");
      })
      .catch(function (error) {
        throw new Error("Failed to setup network: " + error)
      });
  }
  
  //return deployChaincode();

};

/**
 * Enrolls a new user using provider (non-network user)
 * @param username
 * @param affiliation
 * @param callback
 */
var _enrollNewUser = function(newUser){
  console.log("LOG: enrolling, username: " + newUser.username + ", affiliation: " + newUser.affiliation);

  var registrationRequest = {
    enrollmentID: newUser.username,
    //affiliation: newUser.affiliation
    affiliation: "institution_a"
  };

  if (newUser.attributes) {
    registrationRequest.attributes = newUser.attributes;
  }

  return new Promise(function(resolve, reject){
    chain.registerAndEnroll(registrationRequest, function(err, user){

      // Throw an error if user couldn't be registered
      if (err) {
        console.log("Unable to register user: " + err);
        throw Error("Failed to register and enroll " + newUser.username + ": " + err );
      }

      /**
       * Returns the user so that actions can be performed against them
       */
      console.log("LOG: Successfully enrolled User: " + user.name);
      resolve(user);
    });
  })
};

var deployChaincode = function() {
  return new Promise(function(resolve, reject) {
    if (!chaincodeId) {

      console.log("LOG: Deploying Chaincode");
      chain.setDeployWaitTime(config.blockchain.setup.deploy.waitTime);
      var args = getArgs(config.blockchain.setup.deploy.args);
      console.log(args);
      // Construct the deploy request
      var deployRequest = {
        // Function to trigger
        fcn: config.blockchain.setup.deploy.functionName,
        // Arguments to the initializing function
        args: args,
        chaincodePath: config.blockchain.setup.deploy.chaincodePath
      };

      console.log("--- Deploy Request ---");
      console.log(deployRequest);

      // Trigger the deploy transaction
      var deployTx = deployUser.deploy(deployRequest);

      console.log("after deploy")
      // Print the deploy results
      deployTx.on('complete', function (results) {
        // Deploy request completed successfully
        chaincodeID = results.chaincodeID;
        console.log("\nChaincode ID : " + chaincodeID);
        console.log(util.format("\nSuccessfully deployed chaincode: request=%j, response=%j", deployRequest, results));
        // Save the chaincodeID
        fs.writeFileSync(chaincodeId, chaincodeID);
        //invoke();
        resolve(chaincodeId);
      });

      deployTx.on('error', function (err) {
        // Deploy request failed
        console.log(util.format("\nFailed to deploy chaincode: request=%j, error=%j", deployRequest, err));
        process.exit(1);
      });

      deployTx.on('submitted', function (results) {
        // Deploy request failed
        chaincodeID = results.chaincodeID;
        console.log("\nChaincode ID : " + chaincodeID);
        console.log(util.format("\nSuccessfully submitted chaincode: request=%j, response=%j", deployRequest, results));
        // Save the chaincodeID
        fs.writeFileSync(chaincodeIdPath, chaincodeID);
        //invoke();
      });
    } else {
      console.log("LOG: Already deployed");
      resolve();
    }
  });
};

var enrollNewUser = function(newUser){

  if (!chain.getRegistrar()){
    console.log("Registrar not set...enrolling");
    return enrollRegistrarUser(USERS[0])
      .then(configureRegistrar)
      .then(function(){
        return newUser;
      })
      .then(_enrollNewUser)
      .catch(function(error){
        console.error(error);
      });
  } else {
    return _enrollNewUser(newUser);
  }
};

function getArgs(argsConfig) {
  var args = [];
  for (var i = 0; i < argsConfig.length; i++) {
    args.push(argsConfig[i]);
  }
  return args;
}

function fileExists(filePath) {
  try {
    return fs.statSync(filePath).isFile();
  } catch (err) {
    return false;
  }
}

module.exports = {
  setupNetwork: setupNetwork,
  enrollNewUser: enrollNewUser
};
