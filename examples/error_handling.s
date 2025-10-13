println("Running error handling tests...")

# Function that returns value and error
fn divide(a, b) ::
    if b == 0 ::
        return nil, error("division by zero")
    end
    return a / b, nil
end

check "Using multiple return values with success case" ::
  var result, err = divide(10, 2)

  result is 5
  err is nil
end

check "Using multiple return values with error case" ::
  var result, err = divide(10, 0)

  result is nil
  err.to_string is "division by zero"
end

# Function that propagates errors
fn safe_divide(a, b) ::
    var result, err = divide(a, b)
    if !isNull(err) ::
        return nil, err  # Propagate error
    end
    return result * 2, nil  # Transform successful result
end

check "Function that propagates errors, success case" ::
  var result, err = safe_divide(20, 4)

  result is 10
  err is nil
end

check "Function that propagates errors, fail case" ::
  var result, err = safe_divide(20, 0)

  result is nil
  err.to_string is "division by zero"
end

println("✓ All error handling tests passed!")
