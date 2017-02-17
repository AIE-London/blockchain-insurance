var config = require('config');


function ivoke(functonName, args, user, callback){
  var invokeRequest = {
    chaincodeID: condig.blockchain.chaincodeId,
    fcn: functonName,
    args: args
  };

  console.log("Invoking Requests: ");
  console.log(invokeRequest);

  var tx = user.invoke(invokeRequest);

  // Listen for the 'submitted' event
  tx.on('submitted', function(results) {
    console.log("submitted invoke: %j",results);
    callback();
  });
  // Listen for the 'complete' event.
  tx.on('complete', function(results) {
    console.log("completed invoke: %j",results);
  });
  // Listen for the 'error' event.
  tx.on('error', function(err) {
    console.log("error on invoke: %j",err);
    callback(err)
  });

}

module.exports = {
  invoke: invoke
};
