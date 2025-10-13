package evaluator

import (
	"fmt"
	"math"
	"strings"

	"github.com/vpaulo/seda/object"
)

// call_object_method dispatches method calls based on object type
func call_object_method(receiver object.Object, method_name string, args []object.Object) object.Object {
	switch obj := receiver.(type) {
	case *object.String:
		return call_string_method(obj, method_name, args)
	case *object.Array:
		return call_array_method(obj, method_name, args)
	case *object.Number:
		return call_number_method(obj, method_name, args)
	case *object.Map:
		return call_map_method(obj, method_name, args)
	case *object.Boolean:
		return call_boolean_method(obj, method_name, args)
	case *object.Error:
		return call_error_method(obj, method_name, args)
	default:
		return object.NewError("method '%s' not found on %s", method_name, receiver.Type())
	}
}

// String Methods

func call_string_method(str *object.String, method_name string, args []object.Object) object.Object {
	switch method_name {
	case "length":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.length. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(len(str.Value))}

	case "upper":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.upper. got=%d, want=0", len(args))
		}
		return &object.String{Value: strings.ToUpper(str.Value)}

	case "lower":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.lower. got=%d, want=0", len(args))
		}
		return &object.String{Value: strings.ToLower(str.Value)}

	case "substr":
		if len(args) != 2 {
			return object.NewError("wrong number of arguments for String.substr. got=%d, want=2", len(args))
		}

		start, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("first argument to String.substr must be NUMBER, got %s", args[0].Type())
		}

		length, ok := args[1].(*object.Number)
		if !ok {
			return object.NewError("second argument to String.substr must be NUMBER, got %s", args[1].Type())
		}

		start_idx := int(start.Value)
		length_val := int(length.Value)

		if start_idx < 0 || start_idx >= len(str.Value) {
			return object.NewError("start index out of bounds")
		}

		end_idx := start_idx + length_val
		if end_idx > len(str.Value) {
			end_idx = len(str.Value)
		}

		return &object.String{Value: str.Value[start_idx:end_idx]}
	}

	// Check for instance-specific custom properties
	if result, found := check_custom_property(str.Properties, method_name, str, args); found {
		return result
	}

	// Check for user-defined methods in the global string registry
	if result, found := check_type_registry(string_registry, method_name, str, args); found {
		return result
	}

	return object.NewError("method '%s' not found on String", method_name)
}

// Array Methods

func call_array_method(arr *object.Array, method_name string, args []object.Object) object.Object {
	// First check built-in methods
	switch method_name {
	case "length":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.length. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(len(arr.Elements))}

	case "push":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.push. got=%d, want=1", len(args))
		}

		// Mutate the array in place
		arr.Elements = append(arr.Elements, args[0])
		return arr

	case "pop":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.pop. got=%d, want=0", len(args))
		}

		length := len(arr.Elements)
		if length == 0 {
			return object.NULL
		}

		return arr.Elements[length-1]

	case "first":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.first. got=%d, want=0", len(args))
		}

		if len(arr.Elements) == 0 {
			return object.NULL
		}

		return arr.Elements[0]

	case "last":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.last. got=%d, want=0", len(args))
		}

		length := len(arr.Elements)
		if length == 0 {
			return object.NULL
		}

		return arr.Elements[length-1]

	case "rest":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.rest. got=%d, want=0", len(args))
		}

		length := len(arr.Elements)
		if length == 0 {
			return &object.Array{Elements: []object.Object{}}
		}

		// Return new array with all elements except the first
		rest_elements := make([]object.Object, length-1)
		copy(rest_elements, arr.Elements[1:])
		return &object.Array{Elements: rest_elements}
	}

	// Check for instance-specific custom properties
	if result, found := check_custom_property(arr.Properties, method_name, arr, args); found {
		return result
	}

	// Check for user-defined methods in the global array registry
	if result, found := check_type_registry(array_registry, method_name, arr, args); found {
		return result
	}

	return object.NewError("method '%s' not found on Array", method_name)
}

// check_custom_property checks for custom properties on an object and handles them
// Returns (result, found) where found indicates if the property was found
func check_custom_property(properties map[string]object.Object, method_name string, receiver object.Object, args []object.Object) (object.Object, bool) {
	if properties == nil {
		return nil, false
	}

	prop, ok := properties[method_name]
	if !ok {
		return nil, false
	}

	// If it's a function, call it with the receiver instance as first argument (self)
	if fn, ok := prop.(*object.Function); ok {
		method_args := append([]object.Object{receiver}, args...)
		return apply_function_from_method(fn, method_args), true
	}

	// If it's not a function, just return the property value (for zero-arg access)
	if len(args) == 0 {
		return prop, true
	}

	return object.NewError("property '%s' is not a function", method_name), true
}

// check_type_registry checks for user-defined methods in a global type registry
// Returns (result, found) where found indicates if the method was found
func check_type_registry(registry *object.Map, method_name string, receiver object.Object, args []object.Object) (object.Object, bool) {
	if registry == nil {
		return nil, false
	}

	pair, ok := registry.Pairs[method_name]
	if !ok {
		return nil, false
	}

	// The user-defined method should receive the receiver as first argument
	method_args := append([]object.Object{receiver}, args...)

	// Call the function
	switch function := pair.Value.(type) {
	case *object.Function:
		return apply_function_from_method(function, method_args), true
	case *object.Builtin:
		return function.Fn(method_args...), true
	default:
		return object.NewError("'%s' is not a function", method_name), true
	}
}

