

// Load the core build.
var _ = require('lodash/core');
var _array = require('lodash/array');

var deepEquals = function(o1, o2) {
    return _.isEqual(o1, o2);
};

/**
 * Checks if the given item exists in the array
 * @param items
 * @param item
 * @returns {boolean}
 */
var doesItemExist = function(items, item){
  return items.some(function(element){
    return deepEquals(element, item)
  });
};

/**
 * Adds a single occurrence of an item to an array
 * @param items
 * @param item
 * @returns {boolean}
 */
var addItemToItems = function(items, item){
    var startLength = items.length;
    items.push(item);
    return (items.length > startLength);
};

/**
 * Adds a single occurrence of an item to an array, but only if there it not already an occurrence of the item
 * @param items
 * @param item
 * @returns {boolean}
 */
var addItemToItemsIfNotAlreadyExists = function(items, item){
      if (doesItemExist(items, item)){
          return true;
      }else{
          return addItemToItems(items, item);
      }
};

/**
 * Removed every occurrence of the given item form the array
 * @param items
 * @param item
 * @returns {boolean}
 */
var removeAllOccurrencesOfItem = function(items, item){

  var result = true;

  if (doesItemExist(items, item)) {

    result = false;
    var initialLength = items.length;

    _array.remove(items, function (n) {
      return deepEquals(n, item);
    });

    result = items.length < initialLength;

  };

  return result;
};

/**
 * Removed all occurrences of the item in the array up to the maximum given. Removes from the start of the array
 * @param items
 * @param item
 * @param num
 * @returns (boolean)
 */
var removeNumOccurrencesOfItem = function(items, item, num){

  var count = 0;
  var result = true;

  if (doesItemExist(items, item)) {

    result = false;
    var initialLength = items.length;

    _array.remove(items, function (n) {
      if (deepEquals(n, item) && count < num) {
        count++;
        return true;
      } else {
        return false;
      }
    });

    result = items.length < initialLength;
  }

  return result;
};

module.exports = {
    doesItemExist: doesItemExist,
    addItemToItems:  addItemToItems,
    addItemToItemsIfNotAlreadyExists: addItemToItemsIfNotAlreadyExists,
    removeAllOccurrencesOfItem: removeAllOccurrencesOfItem,
    removeNumOccurrencesOfItem: removeNumOccurrencesOfItem
}
