println("Testing Number methods and Math module...")

# =====================================
# Number Methods Tests
# =====================================

check "Number.abs() method" ::
  var x = -5.5
  var y = 3.14

  x.abs() is 5.5
  y.abs() is 3.14
  (-10).abs() is 10
  (0).abs() is 0
end

check "Number.floor() method" ::
  var x = 5.9
  var y = -3.2
  var z = 7.1
  var w = 10

  x.floor() is 5
  y.floor() is -4
  z.floor() is 7
  w.floor() is 10
end

check "Number.ceil() method" ::
  var x = 5.1
  var y = -3.8
  var z = 7.9
  var w = 10

  x.ceil() is 6
  y.ceil() is -3
  z.ceil() is 8
  w.ceil() is 10
end

check "Number.round() method" ::
  var x = 5.4
  var y = 5.6
  var z = -3.5
  var w = 10

  x.round() is 5
  y.round() is 6
  z.round() is -4
  w.round() is 10
end

check "Number.sqrt() method" ::
  var x = 16
  var y = 2.25
  var z = 25
  var w = 0

  x.sqrt() is 4
  y.sqrt() is 1.5
  z.sqrt() is 5
  w.sqrt() is 0
end

check "Number.sqrt() with negative number" ::
  # Should return an error when given negative number
  (-5).sqrt() raises "negative"
end

# =====================================
# Math Module Tests
# =====================================

check "Math.pow() function" ::
  Math.pow(2, 3) is 8
  Math.pow(5, 2) is 25
  Math.pow(10, 0) is 1
  Math.pow(2, -1) is 0.5
end

check "Math.max() function" ::
  Math.max(1, 5, 3) is 5
  Math.max(-10, -20, -5) is -5
  Math.max(0) is 0
  Math.max(3.14, 2.71, 1.41) is 3.14
end

check "Math.min() function" ::
  Math.min(1, 5, 3) is 1
  Math.min(-10, -20, -5) is -20
  Math.min(0) is 0
  Math.min(3.14, 2.71, 1.41) is 1.41
end

check "Math trigonometry functions" ::
  # sin(0) = 0
  Math.sin(0) is 0

  # cos(0) = 1
  Math.cos(0) is 1

  # tan(0) = 0
  Math.tan(0) is 0

  # asin(0) = 0
  Math.asin(0) is 0

  # acos(1) = 0
  Math.acos(1) is 0

  # atan(0) = 0
  Math.atan(0) is 0
end

check "Math.atan2() function" ::
  # atan2(0, 1) = 0
  Math.atan2(0, 1) is 0

  # atan2(1, 0) = π/2 ≈ 1.5708
  var result = Math.atan2(1, 0)
  result isGreater 1.57
  result isLess 1.58
end

check "Math logarithm functions" ::
  # log(1) = 0
  Math.log(1) is 0

  # log10(1) = 0
  Math.log10(1) is 0

  # log2(1) = 0
  Math.log2(1) is 0

  # exp(0) = 1
  Math.exp(0) is 1

  # log(e) = 1
  var result = Math.log(Math.E)
  result isGreater 0.99
  result isLess 1.01
end

check "Math.random() function" ::
  var r1 = Math.random()
  var r2 = Math.random()

  # Random values should be between 0 and 1
  r1 isGreater -0.1
  r1 isLess 1.1

  r2 isGreater -0.1
  r2 isLess 1.1

  # Two random calls should (very likely) be different
  r1 isNot r2
end

check "Math.random_int() function" ::
  var r1 = Math.random_int(1, 10)
  var r2 = Math.random_int(1, 10)

  # Random integers should be in range [1, 10)
  r1 isGreater 0
  r1 isLess 10

  r2 isGreater 0
  r2 isLess 10
end

check "Math constants" ::
  # PI ≈ 3.14159265359
  Math.PI isGreater 3.14
  Math.PI isLess 3.15

  # E ≈ 2.71828182846
  Math.E isGreater 2.71
  Math.E isLess 2.72

  # TAU = 2π ≈ 6.28318530718
  Math.TAU isGreater 6.28
  Math.TAU isLess 6.29
end

# =====================================
# Hybrid Approach Tests
# =====================================

check "Number methods chaining" ::
  var x = -5.7
  var result = x.abs().floor()

  result is 5
end

check "Number method + Math function combination" ::
  var base = 2.5
  var power = 3
  var result = Math.pow(base.ceil(), power)

  # ceil(2.5) = 3, pow(3, 3) = 27
  result is 27
end

check "Complex calculations" ::
  # Calculate area of circle: π * r²
  var radius = 5
  var area = Math.PI * Math.pow(radius, 2)

  area isGreater 78.5
  area isLess 78.6
end

check "Math module in expressions" ::
  var x = Math.max(1, 2, 3) + Math.min(4, 5, 6)

  # max(1,2,3) + min(4,5,6) = 3 + 4 = 7
  x is 7
end

println("✓ All Number methods and Math module tests passed!")
