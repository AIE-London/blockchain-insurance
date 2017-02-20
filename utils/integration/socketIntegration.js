/**
 * Created by CSHEFIK on 14/02/2017.
 */

var socketIO;

var initialise = function (server) {
  socketIO = require('socket.io')(server);
  console.log("Hello");

  socketIO.on('connection', function (socket) {
    socket.on('join', function(room) {
      socket.join(room);
      socketIO.emit('welcome', {})
    });});
}

var sendClaimUpdate = function (claimData,roomID){
  socketIO.to(roomID).emit('claimUpdate', claimData)
}


module.exports = {
  initialise: initialise,
  sendClaimUpdate: sendClaimUpdate
};
