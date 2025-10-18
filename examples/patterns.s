println("Running patterns tests...")
# Custom Type Patterns - Advanced OOP Patterns
# Demonstrating inheritance-like behavior, composition, and factories

# ==========================================
# 1. FACTORY PATTERN - Creating similar objects
# ==========================================

fn ShapeFactory() ::
  const factory = {}

  factory.createCircle = fn(self, radius) ::
    const circle = {"type": "Circle", "radius": radius}

    circle.area = fn(self) ::
      const pi = 3.14159
      return pi * self["radius"] * self["radius"]
    end

    circle.describe = fn(self) ::
      return self["type"] + " with radius " + self["radius"].to_string
    end

    return circle
  end

  factory.createRectangle = fn(self, width, height) ::
    const rect = {"type": "Rectangle", "width": width, "height": height}

    rect.area = fn(self) ::
      return self["width"] * self["height"]
    end

    rect.describe = fn(self) ::
      return self["type"] + " " + self["width"].to_string + "x" + self["height"].to_string
    end

    return rect
  end

  return factory
end

const shapeFactory = ShapeFactory()
const circle = shapeFactory.createCircle(5)
const rect = shapeFactory.createRectangle(10, 20)

check "factory pattern" ::
  circle.describe is "Circle with radius 5"
  circle.area is 78.53975
  rect.describe is "Rectangle 10x20"
  rect.area is 200
end

# ==========================================
# 2. COMPOSITION PATTERN
# ==========================================

fn createEngine(horsepower) ::
  var engine = {"hp": horsepower, "running": false}

  engine.start = fn(self) ::
    self["running"] = true
    return "Engine started! " + self["hp"].to_string + " HP"
  end

  engine.stop = fn(self) ::
    self["running"] = false
    return "Engine stopped"
  end

  engine.isRunning = fn(self) ::
    return self["running"]
  end

  return engine
end

fn createCar(model, horsepower) ::
  var car = {
    "model": model,
    "engine": createEngine(horsepower),
    "speed": 0
  }

  car.start = fn(self) ::
    return self["engine"].start
  end

  car.accelerate = fn(self, amount) ::
    if self["engine"].isRunning ::
      self["speed"] = self["speed"] + amount
      return "Accelerating... Speed: " + self["speed"].to_string
    else ::
      return "Cannot accelerate. Engine is off."
    end
  end

  car.stop = fn(self) ::
    self["speed"] = 0
    return self["engine"].stop
  end

  car.describe = fn(self) ::
    return self["model"] + " (Engine: " + self["engine"]["hp"].to_string + " HP)"
  end

  return car
end

var myCar = createCar("Batman tumbler", 670)

check "Composition pattern" ::
  myCar.describe is "Batman tumbler (Engine: 670 HP)"
  myCar.start is "Engine started! 670 HP"
  myCar.accelerate(50) is "Accelerating... Speed: 50"
  myCar.accelerate(30) is "Accelerating... Speed: 80"
  myCar.stop is "Engine stopped"
end

# ==========================================
# 3. BUILDER PATTERN
# ==========================================

fn PersonBuilder() ::
  var builder = {
    "name": "",
    "age": 0,
    "city": "",
    "job": ""
  }

  builder.withName = fn(self, name) ::
    self["name"] = name
    return self
  end

  builder.withAge = fn(self, age) ::
    self["age"] = age
    return self
  end

  builder.withCity = fn(self, city) ::
    self["city"] = city
    return self
  end

  builder.withJob = fn(self, job) ::
    self["job"] = job
    return self
  end

  builder.build = fn(self) ::
    const person = {
      "name": self["name"],
      "age": self["age"],
      "city": self["city"],
      "job": self["job"]
    }

    person.describe = fn(self) ::
      return self["name"] + ", " + self["age"].to_string + " years old, " +
             self["job"] + " from " + self["city"]
    end

    return person
  end

  return builder
end

var person = PersonBuilder()
  .withName("John Doe")
  .withAge(30)
  .withCity("New York")
  .withJob("Developer")
  .build

check "Builder pattern" ::
  person.describe is "John Doe, 30 years old, Developer from New York"
end

# ==========================================
# 4. STATE PATTERN - Array-based state machine
# ==========================================

fn createTrafficLight() ::
  var light = {
    "currentState": 0,  # 0=red, 1=yellow, 2=green
    "states": ["RED", "YELLOW", "GREEN"]
  }

  light.getCurrentColor = fn(self) ::
    return self["states"][self["currentState"]]
  end

  light.next = fn(self) ::
    self["currentState"] = (self["currentState"] + 1) % 3
    return self
  end

  light.canGo = fn(self) ::
    return self["currentState"] == 2  # green
  end

  light.shouldSlow = fn(self) ::
    return self["currentState"] == 1  # yellow
  end

  light.mustStop = fn(self) ::
    return self["currentState"] == 0  # red
  end

  return light
