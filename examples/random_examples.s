println("Running random tests...")

# Real-World Custom Type Examples
# Practical applications showing custom types in action

# ==========================================
# 1. TODO LIST APPLICATION
# ==========================================

fn createTodo(id, title, completed) ::
  const todo = {"id": id, "title": title, "completed": completed}

  todo.toggle = fn(self) ::
    self["completed"] = not self["completed"]
    return self
  end

  todo.to_string = fn(self) ::
    var status = "[ ]"
    if self["completed"] ::
      status = "[✓]"
    end
    return status + " " + self["id"].to_string + ". " + self["title"]
  end

  return todo
end

fn createTodoList() ::
  const list = {"todos": [], "nextId": 1}

  list.add = fn(self, title) ::
    const todo = createTodo(self["nextId"], title, false)
    self["todos"] = self["todos"].push(todo)
    self["nextId"] = self["nextId"] + 1
    return self
  end

  list.complete = fn(self, id) ::
    # In a real app, we'd find and toggle the specific todo
    return "Todo " + id.to_string + " completed"
  end

  list.count = fn(self) ::
    return self["todos"].length
  end

  list.display = fn(self) ::
    var result = "Todo List (" + self["todos"].length.to_string + " items):"
    for todo in self["todos"] ::
      result = result + "\n  " + todo.to_string
    end
    return result
  end

  return list
end

var myTodos = createTodoList()
myTodos = myTodos.add("Learn the language")
myTodos = myTodos.add("Build a project")
myTodos = myTodos.add("Share with friends")

check "Todo List Application" ::
  myTodos.display is "Todo List (3 items):\n  [ ] 1. Learn the language\n  [ ] 2. Build a project\n  [ ] 3. Share with friends"
  myTodos.count is 3
end

# ==========================================
# 2. SHOPPING CART
# ==========================================

fn createProduct(name, price) ::
  const product = {"name": name, "price": price}

  product.to_string = fn(self) ::
    return self["name"] + " ($" + self["price"].to_string + ")"
  end

  return product
end

fn createCart() ::
  const cart = {"items": []}

  cart.add = fn(self, product) ::
    self["items"] = self["items"].push(product)
    return self
  end

  cart.total = fn(self) ::
    var sum = 0
    for item in self["items"] ::
      sum = sum + item["price"]
    end
    return sum
  end

  cart.count = fn(self) ::
    return self["items"].length
  end

  cart.summary = fn(self) ::
    var result = "Shopping Cart (" + self.count.to_string + " items):"
    for item in self["items"] ::
      result = result + "\n  - " + item.to_string
    end
    result = result + "\nTotal: $" + self.total.to_string
    return result
  end

  return cart
end

var cart = createCart()
cart = cart.add(createProduct("Laptop", 999))
cart = cart.add(createProduct("Mouse", 25))
cart = cart.add(createProduct("Keyboard", 75))

check "Shopping Cart" ::
  cart.summary is "Shopping Cart (3 items):\n  - Laptop ($999)\n  - Mouse ($25)\n  - Keyboard ($75)\nTotal: $1099"
end

# ==========================================
# 3. USER AUTHENTICATION
# ==========================================

fn createUser(username, password) ::
  const user = {
    "username": username,
    "password": password,
    "loggedIn": false,
    "loginAttempts": 0
  }

  user.login = fn(self, pwd) ::
    if pwd == self["password"] ::
      self["loggedIn"] = true
      self["loginAttempts"] = 0
      return "Welcome, " + self["username"] + "!"
    else ::
      self["loginAttempts"] = self["loginAttempts"] + 1
      return "Invalid password. Attempts: " + self["loginAttempts"].to_string
    end
  end

  user.logout = fn(self) ::
    self["loggedIn"] = false
    return "Goodbye, " + self["username"]
  end

  user.isAuthenticated = fn(self) ::
    return self["loggedIn"]
  end

  user.resetPassword = fn(self, oldPwd, newPwd) ::
    if oldPwd == self["password"] ::
      self["password"] = newPwd
      return "Password updated successfully"
    else ::
      return "Current password is incorrect"
    end
  end

  return user
end

