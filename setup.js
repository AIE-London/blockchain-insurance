var path = require('path');
var util = require('util');
var appDir = path.dirname(require.main.filename);
process.env.GOPATH = appDir;

var fs = require('fs');
var config = require('config');
var hfc = require('hfc');

var cert;

/** CERTIFICATES **/

/**
 * Certificate Authority URL
 * @type {string}
 */
var caUrl;

/**
 * Path of the certificate.pem
 * @type {string}
 */
var certPath = __dirname + "/src/" + config.blockchain.deployRequest.chaincodePath + "/certificate.pem";


/**
 * File Name of the Certificate
 * @type {string}
 */
var certFileName = 'us.blockchain.ibm.com.cert';

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

var currentUsername = config.blockchain.user.username;

/** NETWORK **/
var network = {};
var userObj;


/** FUNCTIONS **/

/**
 * @sets {chain, network}
 */
var configureNetworkChain = function(){

  if (fileExists(chaincodeIdPath)) {
    chaincodeId = fs.readFileSync(chaincodeIdPath, 'utf8');
    console.log("LOG: Set ChaincodeId: " + chaincodeId);
  }

  console.log("LOG: Confoiguring Network Chain");

  // Set Chain
  chain = hfc.newChain(config.blockchain.chainName);

  //[TODO] VCAP CREDENTIALS

  // Set Network
  network.credentials = require("./config/credentials/ServiceCredentials.json");

};

var setupUser = function(){
  console.log("LOG: Setting Up User: " + currentUsername);
  return new Promise(function(resolve, reject){
      chain.getUser(currentUsername, function(err, user) {
        if (err) console.log("err");
        if (err) throw Error(" Failed to register and enroll " + currentUsername + ": " + err);
        userObj = user;
        console.log("Hello?");
        resolve(user);
      });
  });
};

/**
 * @sets {caURL, certFile, chain.keyValStore, chain.MemberServicesUrl}
 * @configures {keyValueStore, memberServicesUrl}
 */
var setupCertificates = function(){

  console.log("LOG: Setting Up Certificates");
  /**
   * Determining if we are running on a startup or HSBN (High Security Business Network)
   * based on the url of the discovery host name.  The HSBN will contain the string zone.
   * @type {boolean}
   */
  var isHSBN = network.credentials.peers[0].discovery_host.indexOf('secure') >= 0 ? true : false;

  /**
   * The network ID is an attribute within credentials under ca.
   * We only expect there to be one so take the first
   */

  var network_id = Object.keys(network.credentials.ca)[0];

  /**
   * Generateing the URL for the Certificate Authority
   * @type {string}
   */
  caUrl = "grpcs://" + network.credentials.ca[network_id].discovery_host + ":" + network.credentials.ca[network_id].discovery_port;

  /**
   * Configure the KeyValStore which is used to store sensitive keys.
   * This data needs to be located or accessible any time the users enrollmentID
   * perform any functions on the blockchain.  The users are not usable without
   * This data.
   **/

    // uuid is always first 8 characters
  var uuid = network_id.substring(0, 8);
  console.log("KeyValStoreDestination: " + appDir + "/KeyValStore-" + uuid);

  chain.setKeyValStore(hfc.newFileKeyValStore(appDir + '/keyValStore-' + uuid));

  /**
   * if it's a High Security Business Network (HSBN),
   * overwrite the certFile path.
   */
  if (isHSBN){
    certFileName = '0.secure.blockchain.ibm.com.cert';
  }

  /**
   * Reads in the relevant certificate and adds it to the chaincode directory
   */

  console.log(appDir);

  console.log("certFileName: " + certFileName);
  console.log("certPath: " + certPath);
  fs.createReadStream(certFileName).pipe(fs.createWriteStream(certPath));
  cert = fs.readFileSync(certFileName);

  /**
   * sets the chain's mmeber service url
   */
  chain.setMemberServicesUrl(caUrl, {
    pem: cert
  });

};

