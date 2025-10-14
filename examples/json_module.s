# JSON Module Comprehensive Tests
# Tests all JSON module functions from stdlib_plan.md

println("Testing JSON Module...")

# Test JSON.parse - parse simple object
check "JSON.parse - simple object" ::
  var json_str = "{\"name\": \"Alice\", \"age\": 30}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  obj isA "map"
  obj["name"] is "Alice"
  obj["age"] is 30
end

# Test JSON.parse - simple array
check "JSON.parse - simple array" ::
  var json_str = "[1, 2, 3, 4, 5]"
  var arr, err = JSON.parse(json_str)
  isNull(err) isTrue
  arr isA "array"
  arr.length() is 5
  arr[0] is 1
  arr[4] is 5
end

# Test JSON.parse - nested object
check "JSON.parse - nested object" ::
  var json_str = "{\"person\": {\"name\": \"Bob\", \"age\": 25}}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  obj["person"] isA "map"
  obj["person"]["name"] is "Bob"
  obj["person"]["age"] is 25
end

# Test JSON.parse - nested array
check "JSON.parse - nested array" ::
  var json_str = "{\"numbers\": [1, 2, 3]}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  obj["numbers"] isA "array"
  obj["numbers"].length() is 3
  obj["numbers"][0] is 1
end

# Test JSON.parse - boolean values
check "JSON.parse - boolean values" ::
  var json_str = "{\"active\": true, \"deleted\": false}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  obj["active"] isTrue
  obj["deleted"] isFalse
end

# Test JSON.parse - null value
check "JSON.parse - null value" ::
  var json_str = "{\"value\": null}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  isNull(obj["value"]) isTrue
end

# Test JSON.parse - string values
check "JSON.parse - string values" ::
  var json_str = "{\"message\": \"hello world\", \"greeting\": \"hi\"}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  obj["message"] is "hello world"
  obj["greeting"] is "hi"
end

# Test JSON.parse - number values
check "JSON.parse - number values" ::
  var json_str = "{\"int\": 42, \"float\": 3.14, \"negative\": -10}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  obj["int"] is 42
  obj["float"] is 3.14
  obj["negative"] is -10
end

# Test JSON.parse - error on invalid JSON
check "JSON.parse - invalid JSON" ::
  var obj, err = JSON.parse("{invalid json}")
  isNull(obj) isTrue
  !isNull(err) isTrue
  err.to_string() contains "invalid JSON"
end

# Test JSON.stringify - simple object
check "JSON.stringify - simple object" ::
  var obj = {"name": "Bob", "age": 25}
  var json = JSON.stringify(obj)
  json contains "\"name\""
  json contains "\"Bob\""
  json contains "\"age\""
  json contains "25"
end

# Test JSON.stringify - simple array
check "JSON.stringify - simple array" ::
  var arr = [1, 2, 3, 4, 5]
  var json = JSON.stringify(arr)
  json is "[1,2,3,4,5]"
end

# Test JSON.stringify - nested object
check "JSON.stringify - nested object" ::
  var obj = {"person": {"name": "Alice", "age": 30}}
  var json = JSON.stringify(obj)
  json contains "\"person\""
  json contains "\"name\""
  json contains "\"Alice\""
end

# Test JSON.stringify - with indent (pretty print)
check "JSON.stringify - with indent" ::
  var obj = {"name": "Bob", "age": 25}
  var pretty = JSON.stringify(obj, 2)
  pretty contains "\n"
  pretty contains "  "
  pretty contains "\"name\""
end

# Test JSON.stringify - boolean values
check "JSON.stringify - boolean values" ::
  var obj = {"active": true, "deleted": false}
  var json = JSON.stringify(obj)
  json contains "true"
  json contains "false"
end

# Test JSON.stringify - null value
check "JSON.stringify - null value" ::
  var obj = {"value": nil}
  var json = JSON.stringify(obj)
  json contains "null"
end

# Test JSON.stringify - string values
check "JSON.stringify - string values" ::
  var obj = {"message": "hello"}
  var json = JSON.stringify(obj)
  json contains "\"message\""
  json contains "\"hello\""
end

# Test JSON.stringify - number values
check "JSON.stringify - number values" ::
  var obj = {"int": 42, "float": 3.14}
  var json = JSON.stringify(obj)
  json contains "42"
  json contains "3.14"
end

# Test round-trip: parse then stringify
check "JSON round-trip - parse then stringify" ::
  var original = "{\"name\":\"Alice\",\"age\":30}"
  var obj, err = JSON.parse(original)
  isNull(err) isTrue
  var json = JSON.stringify(obj)
  json contains "\"name\""
  json contains "\"Alice\""
  json contains "\"age\""
  json contains "30"
end

# Test round-trip: stringify then parse
check "JSON round-trip - stringify then parse" ::
  var original_obj = {"x": 10, "y": 20}
  var json = JSON.stringify(original_obj)
  var obj, err = JSON.parse(json)
  isNull(err) isTrue
  obj["x"] is 10
  obj["y"] is 20
end

# Test complex nested structure
check "JSON.parse - complex nested structure" ::
  var json_str = "{\"users\": [{\"name\": \"Alice\", \"age\": 30}, {\"name\": \"Bob\", \"age\": 25}]}"
  var obj, err = JSON.parse(json_str)
  isNull(err) isTrue
  obj["users"] isA "array"
  obj["users"].length() is 2
  obj["users"][0]["name"] is "Alice"
  obj["users"][1]["name"] is "Bob"
end

# Test empty object
check "JSON - empty object" ::
  var empty_obj = {}
  var json = JSON.stringify(empty_obj)
  json is "{}"

  var parsed, err = JSON.parse("{}")
  isNull(err) isTrue
  parsed isA "map"
end

# Test empty array
check "JSON - empty array" ::
  var empty_arr = []
  var json = JSON.stringify(empty_arr)
  json is "[]"

  var parsed, err = JSON.parse("[]")
  isNull(err) isTrue
  parsed isA "array"
  parsed.length() is 0
end

println("All JSON module tests completed!")
