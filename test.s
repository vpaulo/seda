# Single line comment

var x = 5
const name = "hello"
fn add(a, b) :: a + b end

#| This is a multiline comment
that spans multiple lines
and should be ignored |#


# Single line comment
var z = 5
#| Multiline
   comment |#
var y = 10 # Another single line


fn add(x: number, y: number): number ::
	x + y
where ::
	add(2, 3) is 5
	add(0, 0) is 0
end

println(add(4, 44))

check ::

  add(4,44) is 48

end