check "User Authentication System" ::
  const user = createUser("alice", "secret123")

  user.login("wrong") is "Invalid password. Attempts: 1"
  user.login("secret123") is "Welcome, alice!"
  user.isAuthenticated is true
  user.logout is "Goodbye, alice"
end

# ==========================================
# 4. BANK ACCOUNT
# ==========================================

fn createBankAccount(accountNumber, initialBalance) ::
  const account = {
    "accountNumber": accountNumber,
    "balance": initialBalance,
    "transactions": []
  }

  account.deposit = fn(self, amount) ::
    self["balance"] = self["balance"] + amount
    const tx = {"type": "deposit", "amount": amount}
    self["transactions"] = self["transactions"].push(tx)
    return "Deposited $" + amount.to_string + ". New balance: $" + self["balance"].to_string
  end

  account.withdraw = fn(self, amount) ::
    if amount > self["balance"] ::
      return "Insufficient funds. Balance: $" + self["balance"].to_string
    else ::
      self["balance"] = self["balance"] - amount
      const tx = {"type": "withdrawal", "amount": amount}
      self["transactions"] = self["transactions"].push(tx)
      return "Withdrew $" + amount.to_string + ". New balance: $" + self["balance"].to_string
    end
  end

  account.getBalance = fn(self) ::
    return self["balance"]
  end

  account.statement = fn(self) ::
    var result = "Account: " + self["accountNumber"].to_string + "\n"
    result = result + "Balance: $" + self["balance"].to_string + "\n"
    result = result + "Transactions: " + self["transactions"].length.to_string
    return result
  end

  return account
end

check "Bank Account System" ::
  const myAccount = createBankAccount(12345, 1000)

  myAccount.deposit(500) is "Deposited $500. New balance: $1500"
  myAccount.withdraw(200) is "Withdrew $200. New balance: $1300"
  myAccount.withdraw(2000) is "Insufficient funds. Balance: $1300"
  myAccount.statement is "Account: 12345\nBalance: $1300\nTransactions: 2"
end

# ==========================================
# 5. TIMER/STOPWATCH
# ==========================================

fn createStopwatch() ::
  const watch = {"startTime": 0, "elapsed": 0, "running": false}

  watch.start = fn(self) ::
    const running = self["running"]
    if running ::
      return "Stopwatch already running"
    else ::
      self["running"] = true
      return "Stopwatch started"
    end
  end

  watch.stop = fn(self) ::
    const running = self["running"]
    if running ::
      self["running"] = false
      self["elapsed"] = self["elapsed"] + 1  # Simplified: increment by 1
      return "Stopwatch stopped. Elapsed: " + self["elapsed"].to_string + " units"
    else ::
      return "Stopwatch not running"
    end
  end

  watch.reset = fn(self) ::
    self["elapsed"] = 0
    self["running"] = false
    return "Stopwatch reset"
  end

  watch.getElapsed = fn(self) ::
    return self["elapsed"]
  end

  return watch
end

check "Timer/Stopwatch" ::
  var stopwatch = createStopwatch()

  stopwatch.start is "Stopwatch started"
  stopwatch.stop is "Stopwatch stopped. Elapsed: 1 units"
  stopwatch.start is "Stopwatch started"
  stopwatch.stop is "Stopwatch stopped. Elapsed: 2 units"
  stopwatch.getElapsed is 2
  stopwatch.reset is "Stopwatch reset"
end

# ==========================================
# 6. VALIDATOR
# ==========================================

fn createValidator() ::
  const validator = {"errors": []}

  validator.required = fn(self, value, fieldName) ::
    # Simple check for empty string
    if value == "" ::
      self["errors"] = self["errors"].push(fieldName + " is required")
    end
    return self
  end

  validator.minLength = fn(self, value, minLen, fieldName) ::
    if value.length < minLen ::
      const msg = fieldName + " must be at least " + minLen.to_string + " characters"
      self["errors"] = self["errors"].push(msg)
    end
    return self
  end

  validator.isValid = fn(self) ::
    return self["errors"].length == 0
  end

  validator.getErrors = fn(self) ::
    var result = "Validation errors:"
    for error in self["errors"] ::
      result = result + "\n  - " + error
    end
    return result
  end

  validator.reset = fn(self) ::
    self["errors"] = []
    return self
  end

  return validator
