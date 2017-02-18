var config =  require('config');

// Integration
var cloudantIntegration = require('./utils/integration/cloudantIntegration');

var userService = require('./utils/integration/userService');

var auth = require('./utils/integration/authService');

var socketIntegration = require('./utils/integration/socketIntegration');

// Helper Functions
var arrayHelperFunctions = require('./utils/helpers/array');
var objectHelperFunctions = require('./utils/helpers/object');
var routingHelperFunctions = require('./utils/helpers/routing');

var blockchainSetup = require('./utils/blockchain/setup');

var oracle = require('./utils/blockchain/oracle');
var claimService = require('./utils/blockchain/claimService');

// Server Imports
var express = require('express'), http = require('http'), path = require('path'), fs = require('fs');

// Create Server
var app = express();

blockchainSetup.setupNetwork();

/**
 * JSON Schema Validation
 */
var validate = require('express-jsonschema').validate;

// [TODO] Add or modify generated schemas to create custom validation
var schemas = {};
schemas.authSchema = require("./config/schemas/authSchema.json");
schemas.postClaimSchema = require("./config/schemas/postClaimSchema.json");

/**
 * Swagger Configuration
 */

// Swagger Import
var swaggerJSDoc = require('swagger-jsdoc');

// Swagger Definition
var swaggerDefinition =  config.swagger;

// Options for the swagger docs
var options = {
  swaggerDefinition: swaggerDefinition,
  apis: ['app.js'], // Path to API Documentation
};

// Initialize Swagger-jsdoc
var swaggerSpec = swaggerJSDoc(options);

// Re-use validation-schemas for swagger, but delete unneeded attributes
swaggerSpec.definitions = objectHelperFunctions.deReferenceSchema(swaggerSpec.definitions, require("./config/schemas/authSchema.json"), "authSchema");
swaggerSpec.definitions = objectHelperFunctions.deReferenceSchema(swaggerSpec.definitions, require("./config/schemas/postClaimSchema.json"), "postClaimSchema");

/**
 * Environment Configuration
 */

// Imports
var logger = require('morgan');
var errorHandler = require('errorhandler');
var bodyParser = require('body-parser');

// API Base Path
var apiPath = config.app.paths.api;

// All Environments
app.set('port', process.env.PORT || 3000);
app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');
app.engine('html', require('ejs').renderFile);
app.use(logger('dev'));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use('/style', express.static(path.join(__dirname, '/views/style')));
app.use(routingHelperFunctions.unlessRoute(["/auth", "/swagger.json","/socket.io/", "/" + apiPath.base + "/oracle*"], auth.middleware));
app.use(auth.allowOriginsMiddleware);


// socketIntegration.initialise(8080);

// Development Only
if ('development' == app.get('env')) {
	app.use(errorHandler());
}

/**
 * HTTP REST API ENDPOINTS
 */

// Server Swagger
app.get('/swagger.json', function(req, res) {
  res.setHeader('Content-Type', 'application/json');
  res.send(swaggerSpec);
});

/**
 * @swagger
 * /auth/:
 *   post:
 *     tags:
 *       - ExampleAPI
 *     description: Authenticates and returns token as header
 *     produces:
 *       - application/json
 *     parameters:
 *       - name: JSON
 *         description: JSON body of post
 *         in: body
 *         required: true
 *         schema:
 *           $ref: '#/definitions/authSchema'
 *     responses:
 *       200:
 *         description: Successfully created
 */
app.post('/auth/', validate({ body : schemas.authSchema }), function(request, response){
  var responseBody = {};

  var username = request.body.username;
  var password = request.body.password;

  userService.authenticate(username, password, function(rsp){
    if(rsp.details && rsp.details.carInsurance && rsp.details.carInsurance.type === "claimant"){
      auth.signJWT(request.body, {secret:"123", name: username}, function(resp){
        console.log(resp);
        if(resp){
          response.setHeader('Content-Type', 'application/json');
          response.setHeader('Token', resp);
          response.write(JSON.stringify(responseBody));
          response.end();
          return;
        }
        else{
          responseBody.reason = "Server error";
          response.setHeader('Content-Type', 'application/json');
          response.statusCode = 500;
          response.write(JSON.stringify(responseBody));
          response.end();
          return;
        }
      });
    }else{
      responseBody.reason = "Incorrect Credentials";
      response.setHeader('Content-Type', 'application/json');
      response.statusCode = 401;
      response.write(JSON.stringify(responseBody));
      response.end();
      return;
    }
  });
});

