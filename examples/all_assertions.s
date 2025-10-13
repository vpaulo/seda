# Comprehensive Guide to All Assertion Types
# This file demonstrates every assertion operator available in the testing framework

# ====================
# 1. is - Equality Check
# ====================
check "is assertion - equality" ::
  var x = 10
  var y = 10
  var name = "Alice"
  var arr1 = [1, 2, 3]
  var arr2 = [1, 2, 3]

  x is 10
  y is x
  name is "Alice"
  arr1 is arr2
end

# ====================
# 2. isNot - Inequality Check
# ====================
check "isNot assertion - inequality" ::
  var x = 10
  var y = 20
  var arr1 = [1, 2]
  var arr2 = [3, 4]

  x isNot y
  x isNot 5
  "hello" isNot "world"
  arr1 isNot arr2
end

# ====================
# 3. isA - Type Check
# ====================
check "isA assertion - type checking" ::
  var num = 42
  var str = "hello"
  var arr = [1, 2, 3]
  var bool_val = true

  num isA "NUMBER"
  str isA "STRING"
  arr isA "ARRAY"
  bool_val isA "BOOLEAN"
end

# ====================
# 4. contains - Membership Check
# ====================
check "contains assertion - membership" ::
  var arr = [1, 2, 3, 4, 5]
  var text = "hello world"

  arr contains 3
  arr contains 1
  text contains "world"
  text contains "hello"
end

# ====================
# 5. isGreater - Numeric Greater Than
# ====================
check "isGreater assertion - greater than comparison" ::
  var x = 10
  var y = 5

  x isGreater y
  x isGreater 9
  100 isGreater 99
  10.5 isGreater 10.4
end

# ====================
# 6. isLess - Numeric Less Than
# ====================
check "isLess assertion - less than comparison" ::
  var x = 5
  var y = 10

  x isLess y
  x isLess 6
  99 isLess 100
  10.4 isLess 10.5
end

# ====================
# 7. isTrue - Boolean True Check
# ====================
check "isTrue assertion - true check" ::
  var flag = true
  var result = (5 > 3)

  flag isTrue
  result isTrue
  (10 == 10) isTrue
  (true && true) isTrue
end

# ====================
# 8. isFalse - Boolean False Check
# ====================
check "isFalse assertion - false check" ::
  var flag = false
  var result = (5 < 3)

  flag isFalse
  result isFalse
  (10 == 11) isFalse
  (true && false) isFalse
end

# ====================
# 9. isEmpty - Empty Array/String Check
# ====================
check "isEmpty assertion - empty check" ::
  var empty_arr = []
  var empty_str = ""
  var numbers = [1, 2, 3]
  var text = "hello"

  empty_arr isEmpty
  empty_str isEmpty
end

# ====================
# 10. startsWith - String Prefix Check
# ====================
check "startsWith assertion - prefix check" ::
  var greeting = "Hello, World!"
  var path = "/home/user/file.txt"

  greeting startsWith "Hello"
  greeting startsWith "Hello, "
  path startsWith "/home"
end

# ====================
# 11. endsWith - String Suffix Check
# ====================
check "endsWith assertion - suffix check" ::
  var greeting = "Hello, World!"
  var filename = "document.pdf"

  greeting endsWith "World!"
  greeting endsWith "!"
  filename endsWith ".pdf"
  filename endsWith "pdf"
end

# ====================
# 12. raises - Error/Exception Check
# ====================
check "raises assertion - error check" ::
  # Check that an error is raised (any error)
  (10 / 0) raises

  # Check that a specific error message is raised
  (10 / 0) raises "division by zero"
end

# ====================
# Combined Example - Real World Usage
# ====================
fn calculate_discount(price, discount_percent) ::
  if price < 0 ::
    return nil, error("price cannot be negative")
  end

  if discount_percent < 0 || discount_percent > 100 ::
    return nil, error("discount must be between 0 and 100")
  end

  var discount = price * (discount_percent / 100)
  return price - discount, nil
end

check "calculate_discount - comprehensive test" ::
  # Test valid inputs
  var result1, err1 = calculate_discount(100, 10)

  result1 is 90
  isNull(err1) isTrue
  result1 isA "NUMBER"
  result1 isGreater 80
  result1 isLess 100

  # Test error conditions
  var result2, err2 = calculate_discount(-50, 10)

  isNull(err2) isFalse
  err2.to_string() contains "negative"
  err2.to_string() startsWith "price"

  # Test boundary conditions
  var result3, err3 = calculate_discount(100, 0)
  result3 is 100

  var result4, err4 = calculate_discount(100, 100)
  result4 is 0
end

# ====================
# where:: Block Example
# ====================
fn add(a, b) ::
  return a + b
where ::
  var result = add(3, 5)

  result is 8
  result isA "NUMBER"
  result isGreater 7
  result isLess 10
  result isNot 10
end

add(3, 5)

# ====================
# Array and String Assertions
# ====================
check "Array and String operations" ::
  var numbers = [1, 2, 3, 4, 5]
  var text = "The quick brown fox"
  var empty_array = []
  var empty_string = ""

  # Array assertions
  numbers contains 3
  numbers isA "ARRAY"
  empty_array isEmpty

  # String assertions
  text contains "quick"
  text startsWith "The"
  text endsWith "fox"
  text isA "STRING"
  empty_string isEmpty
end

println("✓ All assertion type examples completed!")
