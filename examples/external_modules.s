println("Running external modules tests...")

# Load external modules
using "math.s"              # Same directory
using "utils/string.s"      # Subdirectory

check "Math Module Operations" ::
  Math.pow(3, 4) is 81
  Math.add(15, 25) is 40
  Math.multiply(6, 7) is 42
  Math.PI is 3.14159265359
  Math.E is 2.71828182846
end

check "String Module Operations" ::
  StringUtils.concat("Hello", " External Modules") is "Hello External Modules"
  StringUtils.repeat("Code", 2) is "CodeCode"
  StringUtils.addPrefix("module") is "external_module"
end

check "Combined Module Usage" ::
  var radius = 5
  var area = Math.multiply(Math.PI, Math.pow(radius, 2))
  var description = StringUtils.concat("Circle area: ", StringUtils.addPrefix("calculated"))

  area is 78.5398163397448
end

println("✓ All external modules tests passed!")