// apply_function_from_method is a helper to apply user-defined functions as methods
// This is needed because we can't import from evaluator due to circular dependency
func apply_function_from_method(fn *object.Function, args []object.Object) object.Object {
	env := object.NewEnclosedEnvironment(fn.Env)

	for param_idx, param := range fn.Parameters {
		if param_idx < len(args) {
			env.Set(param.Name.Value, args[param_idx])
		}
	}

	// We need to evaluate the function body
	// But we can't call Eval directly due to circular import
	// So we'll use a workaround by storing a reference to the evaluator
	if eval_func != nil {
		evaluated := eval_func(fn.Body, env)
		if returnValue, ok := evaluated.(*object.ReturnValue); ok {
			return returnValue.Value
		}
		return evaluated
	}

	return object.NewError("internal error: evaluator not available")
}

// eval_func is set by the evaluator package to avoid circular imports
var eval_func func(node interface{}, env *object.Environment) object.Object

// Type registries hold user-defined methods
var array_registry *object.Map
var string_registry *object.Map
var number_registry *object.Map
var map_registry *object.Map

// SetEvaluator sets the eval function reference for method calls
func SetEvaluator(fn func(node interface{}, env *object.Environment) object.Object) {
	eval_func = fn
}

// SetArrayRegistry sets the global Array object registry
func SetArrayRegistry(obj *object.Map) {
	array_registry = obj
}

// SetStringRegistry sets the global String object registry
func SetStringRegistry(obj *object.Map) {
	string_registry = obj
}

// SetNumberRegistry sets the global Number object registry
func SetNumberRegistry(obj *object.Map) {
	number_registry = obj
}

// SetMapRegistry sets the global Map object registry
func SetMapRegistry(obj *object.Map) {
	map_registry = obj
}

// Number Methods

func call_number_method(num *object.Number, method_name string, args []object.Object) object.Object {
	switch method_name {
	case "to_string":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Number.to_string. got=%d, want=0", len(args))
		}
		return &object.String{Value: fmt.Sprintf("%.10g", num.Value)}

	case "abs":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Number.abs. got=%d, want=0", len(args))
		}
		return &object.Number{Value: math.Abs(num.Value)}

	case "floor":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Number.floor. got=%d, want=0", len(args))
		}
		return &object.Number{Value: math.Floor(num.Value)}

	case "ceil":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Number.ceil. got=%d, want=0", len(args))
		}
		return &object.Number{Value: math.Ceil(num.Value)}

	case "round":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Number.round. got=%d, want=0", len(args))
		}
		return &object.Number{Value: math.Round(num.Value)}

	case "sqrt":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Number.sqrt. got=%d, want=0", len(args))
		}
		if num.Value < 0 {
			return object.NewError("cannot take square root of negative number")
		}
		return &object.Number{Value: math.Sqrt(num.Value)}
	}

	// Check for instance-specific custom properties
	if result, found := check_custom_property(num.Properties, method_name, num, args); found {
		return result
	}

	// Check for user-defined methods in the global number registry
	if result, found := check_type_registry(number_registry, method_name, num, args); found {
		return result
	}

	return object.NewError("method '%s' not found on Number", method_name)
}

// Map Methods

func call_map_method(map_obj *object.Map, method_name string, args []object.Object) object.Object {
	// Check for instance-specific custom properties
	if result, found := check_custom_property(map_obj.Properties, method_name, map_obj, args); found {
		return result
	}

	// Check for user-defined methods in the global map registry
	if result, found := check_type_registry(map_registry, method_name, map_obj, args); found {
		return result
	}

	return object.NewError("method '%s' not found on Map", method_name)
}

// Boolean Methods

func call_boolean_method(bool *object.Boolean, method_name string, args []object.Object) object.Object {
	switch method_name {
	case "to_string":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Boolean.to_string. got=%d, want=0", len(args))
		}
		return &object.String{Value: fmt.Sprintf("%t", bool.Value)}
	}

	// Check for instance-specific custom properties
	if result, found := check_custom_property(bool.Properties, method_name, bool, args); found {
		return result
	}

	return object.NewError("method '%s' not found on Boolean", method_name)
}

// Error Methods

func call_error_method(err *object.Error, method_name string, args []object.Object) object.Object {
	switch method_name {
	case "to_string":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Error.to_string. got=%d, want=0", len(args))
		}
		// Return just the error message without "ERROR: " prefix
		return &object.String{Value: err.Message}
	}

	return object.NewError("method '%s' not found on Error", method_name)
}

// Global Functions (kept as builtin functions for print/println)

var global_functions = map[string]*object.Builtin{
	"print":   {Fn: print_builtin},
	"println": {Fn: println_builtin},
	"isNull":  {Fn: is_null_builtin},
	"error":   {Fn: error_builtin},
}

// get_global_function returns a global function by name
func get_global_function(name string) *object.Builtin {
	return global_functions[name]
}

func print_builtin(args ...object.Object) object.Object {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	return object.NULL
}

func println_builtin(args ...object.Object) object.Object {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	fmt.Println()
	return object.NULL
}

func is_null_builtin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0] == object.NULL || args[0].Type() == object.NULL_OBJ {
		return object.TRUE
	}
	return object.FALSE
}

func error_builtin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	// Convert argument to string for error message
	var message string
	switch arg := args[0].(type) {
	case *object.String:
		message = arg.Value
	default:
		message = arg.String()
	}

	// Create a user-created error
	err := object.NewError("%s", message)
	err.IsUserCreated = true
	return err
}