println("Testing Array methods...")

var arr = [3, 1, 4, 1, 5, 9]

# =====================================
# Existing Array Methods Tests
# =====================================

check "Existing Array methods" ::
  arr.length() is 6
  arr.first() is 3
  arr.last() is 9

  var arr2 = [1, 2, 3]
  var _ = arr2.push(4)
  arr2.length() is 4
  arr2.last() is 4

  var rest = arr.rest()
  rest.length() is 5
  rest.first() is 1
end

# =====================================
# Functional Operations
# =====================================

check "Array.map() method" ::
  var doubled = arr.map(fn(x) :: x * 2 end)

  doubled[0] is 6
  doubled[1] is 2
  doubled[2] is 8
  doubled.length() is 6
end

check "Array.filter() method" ::
  var filtered = arr.filter(fn(x) :: x > 3 end)

  filtered.length() is 3
  filtered[0] is 4
  filtered[1] is 5
  filtered[2] is 9
end

check "Array.reduce() method" ::
  var sum = arr.reduce(fn(acc, x) :: acc + x end, 0)
  sum is 23

  var product = [2, 3, 4].reduce(fn(acc, x) :: acc * x end, 1)
  product is 24
end

check "Array.each() method" ::
  var nums = [1, 2, 3]
  var side_effect_ran = false
  # Note: each() is for side effects, but closures can't modify outer scope
  # This test just verifies each() runs without error
  var result = nums.each(fn(x) :: x * 2 end)
  isNull(result) isTrue
end

check "Array.map_with_index() method" ::
  var result = [10, 20, 30].map_with_index(fn(val, idx) :: val + idx end)

  result[0] is 10
  result[1] is 21
  result[2] is 32
end

# =====================================
# Finding Methods
# =====================================

check "Array.find() method" ::
  var found = arr.find(fn(x) :: x > 3 end)
  found is 4

  var notFound = arr.find(fn(x) :: x > 100 end)
  isNull(notFound) isTrue
end

check "Array.find_index() method" ::
  var idx = arr.find_index(fn(x) :: x > 3 end)
  idx is 2

  var notFoundIdx = arr.find_index(fn(x) :: x > 100 end)
  notFoundIdx is (-1)
end

check "Array.any(), all(), none() methods" ::
  arr.any(fn(x) :: x > 8 end) isTrue
  arr.all(fn(x) :: x > 0 end) isTrue
  arr.none(fn(x) :: x > 10 end) isTrue
  arr.none(fn(x) :: x > 3 end) isFalse
end

check "Array.count() method" ::
  var count = arr.count(fn(x) :: x == 1 end)
  count is 2

  var count2 = arr.count(fn(x) :: x > 3 end)
  count2 is 3
end

# =====================================
# Transformation Methods
# =====================================

check "Array.sort() method" ::
  var nums = [3, 1, 4, 1, 5, 9]
  var _ = nums.sort()

  nums[0] is 1
  nums[1] is 1
  nums[2] is 3
  nums[3] is 4
  nums[4] is 5
  nums[5] is 9
end

check "Array.sort_by() method" ::
  var nums = [1, 2, 3, 4, 5]
  var _ = nums.sort_by(fn(x) :: -x end)

  nums[0] is 5
  nums[1] is 4
  nums[2] is 3
  nums[3] is 2
  nums[4] is 1
end

check "Array.reverse() method" ::
  var nums = [1, 2, 3]
  var _ = nums.reverse()

  nums[0] is 3
  nums[1] is 2
  nums[2] is 1
end

check "Array.unique() method" ::
  var nums = [1, 2, 2, 3, 3, 3, 4]
  var uniq = nums.unique()

  uniq.length() is 4
  uniq[0] is 1
  uniq[1] is 2
  uniq[2] is 3
  uniq[3] is 4
end

# =====================================
# Slicing/Combining
# =====================================

check "Array.slice() method" ::
  var sliced = arr.slice(1, 4)

  sliced.length() is 3
  sliced[0] is 1
  sliced[1] is 4
  sliced[2] is 1
end

check "Array.take() method" ::
  var taken = arr.take(3)

  taken.length() is 3
  taken[0] is 3
  taken[1] is 1
  taken[2] is 4
