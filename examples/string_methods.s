println("Testing String methods...")

# =====================================
# Existing String Methods Tests
# =====================================

check "String.length() method" ::
  var str1 = "hello"
  var str2 = ""
  var str3 = "hello world"

  str1.length() is 5
  str2.length() is 0
  str3.length() is 11
end

check "String.upper() method" ::
  var str1 = "hello"
  var str2 = "Hello World"
  var str3 = "ALREADY UPPER"

  str1.upper() is "HELLO"
  str2.upper() is "HELLO WORLD"
  str3.upper() is "ALREADY UPPER"
end

check "String.lower() method" ::
  var str1 = "HELLO"
  var str2 = "Hello World"
  var str3 = "already lower"

  str1.lower() is "hello"
  str2.lower() is "hello world"
  str3.lower() is "already lower"
end

check "String.substr() method" ::
  var str = "hello world"

  str.substr(0, 5) is "hello"
  str.substr(6, 5) is "world"
  str.substr(0, 11) is "hello world"
  str.substr(6, 3) is "wor"
end

# =====================================
# New String Methods Tests (Phase 2)
# =====================================

check "String.split() method" ::
  var str1 = "hello,world,test"
  var parts1 = str1.split(",")

  parts1.length() is 3
  parts1[0] is "hello"
  parts1[1] is "world"
  parts1[2] is "test"

  var str2 = "one two three"
  var parts2 = str2.split(" ")

  parts2.length() is 3
  parts2[0] is "one"
  parts2[1] is "two"
  parts2[2] is "three"
end

check "String.split() with empty delimiter" ::
  var str = "abc"
  var chars = str.split("")

  chars.length() is 3
  chars[0] is "a"
  chars[1] is "b"
  chars[2] is "c"
end

check "String.trim() method" ::
  var str1 = "  hello  "
  var str2 = "\t\nworld\t\n"
  var str3 = "no spaces"

  str1.trim() is "hello"
  str2.trim() is "world"
  str3.trim() is "no spaces"
end

check "String.replace() method" ::
  var str1 = "hello world"
  var str2 = "foo bar foo"
  var str3 = "test"

  str1.replace("world", "universe") is "hello universe"
  str2.replace("foo", "baz") is "baz bar baz"
  str3.replace("x", "y") is "test"
end

check "String.starts_with() method" ::
  var str = "hello world"

  str.starts_with("hello") isTrue
  str.starts_with("world") isFalse
  str.starts_with("h") isTrue
  str.starts_with("") isTrue
end

check "String.ends_with() method" ::
  var str = "hello world"

  str.ends_with("world") isTrue
  str.ends_with("hello") isFalse
  str.ends_with("d") isTrue
  str.ends_with("") isTrue
end

check "String.index_of() method" ::
  var str = "hello world"

  str.index_of("world") is 6
  str.index_of("hello") is 0
  str.index_of("l") is 2
  str.index_of("xyz") is (-1)
end

check "String.char_at() method" ::
  var str = "hello"

  str.char_at(0) is "h"
  str.char_at(1) is "e"
  str.char_at(4) is "o"
end

check "String.char_at() out of bounds" ::
  var str = "hello"

  str.char_at(10) raises "out of bounds"
  str.char_at(-1) raises "out of bounds"
end

# =====================================
# Method Chaining Tests
# =====================================

check "String method chaining" ::
  var str = "  Hello World  "
  var result = str.trim().lower()

  result is "hello world"
end

check "Complex string operations" ::
  var str = "  foo,bar,baz  "
  var parts = str.trim().split(",")

  parts.length() is 3
  parts[0] is "foo"
  parts[1] is "bar"
  parts[2] is "baz"

  var uppercase = parts[0].upper()
  uppercase is "FOO"
end

check "String manipulation pipeline" ::
  var text = "Hello World"
  var modified = text.lower().replace("world", "universe")

  modified is "hello universe"

  var has_hello = modified.starts_with("hello")
  has_hello isTrue
end

