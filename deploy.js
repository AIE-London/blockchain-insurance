const fetch = require('node-fetch');

const blockchain = require('./utils/blockchain/blockchain-helpers.js');

const deployMainChaincode = options => {
  console.log("[DEPLOY] Deploying Main chaincode");
  return blockchain.deploy(options.endpoint,
    "https://github.com/Capgemini-AIE/blockchain-insurance/chaincode/src/insurance_main",
    options.user,
    [options.crudHash])
    .then(result => {
      console.log("[DEPLOY] Chaincode Deployed Successfully");
      // Read existing config return in and object along with api call result
      return  { chaincodeHash: result.hash, crudHash: options.crudHash };
    }).catch(error => {
      console.error("[DEPLOY] Chaincode deployment failed with error:");
      console.error(error);
    });
}

let deploy = (endpoint, user) => {
  blockchain.enroll(endpoint, user).then(() => {
    if (process.env.CRUD_HASH) {
      deployMainChaincode({
        crudHash: process.env.CRUD_HASH,
        user: user.enrollId,
        endpoint
      });
    } else {
      console.log("[DEPLOY] Deploying CRUD chaincode");
      blockchain.deploy(endpoint,
        "https://github.com/Capgemini-AIE/blockchain-insurance/chaincode/src/insurance_crud",
        user.enrollId,
        [])
        .then(hash => ({
          crudHash: hash,
          user: user.enrollId,
          endpoint
        }))
        .then(deployMainChaincode)
        .catch(error => {
          console.error("[DEPLOY] Deploy of CRUD chaincode failed with error:");
          console.error(error);
        });
    }
  });
  
}

module.exports = {
  deploy
};