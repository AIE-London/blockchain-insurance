var config =  require('config');
var cloudant;
var db = {};
var dbCredentials = {};

function initDBConnection() {
    dbCredentials.name = config.database.cloudant.credentials.name;

    dbCredentials.host = config.database.cloudant.credentials.host;
    dbCredentials.port = config.database.cloudant.credentials.port;
    dbCredentials.user = config.database.cloudant.credentials.user;
    dbCredentials.password = config.database.cloudant.credentials.password;
    dbCredentials.url = config.database.cloudant.credentials.url;

    cloudant = require('cloudant')(dbCredentials.url);

    // check if DB exists if not create
    cloudant.db.create(dbCredentials.name, function (err, res) {
        if (err) {
            console.log(config.database.cloudant.messages.dbExists);
            //console.log(err);
        }else{
            console.log(config.database.cloudant.messages.dbCreated);
        }
    });

    db = cloudant.use(dbCredentials.name);
}

initDBConnection();

module.exports = {
    getDb: function(){
        return db;
    },
    getDoc: function(id, callback){
        db.get(id, {revs_info: false}, function(err, doc){
            callback(err, doc);
        });
    },
    saveNewDoc: function(data, callback){
        delete data._rev;
        db.insert(data, function(err, dbresult){
                if (!err && dbresult.ok)
                {
                    console.log("Successfully created doc: " + dbresult._id);
                }else{
                    console.log("Failed to create doc");
                }
                callback(err, dbresult);
            }
        );
    },
    saveExistingDoc: function(data, callback){
        var doc = db.get(data._id, {revs_info: false}, function(err, doc)
        {

            db.insert(data, function(err, dbresult){
                if (!err && dbresult.ok)
                {
                    data._rev = doc._rev;
                    console.log("Successfully saved doc: " + data._id);
                }else{
                    console.log("Failed to save doc: " + data._id);
                }
                callback(err, dbresult);
            });
        });
    },

    doesDocExist: function(id, callback){
        var exists = false;

        db.get(id, {revs_info: false}, function(err, doc)
        {
            if (err){
                if (err.statusCode == '404'){
                    exists = false;
                }
            }else{
                exists = true;
            };
            callback(exists);
        });
    },
    deleteDoc: function(id, callback){
        db.get(id, {revs_info: false}, function(err, doc) {
            var responseBody = {};
            if (!err) {
                db.destroy(doc._id, doc._rev, function(er, body, header) {
                    callback(er);
                });
            }
            else {
                callback(err);
            }
        });
    }

}
