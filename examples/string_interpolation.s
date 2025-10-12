println("Running string interpolation tests...")

# String Interpolation Examples
# This example demonstrates string interpolation with #{expression} syntax

check "Basic variable interpolation" ::
  var name = "Joe"
  var greeting = "Hello, #{name}!"
  
  greeting is "Hello, Joe!"
end

check "Multiple variables" ::
  var firstName = "Alice"
  var lastName = "Smith"
  var fullName = "#{firstName} #{lastName}"
  
  fullName is "Alice Smith"
end

check "Number interpolation" ::
  var age = 30
  var message = "I am #{age} years old"
  
  message is "I am 30 years old"
end

check "Arithmetic expressions" ::
  var x = 5
  var y = 10
  
  "#{x} + #{y} = #{x + y}" is "5 + 10 = 15"
  "#{x} * #{y} = #{x * y}" is "5 * 10 = 50"
end

check "Boolean expressions" ::
  var score = 85
  
  "Passing grade? #{score >= 60}" is "Passing grade? true"
end

check "String concatenation with interpolation" ::
  var city = "New York"
  var country = "USA"
  var location = "Location: #{city}, #{country}"
  
  location is "Location: New York, USA"
end

# Function calls in interpolation
fn double(n) ::
  return n * 2
end

fn square(n) ::
  return n * n
end

check "Function calls in interpolation" ::
  var num = 7
  
  "Double of #{num} is #{double(num)}" is "Double of 7 is 14"
  "Square of #{num} is #{square(num)}" is "Square of 7 is 49"
end

check "Array element access" ::
  var colors = ["red", "green", "blue"]
  
  "First color: #{colors[0]}" is "First color: red"
  "Last color: #{colors[2]}" is "Last color: blue"
end

check "String method calls" ::
  var text = "hello"
  
  "Text: '#{text}' has length #{text.length}" is "Text: 'hello' has length 5"
end

check "Complex expressions" ::
  var price = 100
  var discount = 0.20
  var finalPrice = price * (1 - discount)
  
  "Original price: $#{price}" is "Original price: $100"
  "Discount: #{discount * 100}%" is "Discount: 20%"
  "Final price: $#{finalPrice}" is "Final price: $80"
end

# Building formatted output
fn formatPrice(amount) ::
  return "$#{amount}"
end

check "Building formatted output" ::
  var apple_name = "Apple"
  var apple_price = 1.50
  var banana_name = "Banana"
  var banana_price = 0.75
  
  "Item: #{apple_name}, Price: #{formatPrice(apple_price)}" is "Item: Apple, Price: $1.5"
  "Item: #{banana_name}, Price: #{formatPrice(banana_price)}" is "Item: Banana, Price: $0.75"
end

check "Using interpolation with escape sequences" ::
  var firstName = "Alice"
  var lastName = "Smith"
  var age = 30
  var city = "New York"
  
  "Name:\t#{firstName} #{lastName}" is "Name:\tAlice Smith"
  "Age:\t#{age}" is "Age:\t30"
  "City:\t#{city}" is "City:\tNew York"
end

check "Nested calculations" ::
  var a = 2
  var b = 3
  var c = 4
  
  "Expression: (#{a} + #{b}) * #{c} = #{(a + b) * c}" is "Expression: (2 + 3) * 4 = 20"
end

# Function returning interpolated strings
fn createGreeting(person, time) ::
  return "Good #{time}, #{person}!"
end

check "Function returning interpolated strings" ::
  var person1 = "Bob"
  var time1 = "morning"
  var person2 = "Carol"
  var time2 = "evening"
  
  "#{createGreeting(person1, time1)}" is "Good morning, Bob!"
  "#{createGreeting(person2, time2)}" is "Good evening, Carol!"
end

println("✓ All string interpolation tests passed!")