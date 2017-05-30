let fetch = require('node-fetch');

let deployChaincode = (peerEndpoint, chaincodePath, username, args) => {
    let req = {
        "jsonrpc": "2.0",
        "method": "deploy",
        "params": {
            "type": 1,
            "chaincodeID": {
                "path": chaincodePath
            },
            "ctorMsg": {
                "function": "init",
                "args": args
            },
            "secureContext": username
        },
        "id": 1
    };
    console.log(req);
    return fetch(`https://${peerEndpoint}/chaincode`, {
        method: 'POST',
        body: JSON.stringify(req),
        headers: {
            "Content-Type": "application/json"
        }
    })
        .then(response => response.json().then(json => ({
            json,
            status: response.status
        })))
        .then(response => {
            if (response.status !== 200) {
                throw new Error("Response was not 200, has message : " + JSON.stringify(response.json));
            }
            return response.json;
        })
        .then(response => {
            console.log(response)
            return response;
        })
        .then(jsonResponse => jsonResponse.result.message);
}

let queryChaincode = (peerEndpoint, chaincodeHash, username, functionName, functionArgs) => {
    return fetch(`https://${peerEndpoint}/chaincode`, {
        method: 'POST',
        body: JSON.stringify({
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                "type": 1,
                "chaincodeID": {
                    "name": chaincodeHash
                },
                "ctorMsg": {
                    "function": functionName,
                    "args": functionArgs
                },
                "attributes": ["username","role"],
                "secureContext": username
            },
            "id": 0
        }),
        headers: {
            "Content-Type": "application/json"
        }
    })
        .then(response => response.json());
}


let invokeChaincode = (peerEndpoint, chaincodeHash, username, functionName, functionArgs) => {
    return fetch(`https://${peerEndpoint}/chaincode`, {
        method: 'POST',
        body: JSON.stringify({
            "jsonrpc": "2.0",
            "method": "invoke",
            "params": {
                "type": 1,
                "chaincodeID": {
                    "name": chaincodeHash
                },
                "ctorMsg": {
                    "function": functionName,
                    "args": functionArgs
                },
                "attributes": ["username","role"],
                "secureContext": username
            },
            "id": 0
        }),
        headers: {
            "Content-Type": "application/json"
        }
    })
        .then(response => response.json());
}

let enroll = (endpoint, user) => {
    console.log(user.enrollId);
    return fetch(`https://${endpoint}/registrar`, {
        method: 'POST',
        body: JSON.stringify({
            "enrollId": user.enrollId,
            "enrollSecret": user.enrollSecret
        }),
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then(response => {
        return response.json().then(json => {
            if (response.status !== 200) {
                throw new Error("Response was not 200, has message : " + JSON.stringify(json));
            }
            return json;
        });
    })
}

module.exports = {
    deploy: deployChaincode,
    query: queryChaincode,
    invoke: invokeChaincode,
    enroll
};