/**
 * @swagger
 * /component/test/{username}:
 *   get:
 *     tags:
 *       - blockchain-insurance-node-components
 *     description: Is a test endpoint
 *     produces:
 *       - application/json
 *     parameters:
 *       - name: username
 *         description: the username
 *         in: path
 *         type: string
 *         required: true
 *     responses:
 *       200:
 *         description: Successful
 */
app.get('/' + apiPath.base + '/test/:username', function(request, response){
	var responseBody = {};

  blockchainSetup.setup();

	responseBody.message = "Endpoint hit successfully";
	response.setHeader('Content-Type', 'application/json');
	response.write(JSON.stringify(responseBody));
	response.end();
	return;
});

/**
 * @swagger
 * /claimant/{username}/claim:
 *   post:
 *     tags:
 *       - blockchain-insurance-node-components
 *     description: Is a test endpoint
 *     produces:
 *       - application/json
 *     parameters:
 *       - name: username
 *         description: the username
 *         in: path
 *         type: string
 *         required: true,
 *         schema:
 *           $ref: '/definitions/postClaimSchema'
 *     responses:
 *       200:
 *         description: Successful
 */
app.post('/claimant/:username/claim', validate({ body: schemas.postClaimSchema}), function(request, response){

  var responseBody = {};

  claimService.raiseClaim(request.body, function(res){
    if (res.error){
      responseBody.error = res.error;
      response.statusCode = 500;
    } else if (res.results){
      responseBody.results = res.results;
      response.statusCode = 200;
    } else {
      responseBody.error = "unknown issue";
      response.statusCode = 500;
    }

    response.setHeader('Content-Type', 'application/json');
    response.write(JSON.stringify(responseBody));
    response.end();
    return;

  });
});

/**
 * @swagger
 * /component/oracle/vehicle/{styleId}/value:
 *   get:
 *     tags:
 *       - blockchain-insurance-node-components
 *     description: Obtain an estimated vehicle value based on an Edmunds Api style id
 *     produces:
 *       - application/json
 *     parameters:
 *       - styleId: styleId
 *         description: the edmunds api style id of the vehicle
 *         in: path
 *         type: string
 *         required: true
 *       - mileage: mileage
 *         description: the mileage of the vehicle
 *         in: query
 *         type: string
 *         required: true
 *       - requestId: requestId
 *         description: requests with the same requestId will always return the same result
 *         in: query
 *         type: string
 *         required: true
 *       - callbackFunctionName: callbackFunctionName
 *         description: the name of the chaincode function that should be invoked when a value has been obtained
 *         in: query
 *         type: string
 *         required: true
 *     responses:
 *       202:
 *         description: Successful ???
 */
app.get('/' + apiPath.base + '/oracle/vehicle/:styleId/value', function(request, response){
  var responseBody = {};

  var mileage = request.query.mileage;
  var requestId = request.query.requestId;
  var callbackFunctionName = request.query.callbackFunction;
  var styleId = request.params.styleId;

  oracle.requestVehicleValuation(requestId, styleId, mileage, callbackFunctionName, function() {
    response.setHeader('Content-Type', 'application/json');
    response.write(JSON.stringify(responseBody));
    response.statusCode = 202;
    response.end();
    return;
  });

  return;

});

 // End API

// Express Error Handling for Validation
app.use(function(err, req, res, next) {
  var responseData;

  if (err.name === 'JsonSchemaValidation') {
    // Log the error however you please
    console.log(err.message);
    // logs "express-jsonschema: Invalid data found"

    // Set a bad request http response status or whatever you want
    res.status(400);

    // Format the response body however you want
    /*responseData = {
     statusText: 'Bad Request',
     jsonSchemaValidation: true,
     validations: err.validations  // All of your validation information
     };*/
    responseData = {};
    responseData.reason = "Bad Request";

    res.json(responseData);
  } else {
    // pass error to next error middleware handler
    next(err);
  }
});
app.use(bodyParser.json());
http.createServer(app).listen(app.get('port'), '0.0.0.0', function() {
	console.log('Express server listening on port ' + app.get('port'));
});
