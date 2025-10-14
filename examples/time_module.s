## Time Module Test Suite
## Tests Time.now() factory function and all Time methods

check "Time.now() returns current time" ::
  var t = Time.now()

  !isNull(t)
end

check "Time formatting methods" ::
  var t = Time.now()

  ## format() with custom pattern
  var formatted = t.format("YYYY-MM-DD")
  formatted.length() isGreater 0

  ## to_string() returns ISO 8601 format
  var iso = t.to_string()
  iso.contains("T")  ## ISO 8601 has 'T' separator

  ## unix() returns seconds since epoch
  var unix_time = t.unix()
  unix_time isGreater 0

  ## unix_millis() returns milliseconds since epoch
  var unix_millis_time = t.unix_millis()
  unix_millis_time isGreater unix_time  ## Millis should be larger
end

check "Time component methods" ::
  var t = Time.now()

  ## year() - should be current year (2025 or later)
  var year_val = t.year()
  year_val isGreater 2024

  ## month() - should be between 1 and 12
  var month_val = t.month()
  month_val isGreater 0
  month_val isLess 13

  ## day() - should be between 1 and 31
  var day_val = t.day()
  day_val isGreater 0
  day_val isLess 32

  ## hour() - should be between 0 and 23
  var hour_val = t.hour()
  hour_val isGreater -1
  hour_val isLess 24

  ## minute() - should be between 0 and 59
  var minute_val = t.minute()
  minute_val isGreater -1
  minute_val isLess 60

  ## second() - should be between 0 and 59
  var second_val = t.second()
  second_val isGreater -1
  second_val isLess 60

  ## weekday() - should be between 0 (Sunday) and 6 (Saturday)
  var weekday_val = t.weekday()
  weekday_val isGreater -1
  weekday_val isLess 7
end

check "Time arithmetic methods return new Time objects" ::
  var t = Time.now()

  ## add_seconds() returns new Time object
  var t2 = t.add_seconds(60)
  !isNull(t2)

  ## add_minutes() returns new Time object
  var t3 = t.add_minutes(5)
  !isNull(t3)

  ## add_hours() returns new Time object
  var t4 = t.add_hours(2)
  !isNull(t4)

  ## add_days() returns new Time object
  var t5 = t.add_days(1)
  !isNull(t5)
end

check "Time arithmetic - add_seconds()" ::
  var t1 = Time.now()
  var t2 = t1.add_seconds(60)

  ## Difference should be 60 seconds
  var diff = t2.diff(t1)
  diff is 60
end

check "Time arithmetic - add_minutes()" ::
  var t1 = Time.now()
  var t2 = t1.add_minutes(5)

  ## Difference should be 300 seconds (5 * 60)
  var diff = t2.diff(t1)
  diff is 300
end

check "Time arithmetic - add_hours()" ::
  var t1 = Time.now()
  var t2 = t1.add_hours(2)

  ## Difference should be 7200 seconds (2 * 60 * 60)
  var diff = t2.diff(t1)
  diff is 7200
end

check "Time arithmetic - add_days()" ::
  var t1 = Time.now()
  var t2 = t1.add_days(1)

  ## Difference should be 86400 seconds (24 * 60 * 60)
  var diff = t2.diff(t1)
  diff is 86400
end

check "Time comparison - diff()" ::
  var t1 = Time.now()
  var t2 = t1.add_seconds(100)
  var t3 = t1.add_seconds(-50)

  ## Forward difference
  var diff1 = t2.diff(t1)
  diff1 is 100

  ## Backward difference (negative)
  var diff2 = t3.diff(t1)
  diff2 is -50
end

check "Time comparison - is_before()" ::
  var t1 = Time.now()
  var t2 = t1.add_seconds(60)
  var t3 = t1.add_seconds(-60)

  ## t1 is before t2
  var result1 = t1.is_before(t2)
  result1 isTrue

  ## t1 is NOT before t3 (t3 is before t1)
  var result2 = t1.is_before(t3)
  result2 isFalse
end

check "Time comparison - is_after()" ::
  var t1 = Time.now()
  var t2 = t1.add_seconds(60)
  var t3 = t1.add_seconds(-60)

  ## t1 is after t3
  var result1 = t1.is_after(t3)
  result1 isTrue

  ## t1 is NOT after t2 (t2 is after t1)
  var result2 = t1.is_after(t2)
  result2 isFalse
end

check "Time format conversion - various patterns" ::
  var t = Time.now()

  ## YYYY-MM-DD format
  var date = t.format("YYYY-MM-DD")
  date.length() is 10
  date.contains("-")

  ## HH:mm:ss format
  var time_str = t.format("HH:mm:ss")
  time_str.length() is 8
  time_str.contains(":")

  ## Full datetime format
  var datetime = t.format("YYYY-MM-DD HH:mm:ss")
  datetime.length() is 19
  datetime.contains(" ")
end

check "Time chaining - arithmetic operations" ::
  var t1 = Time.now()

  ## Chain multiple arithmetic operations
  var t2 = t1.add_days(1).add_hours(2).add_minutes(30).add_seconds(45)

  ## Total difference: 1 day + 2 hours + 30 minutes + 45 seconds
  ## = 86400 + 7200 + 1800 + 45 = 95445 seconds
  var diff = t2.diff(t1)
  diff is 95445
end

check "Time immutability - original unchanged after arithmetic" ::
  var t1 = Time.now()
  var unix1 = t1.unix()

  ## Perform arithmetic operation
  var t2 = t1.add_seconds(100)

  ## Check that t1 is unchanged
  var unix1_after = t1.unix()
  unix1 is unix1_after

  ## Check that t2 is different
  var unix2 = t2.unix()
  var diff = unix2 - unix1
  diff is 100
end

check "Parse a specific date" ::
  var t1 = Time.date("DD-MM-YYYY", "20-01-2000")

  t1.to_string() is "2000-01-20T00:00:00Z"
  t1.year() is 2000
  t1.month() is 1
  t1.day() is 20
end

check "Parse a datetime" ::
  var t2 = Time.date("YYYY-MM-DD HH:mm:ss", "2025-10-14 15:30:45")

  t2.to_string is "2025-10-14T15:30:45Z"
  t2.hour is 15
  t2.minute is 30
  t2.second is 45
end

check "Create time from Unix timestamp" ::
  var unix_timestamp = 1609459200  ## 2021-01-01 00:00:00 UTC
  var t3 = Time.unix(unix_timestamp)

  t3.to_string is "2021-01-01T00:00:00Z"
  t3.year is 2021
end

check "Verify round-trip" ::
  var t4 = Time.now()
  var unix_val = t4.unix()
  var t5 = Time.unix(unix_val)

  t4.unix is t5.unix
end

println("All Time module tests completed!")
