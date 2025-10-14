# OS Module Comprehensive Tests
# Tests all OS module functions

println("Testing OS Module...")

# Test OS.getenv - get environment variable
check "OS.getenv - get environment variable" ::
  var path = OS.getenv("PATH")
  path isA "string"
  # PATH should exist on all systems
  path.length() isGreater 0
end

# Test OS.setenv and getenv together
check "OS.setenv and OS.getenv - set and get variable" ::
  var _ = OS.setenv("SEDA_TEST_VAR", "test_value")
  var value = OS.getenv("SEDA_TEST_VAR")
  value is "test_value"
end

# Test OS.env - get all environment variables
check "OS.env - get all environment variables" ::
  var env_map = OS.env()
  env_map isA "map"
  # Should have PATH variable
  var path = env_map["PATH"]
  path isA "string"
  path.length() isGreater 0
end

# Test OS.args - get command line arguments
check "OS.args - get command line arguments" ::
  var args = OS.args()
  args isA "array"
  # First argument should be the script filename
  args.length() isGreater 0
  args[0] contains "os_module.s"
end

# Test OS.pid - get process ID
check "OS.pid - get process ID" ::
  var pid = OS.pid()
  pid isA "number"
  pid isGreater 0
end

# Test OS.platform - get operating system
check "OS.platform - get operating system" ::
  var platform = OS.platform()
  platform isA "string"
  # Should be one of: linux, darwin, windows
  var valid_platforms = ["linux", "darwin", "windows"]
  var is_valid = valid_platforms.contains(platform)
  is_valid isTrue
end

# Test OS.arch - get architecture
check "OS.arch - get architecture" ::
  var arch = OS.arch()
  arch isA "string"
  # Common architectures
  var valid_archs = ["amd64", "arm64", "386", "arm"]
  var has_arch = valid_archs.contains(arch)
  has_arch isTrue
end

# Test OS.hostname - get machine hostname
check "OS.hostname - get machine hostname" ::
  var hostname = OS.hostname()
  hostname isA "string"
  hostname.length() isGreater 0
end

# Test OS.home_dir - get user home directory
check "OS.home_dir - get user home directory" ::
  var home = OS.home_dir()
  home isA "string"
  home.length() isGreater 0
  home contains "/"
end

# Test OS.temp_dir - get temporary directory
check "OS.temp_dir - get temporary directory" ::
  var temp = OS.temp_dir()
  temp isA "string"
  temp.length() isGreater 0
end

# Test OS.cwd - get current working directory
check "OS.cwd - get current working directory" ::
  var cwd = OS.cwd()
  cwd isA "string"
  cwd.length() isGreater 0
  cwd contains "seda"
end

# Test OS.chdir - change current working directory
check "OS.chdir - change current working directory" ::
  var original = OS.cwd()

  # Change to /tmp
  var _ = OS.chdir("/tmp")
  var new_dir = OS.cwd()
  new_dir is "/tmp"

  # Change back to original
  var _ = OS.chdir(original)
  var back = OS.cwd()
  back is original
end

# Test OS.exec - execute command and get output
check "OS.exec - execute command" ::
  var output, err = OS.exec("echo", "hello world")
  isNull(err) isTrue
  output contains "hello world"
end

# Test OS.exec with multiple arguments
check "OS.exec - execute with multiple arguments" ::
  var output, err = OS.exec("echo", "hello", "from", "seda")
  isNull(err) isTrue
  output contains "hello"
  output contains "from"
  output contains "seda"
end

# Test OS.exec - error handling for invalid command
check "OS.exec - error handling" ::
  var output, err = OS.exec("nonexistent_command_xyz")
  isNull(output) isTrue
  !isNull(err) isTrue
  err.to_string() contains "command failed"
end

# Test OS.spawn - spawn background process
check "OS.spawn - spawn background process" ::
  var pid = OS.spawn("sleep", "0.1")
  pid isA "number"
  pid isGreater 0
end

println("All OS module tests completed!")
