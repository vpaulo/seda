# Test Suite

Comprehensive test suite for the PL programming language.

## Running Tests

### Run All Tests
```bash
# Run all test files
./seda examples/basic.s
./seda examples/type_system.s
./seda examples/custom_properties.s
./seda examples/control_flow.s
# Or
go test ./...
```

### Run Individual Tests
```bash
# Run a specific test file
./seda examples/basic.s
```

## Test Files

### `basic.s`
Core language features:
- Variables and constants
- Arithmetic operations
- String operations
- Boolean operations
- Comparison operations
- Array operations
- Map operations
- Functions
- Closures

### `type_system.s`
Type system features:
- Type aliases
- Type checking with primitives
- Type annotations on variables
- Complex type aliases
- isA operator

### `custom_properties.s`
Custom properties system:
- String custom properties
- Number custom properties
- Boolean custom properties
- Array custom properties
- Map custom properties
- Factory function pattern

### `control_flow.s`
Control flow constructs:
- If-else statements
- Nested if-else
- For loops with arrays
- For loops with maps

## Test Output

Each test file outputs progress and results:
```
Running basic language tests...
✓ All basic tests passed!
```

All tests use the built-in `check` blocks for assertions:
```seda
check "test description" ::
  actual is expected
end
```

## Adding New Tests

To add new tests:

1. Create a new `.s` file in the `examples/` directory
2. Use descriptive names: `feature_name.s`
3. Include a print statement at the start: `println("Running feature tests...")`
4. Group related tests in `check` blocks
5. Add a success message at the end: `println("✓ All feature tests passed!")`

Example structure:
```seda
println("Running my feature tests...")

check "feature description" ::
  # test assertions here
  result is expected
end

println("✓ All my feature tests passed!")
```
