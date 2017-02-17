/**
 * Created by dcotton on 17/02/2017.
 */

var read = require('fs-readdir-recursive');
var request = require('then-request');

const _mappingUrl = 'http://aston-wiremock.eu-gb.mybluemix.net/__admin/mappings';
const _resetUrl = 'http://aston-wiremock.eu-gb.mybluemix.net/__admin/reset';

function readFiles(){
  return read(__dirname + '/mappings').map(function (filePath){
    return require(__dirname + '/mappings/' + filePath);
  });
}

function run() {
  request('POST', _resetUrl).then(function (response) {
    console.log(response.statusCode);
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
  });
}

run();
