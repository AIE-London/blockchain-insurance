var objectHelper =  require('../../utils/helpers/object');

// Load the core build.
var _ = require('lodash/core');

var deepEquals = function(o1, o2) {
  return _.isEqual(o1, o2);
};
describe("ObjectHelper", function() {


  beforeEach(function() {

  });

  /**
   * Tests running against doesItemExist(array, <item>>)
   */
  describe("should be able to delete", function () {
    it("a single string property from an object with a single field", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({"a":"123"}, ["a"]), {})).toBeTruthy();
    });
    it("all string properties from an object", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({"a":"123", "b":"456"}, ["a", "b"]), {})).toBeTruthy();
    });
    it("one property from an object with multiple", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({"a":"123", "b":"456"}, ["a"]), {"b":"456"})).toBeTruthy();
    });
    it("one property from an object with no properties", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({}, ["a"]), {})).toBeTruthy();
    });
    it("all properties from an object with different types of properties", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({"a":"123", "b":2, "c":10.2, "d": true}, ["a","b","c","d"]), {})).toBeTruthy();
    });
    it("number property from an object with different types of properties", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({"a":"123", "b":2, "c":10.2, "d": true}, ["b"]), {"a":"123", "c":10.2, "d": true})).toBeTruthy();
    });
    it("boolean property from an object with different types of properties", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({"a":"123", "b":2, "c":10.2, "d": true}, ["d"]), {"a":"123", "b":2, "c":10.2})).toBeTruthy();
    });
    it("string and decimal properties from an object with different types of properties", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deleteProperties({"a":"123", "b":2, "c":10.2, "d": true}, ["a", "c"]), {"b":2, "d": true})).toBeTruthy();
    });
  });
  /**
   * Tests running against doesItemExist(array, <item>>)
   */
  describe("should be able to de-reference", function () {
    it("the definitions with multiple properties", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deReferenceSchema({}, {"name": {}, "definitions": {"a": "123", "b": 5, "c": 4.3, "d": true}}, "testing"), {"a": "123", "b": 5, "c": 4.3, "d": true, "testing": { "name": {} }})).toBeTruthy();
    });
    it("the definitions with one property string", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deReferenceSchema({}, {"name":{}, "definitions": {"a": "123"}}, "testing"), {"a": "123", "testing":{ "name":{} }})).toBeTruthy();
    });
    it("the definitions with one property string and a string property in the schema", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deReferenceSchema({}, {"name":{"name1": "Adam"}, "definitions": {"a": "123"}}, "testing"), {"a": "123", "testing":{ "name":{"name1": "Adam"} }})).toBeTruthy();
    });
    it("the definitions with one property number", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deReferenceSchema({}, {"name":{}, "definitions": {"b": 5}}, "testing"), {"b": 5, "testing":{ "name":{} }})).toBeTruthy();
    });
    it("the definitions with no properties", function() {
      /**
       * Positive String
       */
      expect(deepEquals(objectHelper.deReferenceSchema({}, {"name":{}, "definitions": {}}, "testing"), {"testing":{ "name":{} }})).toBeTruthy();
    });
  });
});
