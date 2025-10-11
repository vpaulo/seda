println("Running calculator tests...")

# Calculator example - a simple calculator with functions

fn add(a, b) ::
  a + b
end

fn subtract(a, b) ::
  a - b
end

fn multiply(a, b) ::
  a * b
end

fn divide(a, b) ::
  if b == 0 ::
    return "Error: Division by zero"
  else ::
    return a / b
  end
end

fn power(base, exp) ::
  base ^ exp
end

check "calculator tests" ::
  add(10, 5) is 15
  subtract(10, 5) is 5
  multiply(10, 5) is 50
  divide(10, 5) is 2
  power(2, 4) is 16
end

check "power operator tests" ::
  power(3, 2) is 9
  power(2, 3) is 8
  power(5, 0) is 1
  power(2, 8) is 256
  power(4, 0.5) is 2

  3^2 is 9
  2^3 is 8
  5^0 is 1
  2^8 is 256
  4^0.5 is 2
end

println("✓ All calculator tests passed!")