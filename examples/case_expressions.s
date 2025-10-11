println("Running case expressions tests...")

# Case Expressions example - case as expressions that return values

# Basic case expression with grades
var grade = "A"
var gpa = case grade ::
  "A" => 4.0
  "B" => 3.0
  "C" => 2.0
  "D" => 1.0
  _ => 0.0
end

check "Basic case expression with grades" ::
  gpa is 4.0
end

# Case expression with numbers
var dayNum = 3
var dayName = case dayNum ::
  1 => "Monday"
  2 => "Tuesday"
  3 => "Wednesday"
  4 => "Thursday"
  5 => "Friday"
  6 => "Saturday"
  7 => "Sunday"
  _ => "Invalid day"
end

check "Case expression with numbers" ::
  dayName is "Wednesday"
end

# Case expression with complex expressions
var x = 10
var y = 5
var operation = case x + y ::
  15 => "Addition result is 15"
  20 => "Addition result is 20"
  _ => "Unexpected addition result"
end

check "Case expression with complex expressions" ::
  operation is "Addition result is 15"
end

# Case expression used in function calls
fn getDiscount(membership) ::
  return case membership ::
    "gold" => 0.20
    "silver" => 0.15
    "bronze" => 0.10
    _ => 0.0
  end
end

check "Case expression used in function calls" ::
  getDiscount("gold") is 0.20
  getDiscount("silver") is 0.15
  getDiscount("bronze") is 0.10
  getDiscount("none") is 0
end

# Case expression in variable assignment with function calls
fn getStatusCode() ::
  return 404
end

var statusMessage = case getStatusCode() ::
  200 => "OK"
  404 => "Not Found"
  500 => "Internal Server Error"
  _ => "Unknown Status"
end

check "Case expression in variable assignment with function calls" ::
  statusMessage is "Not Found"
end

# Case expression with boolean values
var isLoggedIn = true
var accessLevel = case isLoggedIn ::
  true => "full access"
  false => "guest access"
  _ => "invalid state"
end

check "Case expression with boolean values" ::
  accessLevel is "full access"
end

# Case expression with string patterns
var userRole = "admin"
var permissions = case userRole ::
  "admin" => "read, write, delete, manage"
  "editor" => "read, write"
  "viewer" => "read"
  _ => "no permissions"
end

check "Case expression with string patterns" ::
  permissions is "read, write, delete, manage"
end

# Nested case expressions (case expression inside another case expression)
var category = "electronics"
var subcat = "phone"

var department = case category ::
  "electronics" => case subcat ::
                     "phone" => "Mobile Devices"
                     "laptop" => "Computing"
                     _ => "General Electronics"
                   end
  "clothing" => "Fashion"
  _ => "General Store"
end

check "Case expressions nested" ::
  department is "Mobile Devices"
end

# Using case expressions in arithmetic
var size = "medium"
var basePrice = 10.0
var sizeMultiplier = case size ::
  "small" => 0.8
  "medium" => 1.0
  "large" => 1.2
  "xl" => 1.5
  _ => 1.0
end

var finalPrice = basePrice * sizeMultiplier

check "Using case expressions in arithmetic" ::
  finalPrice is 10
end

# Case expression in array indexing
var priorities = ["low", "medium", "high", "critical"]
var urgencyLevel = 2
var priority = case urgencyLevel ::
  0 => priorities[0]
  1 => priorities[1]
  2 => priorities[2]
  3 => priorities[3]
  _ => "unknown"
end

check "Case expression in array indexing" ::
  priority is "high"
end

# Test that case expressions return correct values
check "case expression tests" ::
  case "A" ::
    "A" => "Excellent"
    "B" => "Good"
    _ => "Unknown"
  end is "Excellent"

  case 2 ::
    1 => "One"
    2 => "Two"
    _ => "Other"
  end is "Two"

  case true ::
    true => "Yes"
    false => "No"
    _ => "Maybe"
  end is "Yes"
end

println("✓ All case expressions tests passed!")