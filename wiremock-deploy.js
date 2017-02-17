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
  readFiles().forEach(function (mapping) {
    request('POST', _mappingUrl, {
      headers: {
        'content-type': 'application/json'
      },
      body: JSON.stringify(mapping)
    }).then(function (response) {
      console.log(response.statusCode);
    });
  });
}

run();
