/**
 * Only overwrites fields in master which exist in patch and leaves other master values intact
 * @param master
 * @param patch
 * @returns master
 */
var patchObjects = function(master, patch){

    var masterObj = master;

    for (var key in masterObj){
        if (patch[key]){
            if (isObject(master[key]) && isObject(patch[key])) { // If Object patch recursively
                master[key] = patchObjects(master[key], patch[key]);
            }else if (patch[key] != null){ // if not object and patch is not null, patch value
                master[key] = patch[key];
            }
        }
    }

    var patchObj = patch;

    for (var key in patchObj) {
        if (master[key] == null) { // Add any attributes in patch but not master
            master[key] = patch[key];
        }
    }

    return master;
};

/**
 *
 * @param item
 * @returns true or false (boolean)
 */
var isObject =  function(item){
    return (item !== null && typeof item === "object")
};

module.exports = {
    patchObjects: function(master, patch){
        return patchObjects(master, patch);
    },
    deleteProperties: function (object, array) {
        array.forEach(function (item) {
          delete object[item];
        });
        return object;
    },
    deReferenceSchema: function (swagger, schema, refName) {
      if (schema.definitions && schema.definitions !== true){
        Object.keys(schema.definitions).forEach(function (key) {
          swagger[key] = schema.definitions[key];
        });
        delete schema.definitions;
      }
      schema = module.exports.deleteProperties(schema, ["$schema", "title", "description"]);
      swagger[refName] = schema;
      return swagger;
    }
};
