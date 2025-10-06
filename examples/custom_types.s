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

println("✓ All custom type system tests passed!")