check "Split and join pattern" ::
  var csv = "apple,banana,cherry"
  var fruits = csv.split(",")
  var first = fruits[0]
  var last = fruits[2]

  first is "apple"
  last is "cherry"
  first.upper() is "APPLE"
end

# =====================================
# Additional String Methods Tests
# =====================================

check "String.trim_left() and trim_right() methods" ::
  var str = "  hello world  "

  str.trim_left() is "hello world  "
  str.trim_right() is "  hello world"

  var str2 = "\t\n test \t\n"
  str2.trim_left() is "test \t\n"
  str2.trim_right() is "\t\n test"
end

check "String.contains() method" ::
  var str = "hello world"

  str.contains("world") isTrue
  str.contains("hello") isTrue
  str.contains("xyz") isFalse
  str.contains("") isTrue
end

check "String.last_index_of() method" ::
  var str = "hello world hello"

  str.last_index_of("hello") is 12
  str.last_index_of("l") is 15
  str.last_index_of("xyz") is (-1)
end

check "String.count() method" ::
  var str = "hello world"

  str.count("l") is 3
  str.count("o") is 2
  str.count("xyz") is 0
  str.count("ll") is 1
end

check "String.replace_first() method" ::
  var str = "foo bar foo"

  str.replace_first("foo", "baz") is "baz bar foo"
  str.replace_first("bar", "qux") is "foo qux foo"
end

check "String.reverse() method" ::
  var str1 = "hello"
  var str2 = "abc"
  var str3 = ""

  str1.reverse() is "olleh"
  str2.reverse() is "cba"
  str3.reverse() is ""
end

check "String.repeat() method" ::
  var str = "ab"

  str.repeat(3) is "ababab"
  str.repeat(1) is "ab"
  str.repeat(0) is ""
end

check "String.lines() method" ::
  var text = "line1\nline2\nline3"
  var lines = text.lines()

  lines.length() is 3
  lines[0] is "line1"
  lines[1] is "line2"
  lines[2] is "line3"
end

check "String.chars() method" ::
  var str = "abc"
  var chars = str.chars()

  chars.length() is 3
  chars[0] is "a"
  chars[1] is "b"
  chars[2] is "c"
end

check "String.words() method" ::
  var text = "  hello   world  test  "
  var words = text.words()

  words.length() is 3
  words[0] is "hello"
  words[1] is "world"
  words[2] is "test"
end

check "String.capitalize() method" ::
  var str1 = "hello world"
  var str2 = "test"

  str1.capitalize() is "Hello world"
  str2.capitalize() is "Test"
end

check "String.title_case() method" ::
  var str1 = "hello world"
  var str2 = "the quick brown fox"

  str1.title_case() is "Hello World"
  str2.title_case() is "The Quick Brown Fox"
end

check "String.is_empty() method" ::
  var str1 = ""
  var str2 = "hello"

  str1.is_empty() isTrue
  str2.is_empty() isFalse
end

check "String.is_blank() method" ::
  var str1 = "   "
  var str2 = "hello"
  var str3 = ""

  str1.is_blank() isTrue
  str2.is_blank() isFalse
  str3.is_blank() isTrue
end

check "String.is_numeric() method" ::
  var str1 = "12345"
  var str2 = "123abc"
  var str3 = ""

  str1.is_numeric() isTrue
  str2.is_numeric() isFalse
  str3.is_numeric() isFalse
end

check "String.is_alpha() method" ::
  var str1 = "hello"
  var str2 = "hello123"
  var str3 = ""

  str1.is_alpha() isTrue
  str2.is_alpha() isFalse
  str3.is_alpha() isFalse
end

check "String.pad_left() method" ::
  var str = "hello"

  str.pad_left(10, "-") is "-----hello"
  str.pad_left(8, "x") is "xxxhello"
  str.pad_left(5, "-") is "hello"
  str.pad_left(3, "-") is "hello"
end

check "String.pad_right() method" ::
  var str = "hello"

  str.pad_right(10, "-") is "hello-----"
  str.pad_right(8, "x") is "helloxxx"
  str.pad_right(5, "-") is "hello"
  str.pad_right(3, "-") is "hello"
end

println("✓ All String method tests passed!")
