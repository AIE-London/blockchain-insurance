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