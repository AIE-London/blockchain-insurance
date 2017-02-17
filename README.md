# Blockchain Insurance

A proof of concept Blockchain project around Car Insurance.

## Development Environment

We're building on the docker-compose from [https://github.com/IBM-Blockchain/fabric-images](https://github.com/IBM-Blockchain/fabric-images).

See also [these instructions for using a local hyperledger network](https://github.com/IBM-Blockchain/marbles/blob/master/docs/use_local_hyperledger.md).

To run a single node network using Docker Compose -

```
docker-compose -f docker-compose/single-peer-ca.yaml up -d
```

For a four node + CA network:

```
docker-compose -f four-peer-ca.yaml up -d
```

To test it is up and running properly hit [http://localhost:7050/chain](http://localhost:7050/chain)

### Deploying chaincode for development (Only currently working in the 4 peer environment)

To run chaincode locally without having to deploy from a github url:

1) Build your chaincode

```
go build ./
```

2) Register the chaincode with a peer

```
CORE_CHAINCODE_ID_NAME=insurance CORE_PEER_ADDRESS=0.0.0.0:7051 ./insurance_main
```

3) Send a REST request to enroll a user

```
POST: http://localhost:7050/registrar
{
  "enrollId": "bob",
  "enrollSecret": "NOE63pEQbL25"
}
```

4) Send a REST deploy request

```
POST: http://localhost:7050/chaincode
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name": "insurance"
    },
    "ctorMsg": {
        "function":"init"
    },
    "secureContext": "bob"
  },
  "id": 1
}
```

5) The chaincode should now be deployed and ready to accept INVOKE or QUERY requests.

### Running Blockchain Explorer

If you want to view your blockchian locally you can use the blockchain explorer.

To build the Docker image, run -

```
git clone https://github.com/hyperledger/blockchain-explorer.git
cd blockchain-explorer/explorer_1
docker build -t hyperledger/blockchain-explorer .
```

To run it first you will need to find out which network your hyperledger fabric peers are running on.
To do this run:

```
docker network ls
```

Grab the Docker Compose network name and use it below:

```
docker run -p 9090:9090 -itd --network=dockercompose_default -e HYP_REST_ENDPOINT=http://vp0:7050 hyperledger/blockchain-explorer
```

The explorer should then be available on [http://localhost:9090](http://localhost:9090).

### Configuring the vehicle value oracle

In order to use the car value oracle to obtain actual vehicle values from Edmunds you must:

1) Obtain an Edmunds api key and set the env variable. (Note: this step can be skipped, and the oracle service will then just callback with a hardcoded value.)

```
EDMUNDS_API_KEY
```

2) Deploy the chain code with an init argument specifying the host address of the oracle service (including port)

3) Run the app.js node app

```
npm install
node app.js
```
