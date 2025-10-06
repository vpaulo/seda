# Type System Test Suite

println("Running type system tests...")

# Type Aliases
type UserId = number
type Username = string
type UserAge = number

const id: UserId = 12345
const name: Username = "alice"
const age: UserAge = 25

check "type aliases" ::
  id isA UserId
  name isA Username
  age isA UserAge
end

# Type Checking with Primitives
check "primitive type checking" ::
  var arr = [1, 2, 3]
  42 isA "number"
  "hello" isA "string"
  true isA "boolean"
  arr isA "array"
  {"key": "value"} isA "map"
end

# Type Annotations on Variables
const x: number = 100
const s: string = "test"
const b: boolean = true

check "type annotations" ::
  x is 100
  s is "test"
  b is true
end

# Complex Type Aliases
type Point = map
type Vector = array

const point: Point = {"x": 10, "y": 20}
const vector: Vector = [1, 2, 3]

check "complex type aliases" ::
  point isA Point
  vector isA Vector
  point["x"] is 10
  vector.first is 1
end

println("✓ All type system tests passed!")
