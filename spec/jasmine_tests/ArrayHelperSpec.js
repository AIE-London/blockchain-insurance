describe("ArrayHelper", function() {
  arrayHelper = require('../../utils/helpers/array');

  beforeEach(function () {

  });

  /**
   * Tests running against doesItemExist(array, <item>>)
   */
  describe("should be able to identify", function () {
    it("if string exists within array", function () {
      /**
       * Positive String
       */
      expect(arrayHelper.doesItemExist(["test123", "test1234"], "test123")).toBeTruthy();
      /**
       * Negative String
       */
      expect(arrayHelper.doesItemExist(["test123", "test1234"], "test1235")).toBeFalsy();
      /**
       * Negative Integer
       */
      expect(arrayHelper.doesItemExist(["test123", "test1234", "1"], 1)).toBeFalsy();
      /**
       * Negative Object
       */
      expect(arrayHelper.doesItemExist(["test123", "test1234"], {"test": "test1234"})).toBeFalsy();
      /**
       * Negative Array
       */
      expect(arrayHelper.doesItemExist(["test123", "test1234"], ["test1235"])).toBeFalsy();
    });

    it("if integer exists within array", function () {
      /**
       * Positive Integer
       */
      expect(arrayHelper.doesItemExist([1, 2], 1)).toBeTruthy();
      /**
       * Negative Integer
       */
      expect(arrayHelper.doesItemExist([1, 2], 3)).toBeFalsy();
      /**
       * Negative String
       */
      expect(arrayHelper.doesItemExist([1, 2], "2")).toBeFalsy();
      /**
       * Negative Object
       */
      expect(arrayHelper.doesItemExist([1, 2], {"test": 2})).toBeFalsy();
      /**
       * Negative Array
       */
      expect(arrayHelper.doesItemExist([1, 2], [1, 2])).toBeFalsy();
      /**
       * Negative Decimal
       */
      expect(arrayHelper.doesItemExist([1, 2], 1.1)).toBeFalsy();
    });

    it("if object exists within array", function () {
      /**
       * Positive Object
       */
      expect(arrayHelper.doesItemExist([{"id": 123}, {"id": 12334}], {"id": 123})).toBeTruthy();
      /**
       * Positive Multi-Attributed Object
       */
      expect(arrayHelper.doesItemExist([{"id": 123, "animal": "dog"}, {
        "id": 33,
        "animal": "mouse"
      }], {"animal": "mouse", "id": 33})).toBeTruthy();
      /**
       * Negative Object Save Attribute Different Value
       */
      expect(arrayHelper.doesItemExist([{"id": 123, "animal": "dog"}, {
        "id": 34,
        "animal": "mouse"
      }], {"animal": "mouse", "id": 33})).toBeFalsy();
      /**
       * Negative Object Different Attribute Same Value
       */
      expect(arrayHelper.doesItemExist([{"id": 123, "animal": "dog"}, {
        "id": 33,
        "animals": "mouse"
      }], {"animal": "mouse", "id": 33})).toBeFalsy();
      /**
       * Negative String Matching JSON
       */
      expect(arrayHelper.doesItemExist([{"id": 123, "animal": "dog"}], '{"id":123,"animal":"dog"}')).toBeFalsy();
      /**
       * Negative Int
       */
      expect(arrayHelper.doesItemExist([{"id": 123, "animal": "dog"}], 123)).toBeFalsy();
      /**
       * Negative Array
       */
      expect(arrayHelper.doesItemExist([{"id": 123, "animal": "dog"}], [123, 456])).toBeFalsy();
    });
  });

  describe("should be able to add", function () {

    it("string to items", function () {
      /**
       * Positive Empty Array
       */
      var array = [];
      var expected = ["cat"];
      expect(arrayHelper.addItemToItems(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Possitive String Array
       */
      array = ["dog"];
      expected = ["dog", "cat"];
      expect(arrayHelper.addItemToItems(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Possitive Int Array
       */
      array = [1, 2];
      expected = [1, 2, "cat"];
      expect(arrayHelper.addItemToItems(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Obj Array
       */
      array = [{name: "dog"}, {name: "mouse"}];
      expected = [{name: "dog"}, {name: "mouse"}, "cat"];
      expect(arrayHelper.addItemToItems(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Mixed Array
       */
      array = [{name: "dog"}, 1];
      expected = [{name: "dog"}, 1, "cat"];
      expect(arrayHelper.addItemToItems(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
    });

    it("number to items", function () {
      /**
       * Positive Empty Array
       */
      var array = [];
      var expected = [1];
      expect(arrayHelper.addItemToItems(array, 1)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Int Array
       */
      array = [1];
      expected = [1, 2];
      expect(arrayHelper.addItemToItems(array, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive String Array
       */
      array = ["cat", "dog"];
      expected = ["cat", "dog", 1];
      expect(arrayHelper.addItemToItems(array, 1)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Obj Array
       */
      array = [{name: "dog"}, {name: "mouse"}];
      expected = [{name: "dog"}, {name: "mouse"}, 1];
      expect(arrayHelper.addItemToItems(array, 1)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Mixed Array
       */
      array = [{name: "dog"}, "1"];
      expected = [{name: "dog"}, "1", 2];
      expect(arrayHelper.addItemToItems(array, 2)).toBeTruthy();
      expect(array).toEqual(expected);
    });

    it("object to items", function () {
      /**
       * Positive Empty Array
       */
      var array = [];
      var expected = [{"animal": "cat"}];
      expect(arrayHelper.addItemToItems(array, {"animal": "cat"})).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Object Array
       */
      array = [{"animal": "dog"}];
      expected = [{"animal": "dog"}, {"animal": "mouse"}];
      expect(arrayHelper.addItemToItems(array, {"animal": "mouse"})).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive String Array
       */
      array = ["dog"];
      expected = ["dog", {"animal": "cat"}];
      expect(arrayHelper.addItemToItems(array, {"animal": "cat"})).toBeTruthy();
      expect(array).toEqual(expected);

      /**
       * Positive Int Array
       */
      array = [1];
      expected = [1, {"animal": "cat"}];
      expect(arrayHelper.addItemToItems(array, {"animal": "cat"})).toBeTruthy();
      expect(array).toEqual(expected);
    });
  });

  describe("should be able to add (if not already existing)", function () {

    it("string to items", function () {
      /**
       * Positive Empty Array
       */
      var array = [];
      var expected = ["cat"];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Populated Not Exists
       */
      array = ["dog"];
      expected = ["dog", "cat"];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Already Exists So Don't Add
       */
      array = ["dog"];
      expected = ["dog"];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, "dog")).toBeTruthy();
      expect(array).toEqual(expected);

    });

    it("number to items", function () {
      /**
       * Positive Empty Array
       */
      var array = [];
      var expected = [12];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, 12)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Populated Not Exists
       */
      array = [11];
      expected = [11, 12];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, 12)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive Already Exists So Don't Add
       */
      array = [];
      expected = [];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, 11)).toBeTruthy();
      expect(array).toEqual([11]);
    });

    it("object to items", function () {
      /**
       * Positive Empty Array
       */
      var array = [];
      var expected = [{"animal": "cat"}];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, {"animal": "cat"})).toBeTruthy();
      expect(array).toEqual(expected);

      /**
       * Positive Populated Not Exists
       */
      array = [{"animal": "dog"}];
      expected = [{"animal": "dog"}, {"animal": "mouse"}];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, {"animal": "mouse"})).toBeTruthy();
      expect(array).toEqual(expected);

      /**
       * Positive Already Exists So Don't Add
       */
      array = [{"animal": "dog"}];
      expected = [{"animal": "dog"}];
      expect(arrayHelper.addItemToItemsIfNotAlreadyExists(array, {"animal": "dog"})).toBeTruthy();
      expect(array).toEqual(expected);

    });
  });

  describe("should be able to remove all occurrences of a", function () {

    it("string from items", function () {
      /**
       * Positive remove single occurrences where only that value exists
       */
      var array = ["cat"];
      var expected = [];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive remove single occurrences where other items exist
       */
      array = ["dog", "cat"];
      expected = ["dog"];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive remove multiple occurrences
       */
      array = ["mouse", "cat", "cat", "dog"];
      expected = ["mouse", "dog"];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, "cat")).toBeTruthy();
      expect(array).toEqual(expected);
    });

    it("number from items", function () {
      /**
       * Positive remove single occurrences where only that value exists
       */
      array = [1];
      expected = [];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, 1)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive remove single occurrences where other items exist
       */
      array = [1, 2];
      expected = [1];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive remove multiple occurrences
       */
      array = [1, 2, 2, 3];
      expected = [1, 3];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, 2)).toBeTruthy();
      expect(array).toEqual(expected);
    });

    it("object from items", function () {
      /**
       * Positive remove single occurrences where only that value exists
       */
      var array = [{"animal": "cat"}];
      var expected = [];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, {"animal": "cat"})).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive remove single occurrences where other items exist
       */
      array = [{"animal": "cat"}, {"animal": "dog"}];
      expected = [{"animal": "dog"}];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, {"animal": "cat"})).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive remove multiple occurrences
       */
      array = [{"animal": "mouse"}, {"animal": "cat"}, {"animal": "cat"}, {"animal": "dog"}];
      expected = [{"animal": "mouse"}, {"animal": "dog"}];
      expect(arrayHelper.removeAllOccurrencesOfItem(array, {"animal": "cat"})).toBeTruthy();
      expect(array).toEqual(expected);
    });
  });

  describe("should be able to remove x occurrences of a", function () {

    it("string from items", function () {
      /**
       * Positive 1 where only that value exists
       */
      var array = ["cat"];
      var expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, "cat", 1)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 1 occurrence of value exists
       */
      array = ["cat"];
      expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, "cat", 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 2 occurrence of value exists
       */
      array = ["cat", "cat"];
      expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, "cat", 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 2 occurrence of value exists between other values
       */
      array = ["dog", "cat", "cat", "mouse"];
      expected = ["dog", "mouse"];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, "cat", 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 1 where 2 occurrence of value exists between other values
       */
      array = ["dog", "cat", "cat", "mouse"];
      expected = ["dog", "cat", "mouse"];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, "cat", 1)).toBeTruthy();
      expect(array).toEqual(expected);
    });

    it("number from items", function () {
      /**
       * Positive 1 where only that value exists
       */
      var array = [2];
      var expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, 2, 1)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 1 occurrence of value exists
       */
      array = [2];
      expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, 2, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 2 occurrence of value exists
       */
      array = [2, 2];
      expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, 2, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 2 occurrence of value exists between other values
       */
      array = [1, 2, 2, 1];
      expected = [1, 1];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, 2, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 1 where 2 occurrence of value exists between other values
       */
      array = [1, 2, 2, 1];
      expected = [1, 2, 1];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, 2, 1)).toBeTruthy();
      expect(array).toEqual(expected);
    });

    it("object from items", function () {
      /**
       * Positive 1 where only that value exists
       */
      var array = [{"animal": "cat"}];
      var expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, {"animal": "cat"}, 1)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 1 occurrence of value exists
       */
      array = [{"animal": "cat"}];
      expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, {"animal": "cat"}, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 2 occurrence of value exists
       */
      array = [{"animal": "cat"}, {"animal": "cat"}];
      expected = [];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, {"animal": "cat"}, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 2 where only 2 occurrence of value exists between other values
       */
      array = [{"animal": "dog"}, {"animal": "cat"}, {"animal": "cat"}, {"animal": "dog"}];
      expected = [{"animal": "dog"}, {"animal": "dog"}];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, {"animal": "cat"}, 2)).toBeTruthy();
      expect(array).toEqual(expected);
      /**
       * Positive 1 where 2 occurrence of value exists between other values
       */
      array = [{"animal": "dog"}, {"animal": "cat"}, {"animal": "cat"}, {"animal": "dog"}];
      expected = [{"animal": "dog"}, {"animal": "cat"}, {"animal": "dog"}];
      expect(arrayHelper.removeNumOccurrencesOfItem(array, {"animal": "cat"}, 1)).toBeTruthy();
      expect(array).toEqual(expected);
    });
  });
});
