# File Module Comprehensive Tests
println("Testing File Module...")

# Test File.join - path manipulation
check "File.join - path segments" ::
  var path1 = File.join("/path", "to", "file.txt")
  path1 is "/path/to/file.txt"

  var path2 = File.join("home", "user", "documents")
  path2 contains "home"
  path2 contains "user"
  path2 contains "documents"
end

# Test File.basename - get base name
check "File.basename - extract filename" ::
  var base1 = File.basename("/path/to/file.txt")
  base1 is "file.txt"

  var base2 = File.basename("/home/user/document.md")
  base2 is "document.md"
end

# Test File.dirname - get directory
check "File.dirname - extract directory" ::
  var dir1 = File.dirname("/path/to/file.txt")
  dir1 is "/path/to"

  var dir2 = File.dirname("/home/user/document.md")
  dir2 is "/home/user"
end

# Test File.extname - get extension
check "File.extname - extract extension" ::
  var ext1 = File.extname("/path/to/file.txt")
  ext1 is ".txt"

  var ext2 = File.extname("document.md")
  ext2 is ".md"

  var ext3 = File.extname("noextension")
  ext3 is ""
end

# Test File.cwd - get current directory
check "File.cwd - current directory" ::
  var cwd = File.cwd()
  cwd isA "string"
  cwd contains "seda"
end

# Test File.absolute_path - get absolute path
check "File.absolute_path - resolve path" ::
  var abs = File.absolute_path(".")
  abs isA "string"
  abs contains "/"
end

# Test File write and read operations
check "File.write and File.read - basic operations" ::
  var test_file = "/tmp/seda_test_file.txt"
  var content = "Hello, File Module!"

  # Write content to file
  var write_err = File.write(test_file, content)
  isNull(write_err) isTrue

  # Read content back
  var read_content, read_err = File.read(test_file)
  isNull(read_err) isTrue
  read_content is content

  # Clean up
  var _ = File.delete(test_file)
end

# Test File.read_lines - read file as lines
check "File.read_lines - read as array" ::
  var test_file = "/tmp/seda_test_lines.txt"
  var lines_content = "line1\nline2\nline3"

  # Write multi-line content
  var _ = File.write(test_file, lines_content)

  # Read as lines
  var lines, err = File.read_lines(test_file)
  isNull(err) isTrue
  lines.length() is 3
  lines[0] is "line1"
  lines[1] is "line2"
  lines[2] is "line3"

  # Clean up
  var _ = File.delete(test_file)
end

# Test File.append - append to file
check "File.append - append content" ::
  var test_file = "/tmp/seda_test_append.txt"

  # Write initial content
  var _ = File.write(test_file, "First line\n")

  # Append content
  var append_err = File.append(test_file, "Second line\n")
  isNull(append_err) isTrue

  # Read back and verify
  var content, _ = File.read(test_file)
  content contains "First line"
  content contains "Second line"

  # Clean up
  var _ = File.delete(test_file)
end

# Test File.exists - check existence
check "File.exists - file existence" ::
  var test_file = "/tmp/seda_test_exists.txt"

  # File should not exist initially
  var exists_before = File.exists(test_file)
  exists_before isFalse

  # Create file
  var _ = File.write(test_file, "test")

  # File should exist now
  var exists_after = File.exists(test_file)
  exists_after isTrue

  # Clean up
  var _ = File.delete(test_file)
end

# Test File.size - get file size
check "File.size - file size in bytes" ::
  var test_file = "/tmp/seda_test_size.txt"
  var content = "12345"

  # Write content
  var _ = File.write(test_file, content)

  # Get size
  var size, err = File.size(test_file)
  isNull(err) isTrue
  size is 5

  # Clean up
  var _ = File.delete(test_file)
end

