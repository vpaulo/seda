# Mock Third-Party Library
module TestLib ::
  fn greet(name) ::
    return "Hello, " + name + " from TestLib!"
  end

  fn multiply(a, b) ::
    return a * b
  end

  var VERSION = "1.0.0"
  var AUTHOR = "Test Author"
end
