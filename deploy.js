const fetch = require('node-fetch');
const fs = require('fs');

const blockchain = require('./utils/blockchain/blockchain-helpers.js');

const deployMainChaincode = options => {
  let finalResult = {};
  console.log("[DEPLOY] Deploying Main chaincode");
  return blockchain.deploy(options.endpoint,
    "https://github.com/Capgemini-AIE/blockchain-insurance/chaincode/src/insurance_main",
    options.user,
    [options.crudHash])
    .then(result => {
      console.log("[DEPLOY] Chaincode Deployed Successfully");
      // Write hashes to environment variable
      process.env.CHAINCODE_ID = result;
      process.env.CHAINCODE_CRUD_ID = options.crudHash;
      // Read existing config return in and object along with api call result
      return  { chaincodeHash: result, crudHash: options.crudHash };
    }).then(result => {
      finalResult = result;
      // write output
      return new Promise((resolve, reject) => {
        fs.writeFile(__dirname + "/chaincodeIDs.json", JSON.stringify(result), 'utf8', function (err) {
          if (err) {
              reject(err);
          }

          console.log("The file was saved!");
          resolve(result);
        });
      });
    })
    .then(res => finalResult)
    .catch(error => {
      console.error("[DEPLOY] Chaincode deployment failed with error:");
      console.error(error);
    });
}

let setupInitialData = (endpoint, hash) => {
  console.log("[SETUP] Querying existing policies");
  return blockchain.query(endpoint, hash, "claimant1", "retrieveAllPolicies", [])
    .then(response => {
      let result = [];
      try {
        result = JSON.parse(response.result.message);
      } catch (error) {
      }
      console.log("[SETUP] Query finished.");
      if (result.length === 0) {
        console.log("[SETUP] Invoking addPolicy");
        return blockchain.invoke(endpoint, result.chaincodeHash, "insurer1", "addPolicy", [
            "claimant1",
            "2017-02-08",
            "2018-02-08",
            "100",
            "DZ14TYV"
          ])
          .then(result => {
            console.log("[SETUP] Successfully invoked addPolicy");
          });
      }
    });
}

let deploy = (endpoint, user) => {
    if (process.env.CRUD_HASH) {
      console.log("[DEPLOY] CRUD Chaincode already deployed.");
      deployMainChaincode({
        crudHash: process.env.CRUD_HASH,
        user: user.enrollmentId,
        endpoint
      });
    } else {
      console.log("[DEPLOY] Deploying CRUD chaincode");
      blockchain.deploy(endpoint,
        "https://github.com/Capgemini-AIE/blockchain-insurance/chaincode/src/insurance_crud",
        user.enrollmentId,
        [])
        .then(hash => ({
          crudHash: hash,
          user: user.enrollmentId,
          endpoint
        }))
        .then(options => {
          return options;
        })
        .then(deployMainChaincode)
        .then(result => {
          // Create policy
          return setupInitialData(endpoint, result.chaincodeHash);
        })
        .catch(error => {
          console.error("[DEPLOY] Deploy of CRUD chaincode failed with error:");
          console.error(error);
        });
    }
  
}

module.exports = {
  deploy
};