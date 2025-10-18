## Comprehensive test file for constant immutability
## Tests deep immutability for all const-declared objects

## Test 1: Constant array index assignment should fail
fn test_assign() ::
    const arr = [1, 2, 3]
    arr[0] = 10
end

check "Constant array index assignment should fail" ::
    test_assign() raises "immutable"
end

## Test 2: Constant array push should fail
fn test_push() ::
    const arr = [1, 2, 3]
    arr.push(4)
end

check "Constant array push should fail" ::
    test_push() raises "immutable"
end

## Test 3: Constant array pop should fail
fn test_pop() ::
    const arr = [1, 2, 3]
    arr.pop()
end

check "Constant array pop should fail" ::
    test_pop() raises "immutable"
end

## Test 4: Constant array sort should fail
fn test_sort() ::
    const arr = [3, 1, 2]
    arr.sort()
end

check "Constant array sort should fail" ::
    test_sort() raises "immutable"
end

## Test 5: Constant array sort_by should fail
fn test_sort_by() ::
    const arr = [3, 1, 2]
    arr.sort_by(fn(x) :: x end)
end

check "Constant array sort_by should fail" ::
    test_sort_by() raises "immutable"
end

## Test 6: Constant array reverse should fail
fn test_reverse() ::
    const arr = [1, 2, 3]
    arr.reverse()
end

check "Constant array reverse should fail" ::
    test_reverse() raises "immutable"
end

## Test 7: Constant map index assignment should fail
fn test_map_assign() ::
    const map = {"x": 5}
    map["x"] = 10
end

check "Constant map index assignment should fail" ::
    test_map_assign() raises "immutable"
end

## Test 8: Constant map adding key should fail
fn test_map_add_key() ::
    const map = {"x": 5}
    map["y"] = 10
end

check "Constant map adding key should fail" ::
    test_map_add_key() raises "immutable"
end

## Test 9: Regular var array should still be mutable
check "Regular var array should still be mutable" ::
    var arr = [1, 2, 3]
    var _ = arr[0] = 10
    arr[0] is 10
end

## Test 10: Regular var array push should work
check "Regular var array push should work" ::
    var arr = [1, 2, 3]
    var _ = arr.push(4)
    arr.length() is 4
    arr[3] is 4
end

## Test 11: Regular var array pop should work
check "Regular var array pop should work" ::
    var arr = [1, 2, 3]
    var popped = arr.pop()
    popped is 3
    arr.length() is 2
end

## Test 12: Regular var array sort should work
check "Regular var array sort should work" ::
    var arr = [3, 1, 2]
    var _ = arr.sort()
    arr[0] is 1
    arr[1] is 2
    arr[2] is 3
end

## Test 13: Regular var array reverse should work
check "Regular var array reverse should work" ::
    var arr = [1, 2, 3]
    var _ = arr.reverse()
    arr[0] is 3
    arr[1] is 2
    arr[2] is 1
end

## Test 14: Regular var map should still be mutable
check "Regular var map should still be mutable" ::
    var map = {"x": 5}
    var _ = map["x"] = 10
    map["x"] is 10
end

## Test 15: Regular var map adding key should work
check "Regular var map adding key should work" ::
    var map = {"x": 5}
    var _ = map["y"] = 20
    map["y"] is 20
end

## Test 16: Const primitive reassignment should fail (existing behavior)
fn test_reassign() ::
    const x = 5
    x = 10
end

check "Const primitive reassignment should fail" ::
    test_reassign() raises "constant"
end

## Test 17: Deep immutability - nested array modification should fail
fn test_nested_array() ::
    const outer = [[1, 2], [3, 4]]
    outer[0][0] = 10
end

check "Deep immutability - nested array modification should fail" ::
    test_nested_array() raises "immutable"
end

## Test 18: Deep immutability - nested array push should fail
fn test_nested_push() ::
    const outer = [[1, 2], [3, 4]]
    outer[0].push(5)
end

check "Deep immutability - nested array push should fail" ::
    test_nested_push() raises "immutable"
end

## Test 19: Deep immutability - nested map in array modification should fail
fn test_nested_map_in_array() ::
    const outer = [{"x": 5}, {"y": 10}]
    outer[0]["x"] = 20
end

check "Deep immutability - nested map in array modification should fail" ::
    test_nested_map_in_array() raises "immutable"
end

## Test 20: Deep immutability - nested array in map modification should fail
fn test_nested_array_in_map() ::
    const outer = {"arr": [1, 2, 3]}
    outer["arr"][0] = 10
end

check "Deep immutability - nested array in map modification should fail" ::
    test_nested_array_in_map() raises "immutable"
end

## Test 21: Deep immutability - nested array in map push should fail
fn test_nested_array_push() ::
    const outer = {"arr": [1, 2, 3]}
    outer["arr"].push(4)
end

check "Deep immutability - nested array in map push should fail" ::
    test_nested_array_push() raises "immutable"
end

## Test 22: Deep immutability - three levels deep
fn test_three_levels() ::
    const outer = [[[1, 2]], [[3, 4]]]
    outer[0][0][0] = 10
end

check "Deep immutability - three levels deep" ::
    test_three_levels() raises "immutable"
end

## Test 23: Non-mutating array methods should work on const arrays
check "Non-mutating array methods should work on const arrays" ::
    const arr = [1, 2, 3, 4, 5]
    var first = arr.first()
    var last = arr.last()
    var rest = arr.rest()

    first is 1
    last is 5
    rest.length() is 4
end

## Test 24: Array.map should work on const arrays (returns new array)
check "Array.map should work on const arrays" ::
    const arr = [1, 2, 3]
    var doubled = arr.map(fn(x) :: x * 2 end)

    doubled[0] is 2
    doubled[1] is 4
    doubled[2] is 6
end

## Test 25: Array.filter should work on const arrays (returns new array)
check "Array.filter should work on const arrays" ::
    const arr = [1, 2, 3, 4, 5]
    var evens = arr.filter(fn(x) :: x % 2 == 0 end)

    evens.length() is 2
    evens[0] is 2
    evens[1] is 4
end

## Test 26: Const array reading is allowed
check "Const array reading is allowed" ::
    const arr = [10, 20, 30]
    arr[0] is 10
    arr[1] is 20
    arr[2] is 30
end

## Test 27: Const map reading is allowed
check "Const map reading is allowed" ::
    const map = {"x": 100, "y": 200}
    map["x"] is 100
    map["y"] is 200
end

println("All const immutability tests completed!")
