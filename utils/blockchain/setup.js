process.env.GOPATH = __dirname;

var path = require('path');
var appDir = path.dirname(require.main.filename);

var fs = require('fs');
var config = require('config');
var hfc = require('hfc');

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
var certPath = appDir + "/" + config.blockchain.deployRequest.chaincodePath + "/certificate.pem";


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
var chaincodeId = config.blockchain.chaincodeId;

/**
 * Path of the deployed chaincode
 * @type {string}
 */
var chaincodePath = __dirname + "/" + chaincodeId + "/";

/** NETWORK **/
var network = {};

/** FUNCTIONS **/

/**
 * @sets {chain, network}
 */
var configureNetworkChain = function(){

  console.log("LOG: Confoiguring Network Chain");

  // Set Chain
  chain = hfc.newChain(config.chainName);

  //[TODO] VCAP CREDENTIALS

  // Set Network
  network.credentials = require("../../config/credentials/ServiceCredentials.json");

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
  chain.setKeyValStore(hfc.newFileKeyValStore(__dirname + '/keyValStore-' + uuid));

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
  var cert = fs.readFileSync(certFileName);

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
    console.log("LOG: Enrolling New User: " + user.enrollId);
    chain.enroll(user.enrollId, user.enrollSecret, function(err, networkUser){
      if (err) throw Error ("\nError: failed to enroll network user " + user.enrollId + ": " + err);
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
  console.log("LOG: Configuring Registrar User");
  chain.setRegistrar(registrarUser);
};

var setupNetwork = function(){

  configureNetworkChain();
  setupCertificates();

    return enrollNetworkUser(network.credentials.users[0])
    .then(configureRegistrar)
    .catch(function(error){
      console.error(error);
    }).then(function(){
        console.log("Completed Network Setup");
    });


};

/**
 * Enrolls a new user using provider (non-network user)
 * @param username
 * @param affiliation
 * @param callback
 */
var enrollNewUser = function(username, affiliation, callback){
  var registrationRequest = {
    enrolmentID: username,
    affiliation: affiliation
  };
  return new Promise(function(){
    chain.registerAndEnroll(registrationRequest, function(err, user){

      // Throw an error if user couldn't be registered
      if (err) throw Error("Failed to register and enroll " + username + ": " + err );

      /**
       * Returns the user so that actions can be performed against them
       */
      resolve(user);
    });
  })
};


module.exports = {
  setupNetwork: setupNetwork,
  enrollNewUser: enrollNetworkUser
};
