/**
 * Created by dcotton on 17/02/2017.
 */

var read = require('fs-readdir-recursive');
var request = require('then-request');

const _mappingUrl = 'http://aston-wiremock.eu-gb.mybluemix.net/__admin/mappings';

function readFiles(){
  return read(__dirname + '/mappings').map(function (filePath){
    return require(__dirname + '/mappings/' + filePath);
  });
}

function run() {
  var count = {
    all: 0,
    success: 0
  };
  var files = readFiles();
  files.forEach(function (mapping) {
    request('POST', _mappingUrl, {
      headers: {
        'content-type': 'application/json'
      },
      body: JSON.stringify(mapping)
    }).then(function (response) {
      count.all++;
      if (response.statusCode >= 200 && response.statusCode < 3000) {
        count.success++;
      } else {
        console.error("Failed for mapping: \n \n " + JSON.stringify(mapping) + " \n \n With Response: \n \n" +response);
      }
      if (count.all === files.length - 1 && count.success === files.length - 1) {
        console.log("Successful Deployment");
      } else if (count.all === files.length - 1) {
        console.error("Unsuccessful Deployment");
      }
    });
  });
}

run();