end

check "State pattern" ::
  var trafficLight = createTrafficLight()
  var trafficLight2 = createTrafficLight().next
  var trafficLight3 = createTrafficLight().next.next
  
  trafficLight.getCurrentColor is "RED"
  trafficLight.canGo is false
  
  trafficLight2.getCurrentColor is "YELLOW"
  trafficLight2.shouldSlow is true
  
  trafficLight3.getCurrentColor is "GREEN"
  trafficLight3.canGo is true
end

# ==========================================
# 5. INHERITANCE-LIKE PATTERN
# ==========================================

fn createAnimal(name, species) ::
  var animal = {"name": name, "species": species}

  animal.speak = fn(self) ::
    return self["name"] + " makes a sound"
  end

  animal.describe = fn(self) ::
    return self["name"] + " is a " + self["species"]
  end

  return animal
end

fn createDog(name, breed) ::
  # Create base animal
  var dog = createAnimal(name, "Dog")

  # Add dog-specific properties
  dog["breed"] = breed

  # Override speak method
  dog.speak = fn(self) ::
    return self["name"] + " barks: Woof!"
  end

  # Add new method
  dog.fetch = fn(self) ::
    return self["name"] + " fetches the ball!"
  end

  # Extend describe method
  const originalDescribe = dog.describe
  dog.describe = fn(self) ::
    return originalDescribe + " (Breed: " + self["breed"] + ")"
  end

  return dog
end

check "Inheritance pattern" ::
  var myDog = createDog("Max", "Golden Retriever")

  myDog.describe is "Max is a Dog (Breed: Golden Retriever)"
  myDog.speak is "Max barks: Woof!"
  myDog.fetch is "Max fetches the ball!"
end

# ==========================================
# 6. SINGLETON PATTERN (Number-based)
# ==========================================

fn getConfig() ::
  const config = 1  # Dummy value, we'll use properties

  config.appName = "MyApp"
  config.version = "1.0.0"
  config.debug = true

  config.get = fn(self, key) ::
    # In a real implementation, we'd have a proper lookup
    if key == "appName" ::
      return self.appName
    else ::
      if key == "version" ::
        return self.version
      else ::
        return self.debug
      end
    end
  end

  config.describe = fn(self) ::
    return self.appName + " v" + self.version + " (Debug: " + self.debug.to_string + ")"
  end

  return config
end

check "Singleton pattern" ::
  var config = getConfig()

  config.describe is "MyApp v1.0.0 (Debug: true)"
end

# ==========================================
# 7. STRATEGY PATTERN - Different sorting strategies
# ==========================================

fn createSorter() ::
  var sorter = {"strategy": "bubble"}

  sorter.setStrategy = fn(self, strategy) ::
    self.strategy = strategy
    return self
  end

  sorter.sort = fn(self, data) ::
    if self.strategy == "bubble" ::
      return "Bubble sorting array of " + data.length.to_string + " elements"
    else ::
      if self.strategy == "quick" ::
        return "Quick sorting array of " + data.length.to_string + " elements"
      else ::
        return "Default sorting array of " + data.length.to_string + " elements"
      end
    end
  end

  return sorter
end

var sorter = createSorter()
var numbers = [3, 1, 4, 1, 5, 9]

sorter.setStrategy("quick")

var sort = sorter.sort(numbers)

check "Strategy pattern" ::
  sort is "Quick sorting array of 6 elements"
end

# ==========================================
# 8. OBSERVABLE PATTERN (Array of listeners)
# ==========================================

fn createEventEmitter() ::
  var emitter = {"listeners": []}

  emitter.on = fn(self, listener) ::
    self["listeners"] = self["listeners"].push(listener)
    return self
  end

  emitter.emit = fn(self, event) ::
    var result = "Emitting: " + event
    for listener in self["listeners"] ::
      result = result + " | Listener called"
    end
    return result
  end

  emitter.listenerCount = fn(self) ::
    return self["listeners"].length
  end

  return emitter
end

var emitter = createEventEmitter()
emitter.on("listener1")
emitter.on("listener2")

check "Observable pattern" ::
  emitter.listenerCount is 2
  emitter.emit("data-changed") is "Emitting: data-changed | Listener called | Listener called"
end

println("✓ All pattern tests passed!")