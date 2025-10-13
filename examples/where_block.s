println("Running where blocks tests...")

fn get_user_info() ::
    return "Alice", 30, "alice@example.com"
where ::
  var name, age, email = get_user_info()

  name is "Alice"
  age is 30
  email is "alice@example.com"
end

check "Destructuring multiple return values" ::
  var name, age, email = get_user_info()

  name is "Alice"
  age is 30
  email is "alice@example.com"
end

fn factorial(n) ::
  if n <= 1 :: 1 else :: n * factorial(n - 1) end
where ::
  # Use for loop to collect test data
  var test_cases = [[0, 1], [1, 1], [2, 2], [3, 6], [4, 24]]
  var all_passed = true

  for input, expected in test_cases ::
    var result = factorial(input)
    if result != expected ::
      all_passed = false
    end
  end

  # Assertion at the where block level
  all_passed isTrue
end

fn add(a, b) ::
  return a + b
where ::
  result is 8
  arg0 is 3
  arg1 is 5
end

add(3,5)

println("✓ All where blocks tests passed!")
