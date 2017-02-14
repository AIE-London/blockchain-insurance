/**
 * Created by CSHEFIK on 14/02/2017.
 */

var socketIO;

var initialise = function (server) {
  socketIO = require('socket.io')(server);

  socketIO.on('connection', function (socket) {
    socket.on('join', function(room) {
      socket.join(room);
    });});
}

var sendClaimUpdate = function (claimData,roomID){
  socketIO.to(roomID).emit('claimUpdate', claimData)
}


