var request = require("request");
var config = require('config');

const EMAIL_SERVER = config.email.serverAddress;

var sendEmail = function(subject, htmlBody, recipientList, fromName) {
  var email = {};
  email.subject = subject;
  email.html = htmlBody;
  email.recipients = recipientList;
  email.relativeHost = " "

  var from = {};
  from.name = fromName;
  from.emailAddress = "doesntGetUsed@gmail.com";
  email.from = from;

  console.log("Posting mail request with body: " + JSON.stringify(email));
  console.log("Email server: " + EMAIL_SERVER);
  request.post({headers: {'content-type' : 'application/json'}, url:EMAIL_SERVER, body: JSON.stringify(email)},
    function(error, response, body) {
      if(!error && response.statusCode == 200) {
        console.log("Email sent successfully");
        console.log("Body: " + body.toString());
      } else {
        console.error("Response code: " + response.statusCode)
        console.error("There was an error when sending email: " + error);
      }
  });
};

module.exports = {
  sendEmail: sendEmail
};
