println("Running custom type system tests...")

# ==========================================
# 1. STRING-BASED TYPES
# ==========================================

fn createPerson(name, age) ::
  const person = name

  person.age = age

  person.greet = fn(self) ::
    return "Hello, I'm " + self + " and I'm " + self.age.to_string + " years old"
  end

  person.birthday = fn(self) ::
    self.age = self.age + 1
    return "Happy birthday! Now " + self.age.to_string
  end

  return person
end

#const alice = createPerson("Alice", 25)
#println(alice.greet)
#println(alice.birthday)
#println(alice.greet)

check "String based Person type" ::
  const alice = createPerson("Alice", 25)
  
  alice.greet is "Hello, I'm Alice and I'm 25 years old"
  alice.birthday is "Happy birthday! Now 26"
  alice.age is 26
end

# ==========================================
# 2. NUMBER-BASED TYPES
# ==========================================
# Number are immutable, you can only re-assign a var 
fn createCounter(initial) ::
  const counter = initial

  counter.increment = fn(self) ::
    return self + 1
  end

  counter.decrement = fn(self) ::
    return self - 1
  end

  counter.double = fn(self) ::
    return self * 2
  end

  counter.isPositive = fn(self) ::
    return self > 0
  end

  return counter
end

check "Number based Counter type" ::
  const count = createCounter(5)
  
  count is 5
  count.increment is 6
  count.double is 10
  count.isPositive is true
  count.decrement is 4
end

fn createCounter2(initial) ::
  const counter = {"value": initial}

  counter.increment = fn(self) ::
    self["value"] = self["value"] + 1
    return self
  end

  counter.decrement = fn(self) ::
    self["value"] = self["value"] - 1
    return self
  end

  counter.double = fn(self) ::
    self["value"] = self["value"] * 2
    return self
  end

  counter.isPositive = fn(self) ::
    return self["value"] > 0
  end

  return counter
end

check "Number-Map based Counter type" ::
  const count = createCounter2(5)
  
  count.value is 5
  count.increment.value is 6
  count.double.value is 12
  count.isPositive is true
  count.decrement.value is 11
end

# ==========================================
# 3. ARRAY-BASED TYPES
# ==========================================

fn createStack() ::
  const stack = []

  stack.push = fn(self, value) ::
    return self.push(value)
  end

  stack.pop = fn(self) ::
    const value = self.last
    # Note: In a real implementation, we'd remove the last element
    return value
  end

  stack.peek = fn(self) ::
    return self.last
  end

  stack.is_empty = fn(self) ::
    return self.length == 0
  end

  stack.size = fn(self) ::
    return self.length
  end

  return stack
end

var myStack = createStack()
myStack = myStack.push(10)
myStack = myStack.push(20)
myStack = myStack.push(30)

check "String based Person type" ::
  myStack.peek is 30
  myStack.size is 3
  myStack.is_empty is false
end

# ==========================================
# 4. MAP-BASED TYPES (Classic OOP)
# ==========================================

fn createRectangle(width, height) ::
  const rect = {"width": width, "height": height}

  rect.area = fn(self) ::
    return self["width"] * self["height"]
  end

  rect.perimeter = fn(self) ::
    return 2 * (self["width"] + self["height"])
  end

  rect.isSquare = fn(self) ::
    return self["width"] == self["height"]
  end

  rect.scale = fn(self, factor) ::
    self["width"] = self["width"] * factor
    self["height"] = self["height"] * factor
    return self
  end

  rect.describe = fn(self) ::
    return "Rectangle: " + self["width"].to_string + "x" + self["height"].to_string
  end

  return rect
end

const rect = createRectangle(10, 5)
println(rect.describe)
println("Area: ", rect.area)
println("Perimeter: ", rect.perimeter)
println("Is square: ", rect.isSquare)

check "String based Person type" ::
  rect.describe is "Rectangle: 10x5"
  rect.area is 50
  rect.perimeter is 30
  rect.isSquare is false
end

# ==========================================
# 5. BOOLEAN-BASED TYPES
# ==========================================

print("\n--- Boolean-Based Flag Type ---")

fn createFlag(initialState) ::
  const flag = initialState

  flag.toggle = fn(self) ::
    if self ::
      return false
    else ::
      return true
    end
  end

  flag.describe = fn(self) ::
    if self ::
      return "Flag is ON"
    else ::
      return "Flag is OFF"
    end
  end

  flag.andWith = fn(self, other) ::
    return self and other
  end

  flag.orWith = fn(self, other) ::
    return self or other
  end

  return flag
end

const myFlag = createFlag(true)
const toggled = myFlag.toggle
const toggledFlag = createFlag(toggled)

check "String based Person type" ::
  myFlag.describe is "Flag is ON"
  toggled is false
  toggledFlag is false
end

println("✓ All custom type system tests passed!")