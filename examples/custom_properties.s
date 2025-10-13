# Custom Properties Test Suite

println("Running custom properties tests...")

# String Custom Properties
const str = "hello"
str.to_upper = fn(self) ::
  return self.upper()
end
str.label = "greeting"

check "string custom properties" ::
  str.to_upper is "HELLO"
  str.label is "greeting"
end

# Number Custom Properties
const num = 42
num.double = fn(self) ::
  return self + self
end
num.info = "answer"

check "number custom properties" ::
  num.double is 84
  num.info is "answer"
end

# Boolean Custom Properties
const flag = true
flag.invert = fn(self) ::
  if self :: return false else :: return true end
end
flag.description = "active"

check "boolean custom properties" ::
  flag.invert is false
  flag.description is "active"
end

# Array Custom Properties
const arr = [1, 2, 3, 4, 5]
arr.sum = fn(self) ::
  var total = 0
  for num in self ::
    total = total + num
  end
  return total
end

check "array custom properties" ::
  arr.sum is 15
end

# Map Custom Properties
const obj = {"a": 1, "b": 2, "c": 3}
obj.sum = fn(self) ::
  var total = 0
  for key in self ::
    total = total + self[key]
  end
  return total
end

check "map custom properties" ::
  obj.sum is 6
end

# Factory Function Pattern
fn createPerson(name, age) ::
  const person = {"name": name, "age": age}

  person.greet = fn(self) ::
    return "Hi, I'm " + self["name"]
  end

  person.birthday = fn(self) ::
    self["age"] = self["age"] + 1
  end

  return person
end

const alice = createPerson("Alice", 25)

check "factory function pattern" ::
  alice.greet is "Hi, I'm Alice"
  alice["age"] is 25
end

println("✓ All custom properties tests passed!")
