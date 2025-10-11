println("Running calculator advanced tests...")

#|
Advanced Calculator - Mathematical operations and utilities
Demonstrates:
- Mathematical functions
- Number methods
- Expression evaluation
- Statistics calculations
- Data validation
|#

module MathUtils ::

  const PI = 3.14159265359
  const E = 2.71828182846

  # Absolute value
  fn abs(x) ::
    if x < 0 ::
      return -x
    end
    return x
  end

  # Maximum of two numbers
  fn max(a, b) ::
    if a > b ::
      return a
    end
    return b
  end

  # Minimum of two numbers
  fn min(a, b) ::
    if a < b ::
      return a
    end
    return b
  end

  # Factorial
  fn factorial(n) ::
    if n <= 1 ::
      return 1
    end
    return n * factorial(n - 1)
  end

  # Fibonacci sequence
  fn fibonacci(n) ::
    if n <= 1 ::
      return n
    end
    return fibonacci(n - 1) + fibonacci(n - 2)
  end

  # Greatest common divisor
  fn gcd(a, b) ::
    if b == 0 ::
      return a
    end
    return gcd(b, a % b)
  end

  # Least common multiple
  fn lcm(a, b) ::
    return abs(a * b) / gcd(a, b)
  end

  # Check if prime (simplified)
  fn isPrime(n) ::
    if n <= 1 ::
      return false
    end
    if n == 2 ::
      return true
    end
    if n % 2 == 0 ::
      return false
    end

    # Check odd divisors up to sqrt(n)
    var i = 3
    var maxCheck = n / 2
    for i in 3..maxCheck ::
      if n % i == 0 ::
        return false
      end
    end
    return true
  end

  # Sum of array
  fn sum(numbers) ::
    var total = 0
    for num in numbers ::
      total = total + num
    end
    return total
  end

  # Average of array
  fn average(numbers) ::
    if numbers.length == 0 ::
      return 0
    end
    return sum(numbers) / numbers.length
  end

  # Median of array (assumes sorted)
  fn median(numbers) ::
    var len = numbers.length
    if len == 0 ::
      return 0
    end

    var mid = len / 2
    if len % 2 == 0 ::
      # Even length - average of two middle numbers
      return (numbers[mid - 1] + numbers[mid]) / 2
    end
    # Odd length - middle number
    return numbers[mid]
  end

  # Variance
  fn variance(numbers) ::
    var avg = average(numbers)
    var sumSquares = 0

    for num in numbers ::
      var diff = num - avg
      sumSquares = sumSquares + (diff * diff)
    end

    return sumSquares / numbers.length
  end

  # Standard deviation
  fn stdDev(numbers) ::
    var v = variance(numbers)
    # Simple square root approximation using Newton's method
    return sqrt(v)
  end

  # Square root (Newton's method)
  fn sqrt(x) ::
    if x < 0 ::
      return 0
    end
    if x == 0 ::
      return 0
    end

    var guess = x / 2
    var epsilon = 0.0001

    for i in 0..20 ::
      var newGuess = (guess + x / guess) / 2
      var diff = abs(newGuess - guess)

      if diff < epsilon ::
        return newGuess
      end

      guess = newGuess
    end

    return guess
  end

  # Clamp value between min and max
  fn clamp(value, minVal, maxVal) ::
    if value < minVal ::
      return minVal
    end
    if value > maxVal ::
      return maxVal
    end
    return value
  end

  # Linear interpolation
  fn lerp(start, finish, t) ::
    return start + (finish - start) * t
  end

  # Map value from one range to another
  fn mapRange(value, inMin, inMax, outMin, outMax) ::
    var t = (value - inMin) / (inMax - inMin)
    return lerp(outMin, outMax, t)
  end

end

check "Basic Math Operations" ::
  var a = 15
  var b = 4
  
  a isA "number"
  b isA "number"
  a + b is 19
  a - b is 11
  a * b is 60
  a / b is 3.75
  a % b is 3
  a ^ 2 is 225
end

check "Utility Functions" ::
  MathUtils.abs(-42) is 42
  MathUtils.max(10, 25) is 25
  MathUtils.min(10, 25) is 10
end

check "Factorial and Fibonacci" ::
  MathUtils.factorial(5) is 120
  MathUtils.fibonacci(6) is 8
end

check "GCD and LCM" ::
  var num1 = 48
  var num2 = 18
  MathUtils.gcd(num1, num2) is 6
  MathUtils.lcm(num1, num2) is 144
end

check "Prime Number Checker" ::
  MathUtils.isPrime(2) is true
  MathUtils.isPrime(3) is true
  MathUtils.isPrime(4) is false
  MathUtils.isPrime(5) is true
  MathUtils.isPrime(17) is true
  MathUtils.isPrime(20) is false
  MathUtils.isPrime(23) is true
  MathUtils.isPrime(25) is false
  MathUtils.isPrime(29) is true
  MathUtils.isPrime(30) is false
end

check "Square Root Approximation" ::
  MathUtils.sqrt(1) is 1
  MathUtils.sqrt(2) is 1.414213562
  MathUtils.sqrt(3) is 1.73205081
  MathUtils.sqrt(4) is 2
  MathUtils.sqrt(5) is 2.236067978
  MathUtils.sqrt(6) is 2.449489743
  MathUtils.sqrt(7) is 2.645751311
  MathUtils.sqrt(8) is 2.828427125
  MathUtils.sqrt(9) is 3
  MathUtils.sqrt(10) is 3.16227766
end

check "Statistics" ::
  var dataset = [5, 10, 15, 20, 25, 30, 35, 40]
  
  MathUtils.sum(dataset) is 180
  MathUtils.average(dataset) is 22.5
  MathUtils.median(dataset) is 22.5
  MathUtils.variance(dataset) is 131.25
  MathUtils.stdDev(dataset) is 11.45643924
end

check "Value Clamping" ::
  var val = 150
  var minBound = 0
  var maxBound = 100
  
  MathUtils.clamp(val, minBound, maxBound) is 100
end

check "Range Mapping" ::
  # Map 0-100 to 0-255
  var oldVal = 50
  var newVal = MathUtils.mapRange(oldVal, 0, 100, 0, 255)
  
  newVal is 127.5
end

check "Constants" ::
  #MathUtils.PI is 3.141592654
  MathUtils.E is 2.718281828
  
  MathUtils.PI * 25 is 78.53981634
  2 * MathUtils.PI * 5 is 31.41592654
end

check "Expression Evaluation" ::
  var expr1 = (10 + 5) * 3 - 2
  var expr2 = 2 ^ 3 + 4 * 5
  var expr3 = (100 / 4) % 7
  
  expr1 is 43
  expr2 is 28
  expr3 is 4
end

println("✓ All calculator advanced tests passed!")