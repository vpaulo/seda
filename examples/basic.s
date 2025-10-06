# Basic Language Features Test Suite

println("Running basic language tests...")

# Variables and Constants
check "variables and constants" ::
  var x = 10
  const y = 20
  x + y is 30
end

# Arithmetic Operations
check "arithmetic operations" ::
  5 + 3 is 8
  10 - 4 is 6
  6 * 7 is 42
  20 / 4 is 5
end

# String Operations
check "string operations" ::
  "hello" + " " + "world" is "hello world"
  "test".length is 4
  "UPPER".lower is "upper"
  "lower".upper is "LOWER"
end

# Boolean Operations
check "boolean operations" ::
  true and true is true
  true and false is false
  false or true is true
  #not false is true
  
  true && true is true
  true && false is false
  false || true is true
  !false is true
end

# Comparison Operations
check "comparison operations" ::
  5 < 10 is true
  10 > 5 is true
  5 == 5 is true
  5 != 10 is true
  5 <= 5 is true
  10 >= 10 is true
end

# Array Operations
check "array operations" ::
  var arr = [1, 2, 3]
  arr.length is 3
  arr.first is 1
  arr.last is 3
  arr.rest is [2, 3]
end

# Map Operations
check "map operations" ::
  var obj = {"name": "Alice", "age": 25}
  obj["name"] is "Alice"
  obj["age"] is 25
end

# Functions
fn add(a, b) ::
  return a + b
end

fn greet(name) ::
  return "Hello, " + name
end

check "functions" ::
  add(5, 3) is 8
  greet("Alice") is "Hello, Alice"
end

# Closures
fn makeCounter() ::
  var count = 0
  return fn() ::
    count = count + 1
    return count
  end
end

check "closures" ::
  const counter = makeCounter()
  counter() is 1
  counter() is 2
  counter() is 3
end

println("✓ All basic tests passed!")
