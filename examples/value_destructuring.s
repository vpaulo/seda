println("Running value destructuring tests...")

# Multiple Return Values Examples

fn get_user_info() ::
    return "Alice", 30, "alice@example.com"
end

check "Destructuring multiple return values" ::
  var name, age, email = get_user_info()

  name is "Alice"
  age is 30
  email is "alice@example.com"
end

# Example 3: Function that returns coordinates
fn get_coordinates() ::
    return 10, 20, 30
end

check "Function that returns coordinates" ::
  var x, y, z = get_coordinates()

  x is 10
  y is 20
  z is 30
end

# Example 4: Function that returns quotient and remainder
fn div_mod(a, b) ::
    var quotient = a / b
    var remainder = a % b
    return quotient, remainder
end

check "Function that returns quotient and remainder" ::
  var q, r = div_mod(17, 5)

  q is 3.4
  r is 2
end

println("✓ All value destructuring tests passed!")
