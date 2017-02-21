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
var garageService = require('./utils/blockchain/garageService');

// Server Imports
var express = require('express'), http = require('http'), path = require('path'), fs = require('fs');

// Create Server
var app = express();

var server = http.createServer(app);
socketIntegration.initialise(server);
server.listen(process.env.PORT || 3000);


blockchainSetup.setupNetwork();

/**
 * JSON Schema Validation
 */
var validate = require('express-jsonschema').validate;

// [TODO] Add or modify generated schemas to create custom validation
var schemas = {};
schemas.authSchema = require("./config/schemas/authSchema.json");
schemas.postClaimSchema = require("./config/schemas/postClaimSchema.json");
schemas.postGarageReportSchemas = require('./config/schemas/postGarageReportSchema.json');
schemas.postPayoutAgreementSchema = require('./config/schemas/postPayoutAgreementSchema.json');

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
swaggerSpec.definitions = objectHelperFunctions.deReferenceSchema(swaggerSpec.definitions, require("./config/schemas/postGarageReportSchema.json"), "postGarageReportSchema");
swaggerSpec.definitions = objectHelperFunctions.deReferenceSchema(swaggerSpec.definitions, require("./config/schemas/postPayoutAgreementSchema.json"), "postPayoutAgreementSchema");

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

// Local Only
if ('development' == app.get('env')) {
  app.set('port', process.env.PORT || 3000);
}
app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');
app.engine('html', require('ejs').renderFile);
app.use(logger('dev'));
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.use('/style', express.static(path.join(__dirname, '/views/style')));
app.use(routingHelperFunctions.unlessRoute(["/auth", "/swagger.json","/socket.io/", "/" + apiPath.base + "/oracle*"], auth.middleware));
app.use(auth.allowOriginsMiddleware);



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

    if(rsp.details && rsp.details.carInsurance){//Add following to restrict to claimant: && rsp.details.carInsurance.type === "claimant"){
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
app.get('/' + apiPath.base + '/test/:username', auth.checkAuthorized, function(request, response){
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
 *       - name: post-claim-schema
 *         description: claim content
 *         in: body
 *         required: true
 *         schema:
 *           $ref: '/definitions/postClaimSchema'
 *     responses:
 *       200:
 *         description: Successful
 */
app.post('/claimant/:username/claim', validate({ body: schemas.postClaimSchema}), auth.checkAuthorized, function(request, response){

  var responseBody = {};

  claimService.raiseClaim(request.body, request.params.username, function(res){
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
 * /claimant/{username}/claim/{claimId}/payout/agreement:
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
 *       - name: claimId
 *         description: Id of the claim
 *         in: path
 *         type: string
 *         required: true,
 *       - name: post-payout-agreement-schema
 *         description: agreement
 *         in: body
 *         required: true
 *         schema:
 *           $ref: '/definitions/postPayoutAgreementSchema'
 *     responses:
 *       200:
 *         description: Successful
 */
app.post('/claimant/:username/claim/:claimId/payout/agreement', validate({ body: schemas.postPayoutAgreementSchema}), auth.checkAuthorized, function(request, response){

  var responseBody = {};

  claimService.makeClaimAgreement(request.params.claimId, request.body.agreement, request.params.username, function(res){
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
 * /caller/{username}/garage/{garage}/report:
 *   post:
 *     tags:
 *       - blockchain-insurance
 *     description: Endpoint for submitting garage reports against claims
 *     produces:
 *       - application/json
 *     parameters:
 *       - name: username
 *         description: the username of the person submitting the report
 *         in: path
 *         type: string
 *         required: true,
 *       - name: garageReport
 *         description: the garage report
 *         in: path
 *         type: string
 *         required: true,
 *       - name: post-garage-report-schema
 *         description: claim content
 *         in: body
 *         required: true
 *         schema:
 *           $ref: '/definitions/postGarageReportSchema'
 *     responses:
 *       200:
 *         description: Successful
 */
app.post('/caller/:username/garage/:garage/report', validate({ body: schemas.postGarageReportSchemas}), auth.checkAuthorized, function(request, response){

  var responseBody = {};

  garageService.addGarageReport(request.body, request.params.username, function(res){

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
 * /caller/{username}/history/claims/all:
 *   post:
 *     tags:
 *       - blockchain-insurance
 *     description: Getting all claims authorised for the user
 *     produces:
 *       - application/json
 *     parameters:
 *       - name: username
 *         description: the username of the person submitting the report
 *         in: path
 *         type: string
 *         required: true,
 *     responses:
 *       200:
 *         description: Successful
 */
app.get('/caller/:username/history/claims/all', auth.checkAuthorized, function(request, response){

  var username = request.params.username;

  var responseBody = {};

  claimService.getFullHistory(request.params.username, function(res){

    if (res.error){
      responseBody.error = res.error;
      response.statusCode = 500;
    } else if (res.results){
      responseBody.results = JSON.parse(res.results);
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
