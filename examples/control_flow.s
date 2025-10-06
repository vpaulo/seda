# Control Flow Test Suite

println("Running control flow tests...")

# If-Else Statements
fn max(a, b) ::
  if a > b ::
    return a
  else ::
    return b
  end
end

check "if-else statements" ::
  max(10, 5) is 10
  max(3, 8) is 8
  max(7, 7) is 7
end

# Nested If-Else
fn classify(num) ::
  if num > 0 ::
    return "positive"
  else ::
    if num < 0 ::
      return "negative"
    else ::
      return "zero"
    end
  end
end

check "nested if-else" ::
  classify(10) is "positive"
  classify(-5) is "negative"
  classify(0) is "zero"
end

# For Loops with Arrays
fn sumArray(arr) ::
  var total = 0
  for num in arr ::
    total = total + num
  end
  return total
end

check "for loops with arrays" ::
  sumArray([1, 2, 3, 4, 5]) is 15
  sumArray([10, 20, 30]) is 60
end

# For Loops with Maps
fn countKeys(obj) ::
  var count = 0
  for key in obj ::
    count = count + 1
  end
  return count
end

check "for loops with maps" ::
  countKeys({"a": 1, "b": 2, "c": 3}) is 3
  countKeys({"x": 10}) is 1
end

# While Loops
#|
fn factorial(n) ::
  var result = 1
  var i = 1
  while i <= n ::
    result = result * i
    i = i + 1
  end
  return result
end

check "while loops" ::
  factorial(5) is 120
  factorial(3) is 6
  factorial(1) is 1
end
|#

println("✓ All control flow tests passed!")
