println("Running third party package tests...")

# Test third-party package loading
using "github.com/vpaulo/seda/std" as TestLib

# Test function calls
var greeting = TestLib.greet("World")
var product = TestLib.multiply(5, 7)

# Test accessing variables
var version = TestLib.VERSION
var author = TestLib.AUTHOR

# Output results
check "Third-party module tests" ::
  greeting is "Hello, World from TestLib!"
  product is 35
  version is "1.0.0"
  author is "Test Author"
end

println("✓ All third party package tests passed!")