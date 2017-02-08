var config =  require('config');

// Integration
var cloudantIntegration = require('./custom_modules/integration/module_CloudantIntegration');
var auth = require('./custom_modules/integration/module_AuthService');

// Helper Functions
var arrayHelperFunctions = require('./custom_modules/helpers/module_ArrayHelperFunctions');
var objectHelperFunctions = require('./custom_modules/helpers/module_ObjectHelperFunctions');
var routingHelperFunctions = require('./custom_modules/helpers/module_RoutingHelperFunctions');

var blockchainSetup = require('./custom_modules/blockchain/module_BlockchainSetup');

// Server Imports
var express = require('express'), http = require('http'), path = require('path'), fs = require('fs');

// Create Server
var app = express();

blockchainSetup.setup();

/**
 * JSON Schema Validation
 */
var validate = require('express-jsonschema').validate;

// [TODO] Add or modify generated schemas to create custom validation
var schemas = {};

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
app.use(express.static(path.join(__dirname, 'public')));
app.use('/style', express.static(path.join(__dirname, '/views/style')));
app.use(routingHelperFunctions.unlessRoute(["/auth", "/swagger.json"], auth.middleware));
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
 *           $ref: 'tbc'
 *     responses:
 *       200:
 *         description: Successfully created
 */
app.post('/auth/', function(request, response){
  var responseBody = {};

  // [TODO] Change generated 'true' to create custom authentication logic
  if(true){
    auth.signJWT(request.body, {secret:"123", name: "exampleAPI"}, function(resp){
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
    // [TODO] Add in logic for failed authentication
  }
});

/**
 * @swagger
 * /component/test:
 *   get:
 *     tags:
 *       - blockchain-insurance-node-components
 *     description: Is a test endpoint
 *     produces:
 *       - application/json
 *     responses:
 *       200:
 *         description: Successful
 */
app.get('/' + apiPath.base + '/test', function(request, response){
	var responseBody = {};

  blockchainSetup.setup();

	responseBody.message = "Endpoint hit successfully";
	response.setHeader('Content-Type', 'application/json');
	response.write(JSON.stringify(responseBody));
	response.end();
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