end

check "Array.drop() method" ::
  var dropped = arr.drop(2)

  dropped.length() is 4
  dropped[0] is 4
  dropped[1] is 1
  dropped[2] is 5
  dropped[3] is 9
end

check "Array.concat() method" ::
  var combined = [1, 2].concat([3, 4])

  combined.length() is 4
  combined[0] is 1
  combined[1] is 2
  combined[2] is 3
  combined[3] is 4
end

# =====================================
# Nested Arrays
# =====================================

check "Array.flatten() method" ::
  var nested = [[1, 2], [3, 4], [5]]
  var flat = nested.flatten()

  flat.length() is 5
  flat[0] is 1
  flat[2] is 3
  flat[4] is 5
end

check "Array.flat_map() method" ::
  var result = [1, 2, 3].flat_map(fn(x) :: [x, x * 2] end)

  result.length() is 6
  result[0] is 1
  result[1] is 2
  result[2] is 2
  result[3] is 4
end

# =====================================
# Membership
# =====================================

check "Array.contains() method" ::
  arr.contains(5) isTrue
  arr.contains(10) isFalse
end

check "Array.index_of() and last_index_of() methods" ::
  arr.index_of(1) is 1
  arr.last_index_of(1) is 3

  arr.index_of(99) is (-1)
end

# =====================================
# Conversion
# =====================================

check "Array.join() method" ::
  var str = [1, 2, 3].join(", ")
  str is "1, 2, 3"

  var str2 = ["a", "b", "c"].join("-")
  str2 is "a-b-c"
end

# =====================================
# Statistics
# =====================================

check "Array.sum() method" ::
  var total = arr.sum()
  total is 23

  var total2 = [10, 20, 30].sum()
  total2 is 60
end

check "Array.average() method" ::
  var avg = [2, 4, 6].average()
  avg is 4

  var avg2 = [10, 20, 30].average()
  avg2 is 20
end

check "Array.min() and max() methods" ::
  arr.min() is 1
  arr.max() is 9

  var nums2 = [5, 2, 8, 1, 9]
  nums2.min() is 1
  nums2.max() is 9
end

# =====================================
# Grouping
# =====================================

check "Array.chunk() method" ::
  var chunked = [1, 2, 3, 4, 5].chunk(2)

  chunked.length() is 3
  chunked[0].length() is 2
  chunked[1].length() is 2
  chunked[2].length() is 1

  chunked[0][0] is 1
  chunked[0][1] is 2
  chunked[1][0] is 3
end

check "Array.partition() method" ::
  var parts = arr.partition(fn(x) :: x > 3 end)

  parts.length() is 2

  var trueGroup = parts[0]
  var falseGroup = parts[1]

  trueGroup.length() is 3
  trueGroup[0] is 4
  trueGroup[1] is 5
  trueGroup[2] is 9

  falseGroup.length() is 3
  falseGroup[0] is 3
  falseGroup[1] is 1
  falseGroup[2] is 1
end

# =====================================
# Array Operations
# =====================================

check "Array.zip() method" ::
  var zipped = [1, 2, 3].zip([10, 20, 30])

  zipped.length() is 3
  zipped[0][0] is 1
  zipped[0][1] is 10
  zipped[1][0] is 2
  zipped[1][1] is 20
end

check "Array.compact() method" ::
  var arr3 = [1, nil, 2, nil, 3]
  var compacted = arr3.compact()

  compacted.length() is 3
  compacted[0] is 1
  compacted[1] is 2
  compacted[2] is 3
end

# =====================================
# Method Chaining
# =====================================

check "Array method chaining" ::
  var result = [1, 2, 3, 4, 5, 6]
    .filter(fn(x) :: x > 2 end)
    .map(fn(x) :: x * 2 end)
    .take(2)

  result.length() is 2
  result[0] is 6
  result[1] is 8
end

check "Complex array pipeline" ::
  var nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  var result = nums
    .filter(fn(x) :: x % 2 == 0 end)
    .map(fn(x) :: x * x end)
    .sum()

  result is 220
end

println("✓ All Array method tests passed!")
