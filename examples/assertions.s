# Test file for new assertion operators

check "isNot assertion tests" ::
  var x = 5
  var y = 10
  x isNot y
  x isNot 10
  "hello" isNot "world"
end

check "isGreater assertion tests" ::
  var a = 10
  var b = 5
  a isGreater b
  a isGreater 5
  100 isGreater 50
end

check "isLess assertion tests" ::
  var a = 5
  var b = 10
  a isLess b
  a isLess 10
  50 isLess 100
end

check "isTrue assertion tests" ::
  var t = true
  var condition = 5 > 3
  t isTrue
  condition isTrue
  (10 == 10) isTrue
end

check "isFalse assertion tests" ::
  var f = false
  var condition = 5 > 10
  f isFalse
  condition isFalse
  (10 == 5) isFalse
end

check "isEmpty assertion tests" ::
  var empty_array = []
  var empty_string = ""
  empty_array isEmpty
  empty_string isEmpty
end

check "startsWith assertion tests" ::
  var greeting = "Hello, World!"
  var name = "Alice"
  greeting startsWith "Hello"
  name startsWith "Al"
  "testing" startsWith "test"
end

check "endsWith assertion tests" ::
  var greeting = "Hello, World!"
  var filename = "test.txt"
  greeting endsWith "World!"
  filename endsWith ".txt"
  "testing" endsWith "ing"
end

check "raises assertion tests" ::
  # Test division by zero error
  (10 / 0) raises "division by zero"

  # Test undefined identifier error
  (undefined_variable) raises "identifier not found"

  # Test another division by zero
  (100 / 0) raises

  # Test type mismatch error
  ("hello" + 5) raises "type mismatch"

  # Test invalid operator
  (true + false) raises "unknown operator"
end

println("All assertion tests completed!")
