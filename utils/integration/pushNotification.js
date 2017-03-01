/**
 * Created by dcotton on 20/02/2017.
 */

var req = require("then-request");

const _gcmToken = require('../../config/credentials/gcm').gcm;

const _pushEndpoint = "https://fcm.googleapis.com/fcm/send";

var sendPush = function (id, title, body) {
  return req('POST', _pushEndpoint, {
    headers: {
      Authorization: "key=" + _gcmToken
    },
    json: {
      to: id,
      notification: {
        title: title,
        body: body
      }
    }
  });
};

module.exports = {
  send: sendPush
};
