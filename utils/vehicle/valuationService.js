var path = require('path');
var util = require("util");
var fs = require('fs');
var request = require("request");
var NodeCache = require("node-cache");
var resultCache = new NodeCache( { stdTTL: 3600, checkperiod: 600 } );
var edmundsApiKey = process.env.EDMUNDS_API_KEY;
var dollarToPoundRate = 0.8;
var defaultCarValue = 6000;

var appDir = path.dirname(require.main.filename);
var API_KEY_PATH = appDir + "/config/credentials/.edmund-api-key";

var getVehicleValuation = function(styleId, mileage, callback){
  var cacheKey = getQueryCacheKey(styleId, mileage);
  var cacheResult = resultCache.get(cacheKey);

  if (cacheResult) {
    console.log("Cache hit on styleId/Mileage: " + styleId + "/" + mileage);
    callback(cacheResult);
    return;
  }

  if (!edmundsApiKey){
    try{
      edmundsApiKey = fs.readFileSync(API_KEY_PATH, 'utf8');
    } catch(e){
      console.error("Could not load API Key From File: ", e);
    }
  }

  if (edmundsApiKey) {
    getVehicleValuationFromEdmunds(styleId, mileage, function (valueInDollars) {
      //Result is in dollars so estimate pound value (Not prod ready but PoC ready!)
      var valueInPounds = valueInDollars * dollarToPoundRate;
      resultCache.set(cacheKey, valueInPounds);

      callback(valueInPounds);
      return;
    });
  } else {
    console.error("Edmunds API key env variable (EDMUNDS_API_KEY) not set");
    console.log("Returning dummy value");
    callback(defaultCarValue);
  }
};

var getQueryCacheKey = function(styleId, mileage) {
  return styleId + ":" + mileage;
}

//Retrieve a car valuation from the Edmunds.com api
//This provides US data so is not accurate for the UK market but will suffice for a PoC.
//The condition is hardcoded to 'Clean' and the zip code is also hard coded.
var getVehicleValuationFromEdmunds = function(styleId, mileage, callback) {

  var baseUrl = "https://api.edmunds.com/v1/api/tmv/tmvservice/calculateusedtmv";
  var params = "?styleid=%s&condition=Clean&mileage=%s&zip=32789&fmt=json&api_key=%s";

  var url = baseUrl + util.format(params, styleId, mileage, edmundsApiKey);

  request.get(baseUrl + util.format(params, styleId, mileage, edmundsApiKey), function (error, response, body) {
    if (!error && response.statusCode == 200) {
      console.log(body);
      var result = JSON.parse(body).tmv;
      var retailValue = parseInt(result.nationalBasePrice.usedTmvRetail);
      var mileageAdjustment = parseInt(result.mileageAdjustment.usedTmvRetail);

      callback(retailValue + mileageAdjustment);
    } else {
      console.error("Unable to retrieve car value: " + response + " " + error);
      console.log("Returning dummy value");
      callback(defaultCarValue);
    }
  });
};

module.exports = {
  getVehicleValuation: function(styleId, mileage, callback){
    getVehicleValuation(styleId, mileage, callback);
  }
};
