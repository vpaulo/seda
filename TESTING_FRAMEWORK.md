# Testing Framework Design

## Overview
The language provides built-in testing capabilities through `check::` and `where::` blocks, making testing a first-class language feature rather than an external library concern.

## Core Testing Constructs

### `check::` Blocks
Standalone test blocks that can appear anywhere in the code:

```
check ::
  expression is expected_value
  function_call is expected_result
  variable_name is expected_state
end
```

**Characteristics:**
- Can be placed at module level, inside functions, or in any scope
- Execute immediately when encountered during program execution
- Provide isolated test environments
- Support multiple assertions per block

### `where::` Blocks
Function-specific test blocks attached to function definitions:

```
fn functionName(params) ::
  function_body
where ::
  functionName(input1) is output1
  functionName(input2) is output2
end
```

**Characteristics:**
- Execute immediately after function definition
- Test the specific function they're attached to
- Provide documentation-like examples
- Create a clear association between function and its expected behavior

## Assertion Syntax

### Basic Assertions
```
# Equality assertion
expression is expected_value

# Examples
2 + 2 is 4
"hello" + " world" is "hello world"
array.length is 3
```

### Extended Assertion Types
1. is - Equality check (e.g., x is 10)
2. isNot - Inequality check (e.g., x isNot y)
3. isA - Type checking (e.g., num isA "NUMBER")
4. contains - Membership check for arrays/strings (e.g., arr contains 3)
5. isGreater - Numeric greater than (e.g., x isGreater 5)
6. isLess - Numeric less than (e.g., x isLess 10)
7. isTrue - Boolean true check (e.g., flag isTrue)
8. isFalse - Boolean false check (e.g., flag isFalse)
9. isEmpty - Empty array/string check (e.g., empty_arr isEmpty)
10. startsWith - String prefix check (e.g., text startsWith "Hello")
11. endsWith - String suffix check (e.g., filename endsWith ".pdf")
12. raises - Error/exception check (e.g., (10 / 0) raises "division by zero")
```
# Type checking
value isA Number
obj isA CustomType

# Comparison assertions
number isGreater 10
number isLess 100

# Boolean assertions
condition isTrue
condition isFalse

# Array/Collection assertions
array contains element
array isEmpty

# String assertions
text startsWith "Hello"
text endsWith "world"

# Approximate equality (for floating point)
calculation isCloseTo 3.14159, tolerance: 0.001

# Exception assertions
expression raises "error message"
```

## Test Execution Model

### Execution Timing
1. **`where::` blocks**: Execute immediately after function definition
2. **`check::` blocks**: Execute when encountered in program flow
3. **Module-level tests**: Execute after module loading
4. **Function-level tests**: Execute when function scope is entered

### Test Environment
Each test block creates an isolated environment:
- Local variables don't leak between assertions
- Function calls are tracked for side effects
- Test state is reset between assertions

### Error Handling
```
check ::
  # This assertion fails but doesn't stop other tests
  1 + 1 is 3  # FAIL: Expected 3, got 2

  # This assertion continues to run
  2 + 2 is 4  # PASS
end
```

## Test Result Reporting

### Basic Output Format
```
✓ PASS: 2 + 2 is 4
✗ FAIL: 1 + 1 is 3
  Expected: 3
  Actual: 2
  Location: example.lang:5:7

✓ PASS: array.length is 3
```

### Detailed Reporting
```
Test Results Summary:
==================
File: example.lang
Tests Run: 15
Passed: 12
Failed: 3
Duration: 23ms

Failures:
---------
1. Line 5: 1 + 1 is 3
   Expected: 3
   Actual: 2

2. Line 12: user.name is "Alice"
   Expected: "Alice"
   Actual: "Bob"

3. Line 18: items contains "missing"
   Expected: array to contain "missing"
   Actual: ["item1", "item2", "item3"]
```

## Advanced Testing Features

### Test Context and Setup
```
check ::
  # Setup code
  var user = User.new("Alice", 25)
  var items = ["apple", "banana", "cherry"]

  # Tests use the setup
  user.name is "Alice"
  user.age is 25
  items.length is 3
  items.first is "apple"
end
```

### Parameterized Tests
```
fn factorial(n: Number): Number ::
  if n <= 1 :: 1 else :: n * factorial(n - 1) end
where ::
  for input, expected in [[0, 1], [1, 1], [2, 2], [3, 6], [4, 24]] ::
    factorial(input) is expected
  end
end
```

### Test Groups and Organization
```
check "Basic arithmetic tests" ::
  2 + 2 is 4
  3 * 4 is 12
  10 / 2 is 5
end

check "String manipulation tests" ::
  "hello".toUpperCase() is "HELLO"
  "world".length is 5
  "test".contains("es") is true
end
```

### Conditional Tests
```
check ::
  if platform == "windows" ::
    path.separator is "\\"
  else ::
    path.separator is "/"
  end
end
```