/**
 * Function to enroll a user from the network credentials (i.e. user[0] for admin)
 * @param user
 */
var enrollNetworkUser = function(user){
  return new Promise(function(resolve, reject){
    console.log("LOG: Enrolling New User: " + user.enrollId + " - " + user.enrollSecret);
    chain.enroll(user.enrollId, user.enrollSecret, function(err, networkUser){
      if (err) {
        console.log("ERROR");
        console.error(err);
        throw Error ("\nError: failed to enroll network user " + user.enrollId + ": " + err);
      }
      console.log("LOG: Enrolled New User: " + user.enrollId);
      resolve(networkUser);
    })
  });

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

var addPeersToChain = function(){
  console.log("LOG: addPeersToChain");
  network.credentials.peers.forEach(function(peer){
    var peerUrl = "grpcs://" + peer.discovery_host + ":" + peer.discovery_port;
    console.log("Adding peer: " + peerUrl);
    chain.addPeer(peerUrl, {
      pem: cert
    });

    var eventUrl = "grpcs://" + peer.event_host + ":" + peer.event_port;
    console.log("Adding event: " + eventUrl);
    chain.eventHubConnect(eventUrl, {
      pem: cert
    });
  });
  // Make sure disconnect the eventhub on exit
  process.on('exit', function() {
    chain.eventHubDisconnect();
  });
};

var setupNetwork = function(){

  configureNetworkChain();
  setupCertificates();

    return setupUser()
      .then(function(){return network.credentials.users[0]})
      .then(enrollNetworkUser)
      .then(configureRegistrar)
      .then(addPeersToChain)
      .then(function(){
        return {
          "username": config.blockchain.user.username,
          "affiliation": config.blockchain.user.affiliation
        };
    })
      .then(enrollNewUser)
      .then(function(user){
        userObj = user;
        console.log("Completed Network Setup");
       deployChaincode();
    })
      .catch(function(error){
        throw new Error("Failed to setup network: " + error)
      });

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
    affiliation: newUser.affiliation
  };

  return new Promise(function(resolve, reject){
    chain.registerAndEnroll(registrationRequest, function(err, user){

      // Throw an error if user couldn't be registered
      if (err) throw Error("Failed to register and enroll " + newUser.username + ": " + err );

      /**
       * Returns the user so that actions can be performed against them
       */
      console.log("LOG: Successfully Retrieved User: " + user.name);
      resolve(user);
    });
  })
};

var deployChaincode = function() {
  if (!chaincodeId) {

    console.log("LOG: Deploying Chaincode");
    chain.setDeployWaitTime(config.blockchain.deployWaitTime);
    var args = getArgs(config.blockchain.deployRequest);
    // Construct the deploy request
    var deployRequest = {
      // Function to trigger
      fcn: config.blockchain.deployRequest.functionName,
      // Arguments to the initializing function
      args: args,
      chaincodePath: config.blockchain.deployRequest.chaincodePath,
      // the location where the startup and HSBN store the certificates
      certificatePath: network.credentials.cert_path
    };

    console.log("--- Deploy Request ---");
    console.log(deployRequest);


    // Trigger the deploy transaction
    var deployTx = userObj.deploy(deployRequest);

    // Print the deploy results
    deployTx.on('complete', function (results) {
      // Deploy request completed successfully
      chaincodeID = results.chaincodeID;
      console.log("\nChaincode ID : " + chaincodeID);
      console.log(util.format("\nSuccessfully deployed chaincode: request=%j, response=%j", deployRequest, results));
      // Save the chaincodeID
      fs.writeFileSync(chaincodeId, chaincodeID);
      //invoke();
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
  }
};

var enrollNewUser = function(newUser){

  if (!chain.getRegistrar()){
    return enrollNetworkUser(network.credentials.users[0])
      .then(setupUser)
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

function getArgs(request) {
  var args = [];
  for (var i = 0; i < request.args.length; i++) {
    args.push(request.args[i]);
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