end

var validator = createValidator()
validator = validator.required("", "Username")
validator = validator.minLength("ab", 3, "Password")

check "Form Validator" ::
  validator.isValid is false
  validator.getErrors is "Validation errors:\n  - Username is required\n  - Password must be at least 3 characters"
end

# ==========================================
# 7. CACHE SYSTEM
# ==========================================

fn createCache(maxSize) ::
  const cache = {"data": {}, "maxSize": maxSize, "size": 0}

  cache.set = fn(self, key, value) ::
    if self["size"] < self["maxSize"] ::
      self["data"][key] = value
      self["size"] = self["size"] + 1
      return "Cached: " + key
    else ::
      return "Cache full! Max size: " + self["maxSize"].to_string
    end
  end

  cache.get = fn(self, key) ::
    return self["data"][key]
  end

  cache.has = fn(self, key) ::
    # Simple existence check
    return self["size"] > 0
  end

  cache.clear = fn(self) ::
    self["data"] = {}
    self["size"] = 0
    return "Cache cleared"
  end

  cache.stats = fn(self) ::
    return "Cache: " + self["size"].to_string + "/" + self["maxSize"].to_string + " items"
  end

  return cache
end

check "Simple Cache System" ::
  var cache = createCache(3)

  cache.set("user1", "Alice") is "Cached: user1"
  cache.set("user2", "Bob") is "Cached: user2"
  cache.stats is "Cache: 2/3 items"
  cache.get("user1") is "Alice"
end

# ==========================================
# 8. LOGGER SYSTEM
# ==========================================

fn createLogger(name) ::
  const logger = {"name": name, "logs": []}

  logger.info = fn(self, message) ::
    const log = "[INFO] " + self["name"] + ": " + message
    self["logs"] = self["logs"].push(log)
    return log
  end

  logger.error = fn(self, message) ::
    const log = "[ERROR] " + self["name"] + ": " + message
    self["logs"] = self["logs"].push(log)
    return log
  end

  logger.warn = fn(self, message) ::
    const log = "[WARN] " + self["name"] + ": " + message
    self["logs"] = self["logs"].push(log)
    return log
  end

  logger.getHistory = fn(self) ::
    var result = "Log history (" + self["logs"].length.to_string + " entries):"
    for log in self["logs"] ::
      result = result + "\n  " + log
    end
    return result
  end

  return logger
end

check "Logger System" ::
  var logger = createLogger("MyApp")

  logger.info("Application started") is "[INFO] MyApp: Application started"
  logger.warn("Low memory warning") is "[WARN] MyApp: Low memory warning"
  logger.error("Database connection failed") is "[ERROR] MyApp: Database connection failed"
  logger.getHistory is "Log history (3 entries):\n  [INFO] MyApp: Application started\n  [WARN] MyApp: Low memory warning\n  [ERROR] MyApp: Database connection failed"
end

#|
Todo List Manager - Real-world application example
Demonstrates:
- Data structures (arrays, maps)
- CRUD operations
- User-defined methods
- String interpolation
- Constants
- Module pattern
- Function composition
|#

