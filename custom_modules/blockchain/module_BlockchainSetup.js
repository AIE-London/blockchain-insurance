// Import version to dynamically grab latest chaincode
var version =  require('config').app.version;
var fullVersion = version.major + "." + version.minor + "." + version.bug;

// Import SDK
var Ibc1 = require('ibm-blockchain-js');
var ibc = new Ibc1();

// peers and users
var peers = [];
var users = [];

// Options for loading SDK
var options = {};

// helper functions
/**
 * Checks if the peer wants tls or not
 * @param peer_array
 * @returns {boolean}
 */
var detect_tls_or_not = function(peer_array){
  var tls = false;
  if(peer_array[0] && peer_array[0].api_port_tls){
    if(!isNaN(peer_array[0].api_port_tls)) tls = true;
  }
  return tls;
};

//Chaincode
var chaincode = {
  zip_url: "https://github.com/Capgemini-AIE/blockchain-insurance/archive/" + fullVersion + ".zip",
  unzip_dir: "blockchain-insurance-" + fullVersion + "/chaincode/src/insurance_code/",
  git_url: "http://gopkg.in/Capgemini-AIE/blockchain-insurance.v" + version.major +"/chaincode"
};

var loadCredentials = function(){
  if(process.env.VCAP_SERVICES){
    var servicesObject =- JSON.parse(process.env.VCAP_SERVICES);
    for (var i in servicesObject){
      if(i.indexOf('ibm-blockchain') >= 0){													//looks close enough
        if(servicesObject[i][0].credentials.error){
          console.log('!\n!\n! Error from Bluemix: \n', servicesObject[i][0].credentials.error, '!\n!\n');
          peers = null;
          users = null;
          process.error = {type: 'network', msg: 'Due to overwhelming demand the IBM Blockchain Network service is at maximum capacity.  Please try recreating this service at a later date.'};
        }
        if(servicesObject[i][0].credentials && servicesObject[i][0].credentials.peers){		//found the blob, copy it to 'peers'
          console.log('overwritting peers, loading from a vcap service: ', i);
          peers = servicesObject[i][0].credentials.peers;
          if(servicesObject[i][0].credentials.users){										//user field may or maynot exist, depends on if there is membership services or not for the network
            console.log('overwritting users, loading from a vcap service: ', i);
            users = servicesObject[i][0].credentials.users;
          }
          else users = null;																//no security
          break;
        }
      }
    }
  } else {
    var manual = require('../../config/credentials/ibm_blockchain_credentials.json');

    if (manual) {
      if (manual.credentials.peers) {
        peers = manual.credentials.peers;
      }
      if (manual.credentials.users) {
        users = manual.credentials.users;
      }

      peers = manual.credentials.peers;
      users = manual.credentials.users;
    } else {
      console.error("You do not have a local copy of ibm_blockchain_credentials, please ask an admin for this file");
    }
  }
};

var setOptions = function(){
  options = {
    network:{
      peers: [peers[0]],
      users: [users[0]],
      options: {
        quiet: true,
        tls: detect_tls_or_not(peers),
        maxRetry: 1
      }
    },
    chaincode: chaincode
  };
};

var sdkLoaded = function(err, cc){
  console.log("--error--");
  console.log(err);

  if (err || !cc) {
    console.log("There was an error loading the chaincode / network");
    console.error(err);
  } else {
    console.log("There was no error loading the chaincode / network");
    chaincode = cc;

    console.log("---- CC ----");
    console.log(cc);

    // ---- To Deploy or Not to Deploy ---- //
    if(!cc.details.deployed_name || cc.details.deployed_name === ''){
      cc.deploy('init', ['99'], {delay_ms: 30000}, function(e){ 						// delay_ms is milliseconds to wait after deploy for conatiner to start, 50sec recommended
        check_if_deployed(e, 1);
      });
    } else {
      console.log('chaincode summary file indicates chaincode has been previously deployed');
      check_if_deployed(null, 1);
    }
  }
};

var loadSDK = function(){
  if(process.env.VCAP_SERVICES){
    console.log('\n[!] looks like you are in bluemix, I am going to clear out the deploy_name so that it deploys new cc.\n[!] hope that is ok budddy\n');
    options.chaincode.deployed_name = '';
  };

  console.log(users);
  console.log(peers);


  console.log(options.peers);
  console.log(options.users);
  console.log(options.chaincode)

  if (options.chaincode) {
    ibc.load(options, sdkLoaded());
  } else {
    console.error("Didn't try load SDK because there is no valid chaincode");
  }
};

module.exports = {
  setup: function(){
    loadCredentials();
    setOptions();
    loadSDK();
  }
};
