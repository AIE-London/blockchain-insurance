#!/bin/bash
printf "*** Registering User bob ***\n"
curl -H "Content-Type: application/json" -X POST -d '{"enrollId": "garage1","enrollSecret": "123456789123"}' http://localhost:7050/registrar
sleep 1
curl -H "Content-Type: application/json" -X POST -d '{"enrollId": "claimant1","enrollSecret": "123456789123"}' http://localhost:7050/registrar
sleep 1
curl -H "Content-Type: application/json" -X POST -d '{"enrollId": "claimant2","enrollSecret": "123456789123"}' http://localhost:7050/registrar
sleep 1
curl -H "Content-Type: application/json" -X POST -d '{"enrollId": "insurer1","enrollSecret": "123456789123"}' http://localhost:7050/registrar
sleep 1
curl -H "Content-Type: application/json" -X POST -d '{"enrollId": "insurer2","enrollSecret": "123456789123"}' http://localhost:7050/registrar
sleep 1


printf "\n*** Deploying CRUD chaincode ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name": "insurancecrud"
    },
    "ctorMsg": {
        "function":"init"
    },
    "secureContext": "insurer1"
  },
  "id": 1
}' http://localhost:7050/chaincode
sleep 3

#"args": ["localhost:3000"]
printf "\n*** Deploying main chaincode ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name": "insurance"
    },
    "ctorMsg": {
        "function":"init",
        "args": ["insurancecrud"]
    },
    "secureContext": "insurer1"
  },
  "id": 1
}' http://localhost:7050/chaincode
sleep 3

printf "\n*** Creating policy ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "addPolicy",
             "args": [
                 "claimant1",
                 "2017-02-08",
                 "2018-02-08",
                 "100",
                 "BP08BRV"
             ]
         },
         "secureContext": "insurer1",
				 "attributes": ["username","role"]
     },
     "id": 2
 }' http://localhost:7050/chaincode
sleep 3
 
printf "\n*** Creating claim ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "createClaim",
             "args": [
                 "P3",
                 "I crashed into a tree",
                 "2017-02-10",
                 "multiple_parties",
                 "DZ14TYV",
                 "false"
             ]
         },
         "secureContext": "claimant1",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Declare Liability ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "declareLiability",
             "args": [
                 "C2",
                 "true"
             ]
         },
         "secureContext": "claimant2",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Adding garage report for C1 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "addGarageReport",
             "args": [
                 "C1",
                 "15000",
                 "false",
                 "Half a tree sticking out of the engine",
                 "BP08BRV"
             ]
         },
         "secureContext": "garage1",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Adding garage report for C2 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "addGarageReport",
             "args": [
                 "C2",
                 "12000",
                 "false",
                 "In a proper mess",
                 "DZ14TYV"
             ]
         },
         "secureContext": "garage1",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Agreeing to valuation for C1 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "agreePayoutAmount",
             "args": [
                 "C1",
                 "true"
             ]
         },
         "secureContext": "claimant1",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Agreeing to valuation for C2 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "agreePayoutAmount",
             "args": [
                 "C2",
                 "true"
             ]
         },
         "secureContext": "claimant1",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Confirming Paid C1 - 1 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "confirmPaidOut",
             "args": [
                 "C1",
                 "1"
             ]
         },
         "secureContext": "insurer1",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Confirming Paid C2 - 1 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "confirmPaidOut",
             "args": [
                 "C2",
                 "1"
             ]
         },
         "secureContext": "insurer2",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Confirming Paid C2 - 2 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "confirmPaidOut",
             "args": [
                 "C2",
                 "2"
             ]
         },
         "secureContext": "insurer2",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Close Claim C1 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "closeClaim",
             "args": [
                 "C1"
             ]
         },
         "secureContext": "insurer1",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

printf "\n*** Close Claim C2 ***\n"
curl -H "Content-Type: application/json" -X POST -d '{
     "jsonrpc": "2.0",
     "method": "invoke",
     "params": {
         "type": 1,
         "chaincodeID": {
             "name": "insurance"
         },
         "ctorMsg": {
             "function": "closeClaim",
             "args": [
                 "C2"
             ]
         },
         "secureContext": "insurer2",
				 "attributes": ["username","role"]
     },
     "id": 3
 }' http://localhost:7050/chaincode
sleep 3

 printf "\n*** Retrieving claims for claimant1 ***\n"
 curl -H "Content-Type: application/json" -X POST -d '{
      "jsonrpc": "2.0",
      "method": "query",
      "params": {
          "type": 1,
          "chaincodeID": {
              "name": "insurance"
          },
          "ctorMsg": {
              "function": "retrieveAllClaims"
          },
          "secureContext": "claimant1",
 				 "attributes": ["username","role"]
      },
      "id": 3
  }' http://localhost:7050/chaincode

   printf "\n*** Retrieving claims for claimant2***\n"
   curl -H "Content-Type: application/json" -X POST -d '{
        "jsonrpc": "2.0",
        "method": "query",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": "insurance"
            },
            "ctorMsg": {
                "function": "retrieveAllClaims"
            },
            "secureContext": "claimant2",
   				 "attributes": ["username","role"]
        },
        "id": 3
    }' http://localhost:7050/chaincode