#| Todo Module - Encapsulates all todo list functionality |#
module TodoList ::

  # Storage for todos
  var todos = []
  var nextId = 1

  # Create a new todo item
  fn create(title, description) ::
    var todo = {
      "id": nextId,
      "title": title,
      "description": description,
      "completed": false,
      "createdAt": nextId  # Using ID as timestamp for now
    }

    todos = todos.push(todo)
    nextId = nextId + 1

    return todo
  end

  # Get all todos
  fn getAll() ::
    return todos
  end

  # Get a todo by ID
  fn getById(id) ::
    for todo in todos ::
      if todo["id"] == id ::
        return todo
      end
    end
    return null
  end

  # Update a todo
  fn update(id, updates) ::
    for i in 0..todos.length ::
      if todos[i]["id"] == id ::
        # Update title if provided
        if updates["title"] ::
          todos[i]["title"] = updates["title"]
        end

        # Update description if provided
        if updates["description"] ::
          todos[i]["description"] = updates["description"]
        end

        # Update completed status if provided
        if updates["completed"] ::
          todos[i]["completed"] = updates["completed"]
        end

        return todos[i]
      end
    end
    return null
  end

  # Delete a todo
  fn delete(id) ::
    var newTodos = []
    var found = false

    for todo in todos ::
      if todo["id"] == id ::
        found = true
      end
      if todo["id"] != id ::
        newTodos = newTodos.push(todo)
      end
    end

    todos = newTodos
    return found
  end

  # Mark todo as completed
  fn complete(id) ::
    return update(id, {"completed": true})
  end

  # Mark todo as incomplete
  fn uncomplete(id) ::
    return update(id, {"completed": false})
  end

  # Get only completed todos
  fn getCompleted() ::
    var completed = []
    for todo in todos ::
      if todo["completed"] ::
        completed = completed.push(todo)
      end
    end
    return completed
  end

  # Get only pending todos
  fn getPending() ::
    var pending = []
    for todo in todos ::
      if !todo["completed"] ::
        pending = pending.push(todo)
      end
    end
    return pending
  end

  # Get count of todos
  fn count() ::
    return todos.length
  end

  # Get count of completed todos
  fn countCompleted() ::
    var count = 0
    for todo in todos ::
      if todo["completed"] ::
        count = count + 1
      end
    end
    return count
  end

  # Get count of pending todos
  fn countPending() ::
    return count() - countCompleted()
  end

  # Clear all todos
  fn clear() ::
    todos = []
    nextId = 1
  end

  # Display a single todo
  fn displayTodo(todo) ::
    var status = case todo["completed"] ::
      true => "[✓]"
      false => "[ ]"
      _ => "[?]"
    end
    
    return "#{status} ##{todo.id}: #{todo.title} - #{todo.description}\n"
  end

  # Display all todos
  fn displayAll() ::
    if todos.length == 0 ::
      print("No todos yet!")
      return ""
    end

    var tds = ""
    for todo in todos ::
      tds = tds + displayTodo(todo)
    end
    
    return tds
  end

  # Display summary
  fn displaySummary() ::
    return "Total: #{count()}\nCompleted: #{countCompleted()}\nPending: #{countPending()}\n"
  end

end

#| Demo: Using the Todo List Manager |#

var todo1 = TodoList.create("Buy groceries", "Milk, eggs, bread, coffee")
var todo2 = TodoList.create("Finish project", "Complete the programming language implementation")
var todo3 = TodoList.create("Exercise", "Go for a 30-minute run")
var todo4 = TodoList.create("Read book", "Read chapter 5 of 'The Pragmatic Programmer'")

TodoList.complete(1)
TodoList.complete(3)

TodoList.update(2, {
  "title": "Finish language implementation",
  "description": "Complete the programming language with all features and examples"
})

var completedTodos = ""
for todo in TodoList.getCompleted() ::
  completedTodos = completedTodos + TodoList.displayTodo(todo)
end

var pendingTodos = ""
for todo in TodoList.getPending() ::
  pendingTodos = pendingTodos + TodoList.displayTodo(todo)
end

TodoList.delete(3)

TodoList.create("Learn new language feature", "Study string interpolation syntax")
TodoList.create("Write documentation", "Document all language features")

check "Todo List Manager" ::
  TodoList.count() is 5
  TodoList.displayAll() is "[✓] #1: Buy groceries - Milk, eggs, bread, coffee\n[ ] #2: Finish language implementation - Complete the programming language with all features and examples\n[ ] #4: Read book - Read chapter 5 of 'The Pragmatic Programmer'\n[ ] #5: Learn new language feature - Study string interpolation syntax\n[ ] #6: Write documentation - Document all language features\n"
  TodoList.displaySummary() is "Total: 5\nCompleted: 1\nPending: 4\n"
  TodoList.displayTodo(TodoList.getById(2)) is "[ ] #2: Finish language implementation - Complete the programming language with all features and examples\n"
  completedTodos is "[✓] #1: Buy groceries - Milk, eggs, bread, coffee\n[✓] #3: Exercise - Go for a 30-minute run\n"
  pendingTodos is "[ ] #2: Finish language implementation - Complete the programming language with all features and examples\n[ ] #4: Read book - Read chapter 5 of 'The Pragmatic Programmer'\n"
end

println("✓ All random tests passed!")