# Test File.is_file and File.is_dir
check "File.is_file and File.is_dir - type checking" ::
  var test_file = "/tmp/seda_test_is_file.txt"
  var test_dir = "/tmp/seda_test_dir"

  # Create file and directory
  var _ = File.write(test_file, "test")
  var _ = File.mkdir(test_dir)

  # Test file type
  var is_file = File.is_file(test_file)
  is_file isTrue

  var is_dir1 = File.is_dir(test_file)
  is_dir1 isFalse

  # Test directory type
  var is_dir2 = File.is_dir(test_dir)
  is_dir2 isTrue

  var is_file2 = File.is_file(test_dir)
  is_file2 isFalse

  # Clean up
  var _ = File.delete(test_file)
  var _ = File.remove_dir(test_dir)
end

# Test File.mkdir and File.list_dir
check "File.mkdir and File.list_dir - directory operations" ::
  var test_dir = "/tmp/seda_test_list_dir"

  # Create directory
  var mkdir_err = File.mkdir(test_dir)
  isNull(mkdir_err) isTrue

  # Create some files in the directory
  var _ = File.write(File.join(test_dir, "file1.txt"), "content1")
  var _ = File.write(File.join(test_dir, "file2.txt"), "content2")
  var _ = File.write(File.join(test_dir, "file3.txt"), "content3")

  # List directory contents
  var files, err = File.list_dir(test_dir)
  isNull(err) isTrue
  files.length() is 3

  # Check if files are in the list
  var files_str = files.join(",")
  files_str contains "file1.txt"
  files_str contains "file2.txt"
  files_str contains "file3.txt"

  # Clean up
  var _ = File.delete(File.join(test_dir, "file1.txt"))
  var _ = File.delete(File.join(test_dir, "file2.txt"))
  var _ = File.delete(File.join(test_dir, "file3.txt"))
  var _ = File.remove_dir(test_dir)
end

# Test File.mkdir_all - create nested directories
check "File.mkdir_all - nested directories" ::
  var nested_dir = "/tmp/seda_test/nested/deeply/nested"

  # Create all directories at once
  var mkdir_all_err = File.mkdir_all(nested_dir)
  isNull(mkdir_all_err) isTrue

  # Verify directory exists
  var exists = File.is_dir(nested_dir)
  exists isTrue

  # Clean up (remove_dir removes all subdirectories)
  var _ = File.remove_dir("/tmp/seda_test")
end

# Test File.delete - delete file
check "File.delete - remove file" ::
  var test_file = "/tmp/seda_test_delete.txt"

  # Create file
  var _ = File.write(test_file, "test")
  var exists_before = File.exists(test_file)
  exists_before isTrue

  # Delete file
  var delete_err = File.delete(test_file)
  isNull(delete_err) isTrue

  # Verify file is deleted
  var exists_after = File.exists(test_file)
  exists_after isFalse
end

# Test File.remove_dir - remove directory tree
check "File.remove_dir - remove directory tree" ::
  var test_dir = "/tmp/seda_test_remove"

  # Create directory with files
  var _ = File.mkdir_all(File.join(test_dir, "subdir"))
  var _ = File.write(File.join(test_dir, "file.txt"), "test")
  var _ = File.write(File.join(test_dir, "subdir", "file2.txt"), "test2")

  # Remove entire directory tree
  var remove_err = File.remove_dir(test_dir)
  isNull(remove_err) isTrue

  # Verify directory is removed
  var exists = File.exists(test_dir)
  exists isFalse
end

# Test error handling - read non-existent file
check "Error handling - non-existent file" ::
  var content, err = File.read("/tmp/nonexistent_file_xyz.txt")
  isNull(content) isTrue
  !isNull(err) isTrue
  err.to_string() contains "no such file"
end

# Test error handling - list non-existent directory
check "Error handling - non-existent directory" ::
  var files, err = File.list_dir("/tmp/nonexistent_dir_xyz")
  isNull(files) isTrue
  !isNull(err) isTrue
end

println("All File module tests completed!")
