package object

import (
	"fmt"
	"strings"

	"github.com/vpaulo/seda/ast"
)

// ObjectType represents the type of an object
type ObjectType string

const (
	// Basic types
	NUMBER_OBJ  = "NUMBER"
	STRING_OBJ  = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"

	// Collection types
	ARRAY_OBJ = "ARRAY"
	MAP_OBJ   = "MAP"
	RANGE_OBJ = "RANGE"

	// Function types
	FUNCTION_OBJ = "FUNCTION"
	METHOD_OBJ   = "METHOD"
	BUILTIN_OBJ  = "BUILTIN"
	MODULE_OBJ   = "MODULE"

	// Control flow
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	BREAK_OBJ        = "BREAK"

	// Testing types
	TEST_RESULT_OBJ = "TEST_RESULT"

	// Special types
	NULL_OBJ       = "NULL"
	ERROR_OBJ      = "ERROR"
	TYPE_ALIAS_OBJ = "TYPE_ALIAS"
)

// Object represents any value in the language
type Object interface {
	Type() ObjectType
	Inspect() string
	String() string
}

// Number represents a numeric value
type Number struct {
	Value      float64
	Properties map[string]Object // Custom properties/methods
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string  { return fmt.Sprintf("%.10g", n.Value) }
func (n *Number) String() string   { return n.Inspect() }

// String represents a string value
type String struct {
	Value      string
	Properties map[string]Object // Custom properties/methods
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf("\"%s\"", s.Value) }
func (s *String) String() string   { return s.Value }

// Boolean represents a boolean value
type Boolean struct {
	Value      bool
	Properties map[string]Object // Custom properties/methods
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) String() string   { return b.Inspect() }

// Array represents an array of objects
type Array struct {
	Elements   []Object
	Properties map[string]Object // Custom properties/methods
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
func (a *Array) String() string { return a.Inspect() }

// Range represents a range of numbers
type Range struct {
	Start     int
	End       int
	Inclusive bool
}

func (r *Range) Type() ObjectType { return RANGE_OBJ }
func (r *Range) Inspect() string {
	if r.Inclusive {
		return fmt.Sprintf("%d...%d", r.Start, r.End)
	}
	return fmt.Sprintf("%d..%d", r.Start, r.End)
}
func (r *Range) String() string { return r.Inspect() }

// MapPair represents a key-value pair in a map
type MapPair struct {
	Key   Object
	Value Object
}

// Map represents a map/hash/dictionary
type Map struct {
	Pairs map[string]MapPair
	Properties map[string]Object // Custom properties/methods
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var pairs []string
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}
func (m *Map) String() string { return m.Inspect() }

// Null represents a null/nil value
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) String() string   { return "null" }

// Error represents an error object
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return fmt.Sprintf("ERROR: %s", e.Message) }
func (e *Error) String() string   { return e.Message }

// ReturnValue wraps other objects when returned from functions
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) String() string   { return rv.Value.String() }

// Break signals a break from a loop
type Break struct {
}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }
func (b *Break) String() string   { return "break" }

// Global constants for common values
var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

// Helper functions for type checking
func IsTruthy(obj Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

// NewError creates a new error object
func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// Environment represents a variable binding environment
type Environment struct {
	store            map[string]Object
	constants        map[string]bool // Track which identifiers are constants
	outer            *Environment
	InWhereBlockTest bool // Flag to prevent infinite recursion in where block tests
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Environment{store: s, constants: c, outer: nil}
}

// NewEnclosedEnvironment creates a new environment with an outer scope
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a value from the environment
func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

// IsInWhereBlockTest checks if this environment or any outer environment is in a where block test
func (e *Environment) IsInWhereBlockTest() bool {
	if e.InWhereBlockTest {
		return true
	}
	if e.outer != nil {
		return e.outer.IsInWhereBlockTest()
	}
	return false
}

// Set stores a value in the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// SetConstant stores a constant value in the environment
func (e *Environment) SetConstant(name string, val Object) Object {
	e.store[name] = val
	e.constants[name] = true
	return val
}

// IsConstant checks if a name is a constant in this environment or outer scopes
func (e *Environment) IsConstant(name string) bool {
	if e.constants[name] {
		return true
	}
	if e.outer != nil {
		return e.outer.IsConstant(name)
	}
	return false
}

// Update updates an existing variable in the environment chain, or creates it in current scope if not found
func (e *Environment) Update(name string, val Object) Object {
	// Check if it's a constant in current scope
	if e.constants[name] {
		return NewError("cannot reassign constant '%s'", name)
	}

	// Check if variable exists in current scope
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val
	}

	// Check outer scopes
	if e.outer != nil {
		// Check if it's a constant in outer scope
		if e.outer.IsConstant(name) {
			return NewError("cannot reassign constant '%s'", name)
		}
		if _, ok := e.outer.Get(name); ok {
			return e.outer.Update(name, val)
		}
	}

	// Variable doesn't exist anywhere, create it in current scope
	e.store[name] = val
	return val
}

// GetStore returns the internal store (for module extraction)
func (e *Environment) GetStore() map[string]Object {
	return e.store
}

// Function represents a user-defined function
type Function struct {
	Parameters []*ast.Parameter
	Body       *ast.BlockStatement
	Env        *Environment
	WhereBlock *ast.WhereBlock
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.Name.String())
	}
	return fmt.Sprintf("fn(%s) {\n%s\n}", strings.Join(params, ", "), f.Body.String())
}
func (f *Function) String() string { return f.Inspect() }

// Builtin represents a built-in function
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) String() string   { return b.Inspect() }

// Module represents a module with its own namespace
type Module struct {
	Name        string
	Environment *Environment
}

func (m *Module) Type() ObjectType { return MODULE_OBJ }
func (m *Module) Inspect() string  { return fmt.Sprintf("module %s", m.Name) }
func (m *Module) String() string   { return m.Inspect() }

// TypeAlias represents a type alias declaration
type TypeAlias struct {
	Name           string
	TypeAnnotation *ast.TypeAnnotation
}

func (t *TypeAlias) Type() ObjectType { return TYPE_ALIAS_OBJ }
func (t *TypeAlias) Inspect() string  { return fmt.Sprintf("type %s", t.Name) }
func (t *TypeAlias) String() string   { return t.Inspect() }

// TestResult represents the result of test assertions
type TestResult struct {
	Passed   int
	Failed   int
	Failures []string
	Label    string
}

func (tr *TestResult) Type() ObjectType { return TEST_RESULT_OBJ }
func (tr *TestResult) Inspect() string {
	total := tr.Passed + tr.Failed
	status := "PASSED"
	if tr.Failed > 0 {
		status = "FAILED"
	}

	result := fmt.Sprintf("Test %s: %d/%d assertions passed", status, tr.Passed, total)
	if tr.Label != "" {
		result = fmt.Sprintf("Test '%s' %s: %d/%d assertions passed", tr.Label, status, tr.Passed, total)
	}

	if tr.Failed > 0 {
		result += "\nFailures:\n"
		for _, failure := range tr.Failures {
			result += fmt.Sprintf("  - %s\n", failure)
		}
	}

	return result
}
func (tr *TestResult) String() string { return tr.Inspect() }
