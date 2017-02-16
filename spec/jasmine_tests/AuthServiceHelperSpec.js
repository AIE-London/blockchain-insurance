describe("AuthServiceHelper", function() {
  authHelper = require('../../utils/integration/authService');


  beforeEach(function () {

  });
  /**
   * Tests running against auth middleware
   */
  describe("middleware should authenticate", function () {
    it("successfully for http://aston-swagger-ui.eu-gb.mybluemix.net (over http)", function () {
      /**
       * Swagger UI Passes validation
       */
      var req = {
        get: function(str){
          return "";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "http://aston-swagger-ui.eu-gb.mybluemix.net"
        },
        method: 'POST'
      };
      var resp = {
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
    it("successfully for https://aston-swagger-ui.eu-gb.mybluemix.net (over https)", function () {
      /**
       * Swagger UI Passes validation
       */
      var req = {
        get: function(str){
          return "";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "https://aston-swagger-ui.eu-gb.mybluemix.net"
        },
        method: 'POST'
      };
      var resp = {
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
    it("successfully for https://random-url.com (over https with valid token)", function () {
      /**
       * With token any url passes validation
       */
      var req = {
        get: function(str){
          return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0ODA0Mzk3MDN9.B4p-G5vMimHKcl61dXOQvO5zLgKmnzGplM6YKk38UIo";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "https://random-url.com"
        },
        method: 'POST'
      };
      var resp = {
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
    it("successfully for http://random-url.com (over http with valid token)", function () {
      /**
       * Over http any url passes with valid token
       */
      var req = {
        get: function(str){
          return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0ODA0Mzk3MDN9.B4p-G5vMimHKcl61dXOQvO5zLgKmnzGplM6YKk38UIo";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "http://random-url.com"
        },
        method: 'POST'
      };
      var resp = {
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
    it("unsuccessfully for https://random-url.com (over https with invalid bearer token)", function () {
      /**
       * Any url fails auth without valid token
       */
      var req = {
        get: function(str){
          return "Bearer JhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0ODA0Mzk3MDN9.B4p-G5vMimHKcl61dXOQvO5zLgKmnzGplM6YKk38UIo";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "https://random-url.com"
        },
        method: 'POST'
      };
      var resp = {
        end: function(){
          expect(true).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(false).toBeTruthy();
      });
    });
    it("successfully for OPTIONS https://random-url.com (over https with invalid bearer token)", function (onComplete) {
      /**
       * Any url fails auth without valid token
       */
      var req = {
        get: function(str){
          return "Bearer JhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0ODA0Mzk3MDN9.B4p-G5vMimHKcl61dXOQvO5zLgKmnzGplM6YKk38UIo";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "https://random-url.com"
        },
        method: 'OPTIONS'
      };
      var resp = {
        end: function(){
          expect(true).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        // we got called back
        expect(true).toBeTruthy();
        onComplete();
      });
    });
    it("unsuccessfully for http://random-url.com (over http with invalid bearer token)", function () {
      /**
       * Any URL fails auth without valid token (over http)
       */
      var req = {
        get: function(str){
          return "Bearer JhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0ODA0Mzk3MDN9.B4p-G5vMimHKcl61dXOQvO5zLgKmnzGplM6YKk38UIo";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "http://random-url.com"
        },
        method: 'POST'
      };
      var resp = {
        end: function(){
          expect(true).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(false).toBeTruthy();
      });
    });
  });
  describe("cross-origin middleware should return header", function () {
    it("successfully for http://aston-swagger-ui.eu-gb.mybluemix.net (over http)", function () {
      /**
       * Swagger UI Passes CORS
       */
      var req = {
        get: function(str){
          return "";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "http://aston-swagger-ui.eu-gb.mybluemix.net"
        },
        method: 'POST'
      };
      var resp = {
        setHeader: function(key, value) {
          if (key === 'Access-Control-Allow-Origin') {
            expect(value).toEqual('http://aston-swagger-ui.eu-gb.mybluemix.net');
          }
          if (key === 'Access-Control-Allow-Methods') {
            expect(value).toEqual('GET, POST, OPTIONS, PUT, PATCH, DELETE');
          }
          if(key === 'Access-Control-Allow-Headers') {
            expect(value).toEqual('X-Requested-With,content-type');
          }
        },
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
    it("successfully for http://aston-swagger-ui.eu-gb.mybluemix.net (over https)", function () {
      /**
       * Swagger UI Passes CORS
       */
      var req = {
        get: function(str){
          return "";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "https://aston-swagger-ui.eu-gb.mybluemix.net"
        },
        method: 'POST'
      };
      var resp = {
        setHeader: function(key, value) {
          if (key === 'Access-Control-Allow-Origin') {
            expect(value).toEqual('https://aston-swagger-ui.eu-gb.mybluemix.net');
          }
          if (key === 'Access-Control-Allow-Methods') {
            expect(value).toEqual('GET, POST, OPTIONS, PUT, PATCH, DELETE');
          }
          if(key === 'Access-Control-Allow-Headers') {
            expect(value).toEqual('X-Requested-With,content-type');
          }
        },
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
    it("successfully for http://localhost:8080 (over http)", function () {
      /**
       * Swagger UI Passes CORS
       */
      var req = {
        get: function(str){
          return "";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "http://localhost:8080"
        },
        method: 'POST'
      };
      var resp = {
        setHeader: function(key, value) {
          if (key === 'Access-Control-Allow-Origin') {
            expect(value).toEqual('http://localhost:8080');
          }
          if (key === 'Access-Control-Allow-Methods') {
            expect(value).toEqual('GET, POST, OPTIONS, PUT, PATCH, DELETE');
          }
          if(key === 'Access-Control-Allow-Headers') {
            expect(value).toEqual('X-Requested-With,content-type');
          }
        },
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
    it("successfully for https://localhost:8080 (over https)", function () {
      /**
       * Swagger UI Passes CORS
       */
      var req = {
        get: function(str){
          return "";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "https://localhost:8080"
        },
        method: 'POST'
      };
      var resp = {
        setHeader: function(key, value) {
          if (key === 'Access-Control-Allow-Origin') {
            expect(value).toEqual('https://localhost:8080');
          }
          if (key === 'Access-Control-Allow-Methods') {
            expect(value).toEqual('GET, POST, OPTIONS, PUT, PATCH, DELETE');
          }
          if(key === 'Access-Control-Allow-Headers') {
            expect(value).toEqual('X-Requested-With,content-type');
          }
        },
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });

    it("successfully for https://localhost:8080 (over https)", function () {
      /**
       * Swagger UI Passes CORS
       */
      var req = {
        get: function(str){
          return "";
        },
        params: {
          username: "exampleAPI"
        },
        headers: {
          origin: "https://localhost:8080"
        },
        method: 'POST'
      };
      var resp = {
        setHeader: function(key, value) {
          if (key === 'Access-Control-Allow-Origin') {
            expect(value).toEqual('https://localhost:8080');
          }
          if (key === 'Access-Control-Allow-Methods') {
            expect(value).toEqual('GET, POST, OPTIONS, PUT, PATCH, DELETE');
          }
          if(key === 'Access-Control-Allow-Headers') {
            expect(value).toEqual('X-Requested-With,content-type');
          }
        },
        end: function(){
          expect(false).toBeTruthy();
        }
      };
      authHelper.middleware(req,resp,function(){
        expect(true).toBeTruthy();
      });
    });
  });
});
