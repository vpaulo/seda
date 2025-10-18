package evaluator

import (
	"fmt"
	"math"
	"strings"
	"time"

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
	case *object.Time:
		return call_time_method(obj, method_name, args)
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

	case "split":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.split. got=%d, want=1", len(args))
		}

		delimiter, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to String.split must be STRING, got %s", args[0].Type())
		}

		parts := strings.Split(str.Value, delimiter.Value)
		elements := make([]object.Object, len(parts))
		for i, part := range parts {
			elements[i] = &object.String{Value: part}
		}
		return &object.Array{Elements: elements}

	case "trim":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.trim. got=%d, want=0", len(args))
		}
		return &object.String{Value: strings.TrimSpace(str.Value)}

	case "replace":
		if len(args) != 2 {
			return object.NewError("wrong number of arguments for String.replace. got=%d, want=2", len(args))
		}

		old, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("first argument to String.replace must be STRING, got %s", args[0].Type())
		}

		new, ok := args[1].(*object.String)
		if !ok {
			return object.NewError("second argument to String.replace must be STRING, got %s", args[1].Type())
		}

		return &object.String{Value: strings.ReplaceAll(str.Value, old.Value, new.Value)}

	case "starts_with":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.starts_with. got=%d, want=1", len(args))
		}

		prefix, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to String.starts_with must be STRING, got %s", args[0].Type())
		}

		if strings.HasPrefix(str.Value, prefix.Value) {
			return object.TRUE
		}
		return object.FALSE

	case "ends_with":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.ends_with. got=%d, want=1", len(args))
		}

		suffix, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to String.ends_with must be STRING, got %s", args[0].Type())
		}

		if strings.HasSuffix(str.Value, suffix.Value) {
			return object.TRUE
		}
		return object.FALSE

	case "index_of":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.index_of. got=%d, want=1", len(args))
		}

		substring, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to String.index_of must be STRING, got %s", args[0].Type())
		}

		index := strings.Index(str.Value, substring.Value)
		return &object.Number{Value: float64(index)}

	case "char_at":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.char_at. got=%d, want=1", len(args))
		}

		idx, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to String.char_at must be NUMBER, got %s", args[0].Type())
		}

		index := int(idx.Value)
		if index < 0 || index >= len(str.Value) {
			return object.NewError("index out of bounds")
		}

		return &object.String{Value: string(str.Value[index])}

	case "trim_left":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.trim_left. got=%d, want=0", len(args))
		}
		return &object.String{Value: strings.TrimLeft(str.Value, " \t\n\r")}

	case "trim_right":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.trim_right. got=%d, want=0", len(args))
		}
		return &object.String{Value: strings.TrimRight(str.Value, " \t\n\r")}

	case "contains":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.contains. got=%d, want=1", len(args))
		}

		substring, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to String.contains must be STRING, got %s", args[0].Type())
		}

		if strings.Contains(str.Value, substring.Value) {
			return object.TRUE
		}
		return object.FALSE

	case "last_index_of":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.last_index_of. got=%d, want=1", len(args))
		}

		substring, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to String.last_index_of must be STRING, got %s", args[0].Type())
		}

		index := strings.LastIndex(str.Value, substring.Value)
		return &object.Number{Value: float64(index)}

	case "count":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.count. got=%d, want=1", len(args))
		}

		substring, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to String.count must be STRING, got %s", args[0].Type())
		}

		count := strings.Count(str.Value, substring.Value)
		return &object.Number{Value: float64(count)}

	case "replace_first":
		if len(args) != 2 {
			return object.NewError("wrong number of arguments for String.replace_first. got=%d, want=2", len(args))
		}

		old, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("first argument to String.replace_first must be STRING, got %s", args[0].Type())
		}

		new, ok := args[1].(*object.String)
		if !ok {
			return object.NewError("second argument to String.replace_first must be STRING, got %s", args[1].Type())
		}

		return &object.String{Value: strings.Replace(str.Value, old.Value, new.Value, 1)}

	case "reverse":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.reverse. got=%d, want=0", len(args))
		}

		runes := []rune(str.Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return &object.String{Value: string(runes)}

	case "repeat":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for String.repeat. got=%d, want=1", len(args))
		}

		count, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to String.repeat must be NUMBER, got %s", args[0].Type())
		}

		return &object.String{Value: strings.Repeat(str.Value, int(count.Value))}

	case "lines":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.lines. got=%d, want=0", len(args))
		}

		lines := strings.Split(str.Value, "\n")
		elements := make([]object.Object, len(lines))
		for i, line := range lines {
			elements[i] = &object.String{Value: line}
		}
		return &object.Array{Elements: elements}

	case "chars":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.chars. got=%d, want=0", len(args))
		}

		runes := []rune(str.Value)
		elements := make([]object.Object, len(runes))
		for i, r := range runes {
			elements[i] = &object.String{Value: string(r)}
		}
		return &object.Array{Elements: elements}

	case "words":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.words. got=%d, want=0", len(args))
		}

		words := strings.Fields(str.Value) // Splits by whitespace and trims
		elements := make([]object.Object, len(words))
		for i, word := range words {
			elements[i] = &object.String{Value: word}
		}
		return &object.Array{Elements: elements}

	case "capitalize":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.capitalize. got=%d, want=0", len(args))
		}

		if len(str.Value) == 0 {
			return str
		}

		runes := []rune(str.Value)
		runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
		return &object.String{Value: string(runes)}

	case "title_case":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.title_case. got=%d, want=0", len(args))
		}

		return &object.String{Value: strings.Title(str.Value)}

	case "is_empty":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.is_empty. got=%d, want=0", len(args))
		}

		if len(str.Value) == 0 {
			return object.TRUE
		}
		return object.FALSE

	case "is_blank":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.is_blank. got=%d, want=0", len(args))
		}

		if len(strings.TrimSpace(str.Value)) == 0 {
			return object.TRUE
		}
		return object.FALSE

	case "is_numeric":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.is_numeric. got=%d, want=0", len(args))
		}

		for _, r := range str.Value {
			if !strings.ContainsRune("0123456789", r) {
				return object.FALSE
			}
		}
		if len(str.Value) == 0 {
			return object.FALSE
		}
		return object.TRUE

	case "is_alpha":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for String.is_alpha. got=%d, want=0", len(args))
		}

		for _, r := range str.Value {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
				return object.FALSE
			}
		}
		if len(str.Value) == 0 {
			return object.FALSE
		}
		return object.TRUE

	case "pad_left":
		if len(args) != 2 {
			return object.NewError("wrong number of arguments for String.pad_left. got=%d, want=2", len(args))
		}

		width, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("first argument to String.pad_left must be NUMBER, got %s", args[0].Type())
		}

		pad_str, ok := args[1].(*object.String)
		if !ok {
			return object.NewError("second argument to String.pad_left must be STRING, got %s", args[1].Type())
		}

		target_width := int(width.Value)
		current_len := len(str.Value)
		if current_len >= target_width {
			return str
		}

		pad_len := target_width - current_len
		if len(pad_str.Value) == 0 {
			return str
		}

		// Repeat padding string enough times
		padding := strings.Repeat(pad_str.Value, (pad_len/len(pad_str.Value))+1)
		padding = padding[:pad_len] // Trim to exact length needed

		return &object.String{Value: padding + str.Value}

	case "pad_right":
		if len(args) != 2 {
			return object.NewError("wrong number of arguments for String.pad_right. got=%d, want=2", len(args))
		}

		width, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("first argument to String.pad_right must be NUMBER, got %s", args[0].Type())
		}

		pad_str, ok := args[1].(*object.String)
		if !ok {
			return object.NewError("second argument to String.pad_right must be STRING, got %s", args[1].Type())
		}

		target_width := int(width.Value)
		current_len := len(str.Value)
		if current_len >= target_width {
			return str
		}

		pad_len := target_width - current_len
		if len(pad_str.Value) == 0 {
			return str
		}

		// Repeat padding string enough times
		padding := strings.Repeat(pad_str.Value, (pad_len/len(pad_str.Value))+1)
		padding = padding[:pad_len] // Trim to exact length needed

		return &object.String{Value: str.Value + padding}
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
		// Check immutability
		if arr.IsImmutable {
			return object.NewError("cannot call push() on immutable array")
		}

		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.push. got=%d, want=1", len(args))
		}

		// Mutate the array in place
		arr.Elements = append(arr.Elements, args[0])
		return arr

	case "pop":
		// Check immutability
		if arr.IsImmutable {
			return object.NewError("cannot call pop() on immutable array")
		}

		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.pop. got=%d, want=0", len(args))
		}

		length := len(arr.Elements)
		if length == 0 {
			return object.NULL
		}

		// Get the last element
		lastElement := arr.Elements[length-1]

		// Remove the last element from the array
		arr.Elements = arr.Elements[:length-1]

		return lastElement

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

	// Functional operations
	case "map":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.map. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.map must be FUNCTION, got %s", args[0].Type())
		}

		result := make([]object.Object, len(arr.Elements))
		for i, elem := range arr.Elements {
			mapped := apply_function_from_method(fn, []object.Object{elem})
			if is_error(mapped) {
				return mapped
			}
			result[i] = mapped
		}
		return &object.Array{Elements: result}

	case "filter":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.filter. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.filter must be FUNCTION, got %s", args[0].Type())
		}

		result := []object.Object{}
		for _, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if is_truthy(condition) {
				result = append(result, elem)
			}
		}
		return &object.Array{Elements: result}

	case "reduce":
		if len(args) != 2 {
			return object.NewError("wrong number of arguments for Array.reduce. got=%d, want=2", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("first argument to Array.reduce must be FUNCTION, got %s", args[0].Type())
		}

		accumulator := args[1]
		for _, elem := range arr.Elements {
			accumulator = apply_function_from_method(fn, []object.Object{accumulator, elem})
			if is_error(accumulator) {
				return accumulator
			}
		}
		return accumulator

	case "each":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.each. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.each must be FUNCTION, got %s", args[0].Type())
		}

		for _, elem := range arr.Elements {
			result := apply_function_from_method(fn, []object.Object{elem})
			if is_error(result) {
				return result
			}
		}
		return object.NULL

	case "map_with_index":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.map_with_index. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.map_with_index must be FUNCTION, got %s", args[0].Type())
		}

		result := make([]object.Object, len(arr.Elements))
		for i, elem := range arr.Elements {
			mapped := apply_function_from_method(fn, []object.Object{elem, &object.Number{Value: float64(i)}})
			if is_error(mapped) {
				return mapped
			}
			result[i] = mapped
		}
		return &object.Array{Elements: result}

	// Finding methods
	case "find":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.find. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.find must be FUNCTION, got %s", args[0].Type())
		}

		for _, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if is_truthy(condition) {
				return elem
			}
		}
		return object.NULL

	case "find_index":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.find_index. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.find_index must be FUNCTION, got %s", args[0].Type())
		}

		for i, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if is_truthy(condition) {
				return &object.Number{Value: float64(i)}
			}
		}
		return &object.Number{Value: -1}

	case "any":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.any. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.any must be FUNCTION, got %s", args[0].Type())
		}

		for _, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if is_truthy(condition) {
				return object.TRUE
			}
		}
		return object.FALSE

	case "all":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.all. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.all must be FUNCTION, got %s", args[0].Type())
		}

		for _, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if !is_truthy(condition) {
				return object.FALSE
			}
		}
		return object.TRUE

	case "none":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.none. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.none must be FUNCTION, got %s", args[0].Type())
		}

		for _, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if is_truthy(condition) {
				return object.FALSE
			}
		}
		return object.TRUE

	case "count":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.count. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.count must be FUNCTION, got %s", args[0].Type())
		}

		count := 0
		for _, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if is_truthy(condition) {
				count++
			}
		}
		return &object.Number{Value: float64(count)}

	// Transformation methods
	case "sort":
		// Check immutability
		if arr.IsImmutable {
			return object.NewError("cannot call sort() on immutable array")
		}

		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.sort. got=%d, want=0", len(args))
		}

		// Sort mutates the array in place
		// Simple number sorting for now
		sorted := make([]object.Object, len(arr.Elements))
		copy(sorted, arr.Elements)

		// Bubble sort for simplicity
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				num1, ok1 := sorted[i].(*object.Number)
				num2, ok2 := sorted[j].(*object.Number)
				if ok1 && ok2 && num1.Value > num2.Value {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}

		arr.Elements = sorted
		return arr

	case "sort_by":
		// Check immutability
		if arr.IsImmutable {
			return object.NewError("cannot call sort_by() on immutable array")
		}

		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.sort_by. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.sort_by must be FUNCTION, got %s", args[0].Type())
		}

		// Create copy with sort keys
		type sortPair struct {
			elem object.Object
			key  float64
		}

		pairs := make([]sortPair, len(arr.Elements))
		for i, elem := range arr.Elements {
			keyObj := apply_function_from_method(fn, []object.Object{elem})
			if is_error(keyObj) {
				return keyObj
			}
			keyNum, ok := keyObj.(*object.Number)
			if !ok {
				return object.NewError("sort_by function must return NUMBER, got %s", keyObj.Type())
			}
			pairs[i] = sortPair{elem: elem, key: keyNum.Value}
		}

		// Bubble sort
		for i := 0; i < len(pairs); i++ {
			for j := i + 1; j < len(pairs); j++ {
				if pairs[i].key > pairs[j].key {
					pairs[i], pairs[j] = pairs[j], pairs[i]
				}
			}
		}

		// Extract sorted elements
		sorted := make([]object.Object, len(pairs))
		for i, p := range pairs {
			sorted[i] = p.elem
		}

		arr.Elements = sorted
		return arr

	case "reverse":
		// Check immutability
		if arr.IsImmutable {
			return object.NewError("cannot call reverse() on immutable array")
		}

		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.reverse. got=%d, want=0", len(args))
		}

		// Reverse mutates in place
		for i, j := 0, len(arr.Elements)-1; i < j; i, j = i+1, j-1 {
			arr.Elements[i], arr.Elements[j] = arr.Elements[j], arr.Elements[i]
		}
		return arr

	case "unique":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.unique. got=%d, want=0", len(args))
		}

		seen := make(map[string]bool)
		result := []object.Object{}

		for _, elem := range arr.Elements {
			key := elem.String()
			if !seen[key] {
				seen[key] = true
				result = append(result, elem)
			}
		}

		return &object.Array{Elements: result}

	// Slicing/Combining
	case "slice":
		if len(args) != 2 {
			return object.NewError("wrong number of arguments for Array.slice. got=%d, want=2", len(args))
		}

		start, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("first argument to Array.slice must be NUMBER, got %s", args[0].Type())
		}

		end, ok := args[1].(*object.Number)
		if !ok {
			return object.NewError("second argument to Array.slice must be NUMBER, got %s", args[1].Type())
		}

		startIdx := int(start.Value)
		endIdx := int(end.Value)

		if startIdx < 0 || startIdx >= len(arr.Elements) {
			startIdx = 0
		}
		if endIdx < 0 || endIdx > len(arr.Elements) {
			endIdx = len(arr.Elements)
		}
		if endIdx < startIdx {
			endIdx = startIdx
		}

		result := make([]object.Object, endIdx-startIdx)
		copy(result, arr.Elements[startIdx:endIdx])
		return &object.Array{Elements: result}

	case "take":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.take. got=%d, want=1", len(args))
		}

		n, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to Array.take must be NUMBER, got %s", args[0].Type())
		}

		count := int(n.Value)
		if count < 0 {
			count = 0
		}
		if count > len(arr.Elements) {
			count = len(arr.Elements)
		}

		result := make([]object.Object, count)
		copy(result, arr.Elements[:count])
		return &object.Array{Elements: result}

	case "drop":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.drop. got=%d, want=1", len(args))
		}

		n, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to Array.drop must be NUMBER, got %s", args[0].Type())
		}

		count := int(n.Value)
		if count < 0 {
			count = 0
		}
		if count > len(arr.Elements) {
			count = len(arr.Elements)
		}

		result := make([]object.Object, len(arr.Elements)-count)
		copy(result, arr.Elements[count:])
		return &object.Array{Elements: result}

	case "concat":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.concat. got=%d, want=1", len(args))
		}

		other, ok := args[0].(*object.Array)
		if !ok {
			return object.NewError("argument to Array.concat must be ARRAY, got %s", args[0].Type())
		}

		result := make([]object.Object, len(arr.Elements)+len(other.Elements))
		copy(result, arr.Elements)
		copy(result[len(arr.Elements):], other.Elements)
		return &object.Array{Elements: result}

	// Nested arrays
	case "flatten":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.flatten. got=%d, want=0", len(args))
		}

		var flatten func([]object.Object) []object.Object
		flatten = func(elems []object.Object) []object.Object {
			result := []object.Object{}
			for _, elem := range elems {
				if subArr, ok := elem.(*object.Array); ok {
					result = append(result, flatten(subArr.Elements)...)
				} else {
					result = append(result, elem)
				}
			}
			return result
		}

		return &object.Array{Elements: flatten(arr.Elements)}

	case "flat_map":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.flat_map. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.flat_map must be FUNCTION, got %s", args[0].Type())
		}

		result := []object.Object{}
		for _, elem := range arr.Elements {
			mapped := apply_function_from_method(fn, []object.Object{elem})
			if is_error(mapped) {
				return mapped
			}
			if subArr, ok := mapped.(*object.Array); ok {
				result = append(result, subArr.Elements...)
			} else {
				result = append(result, mapped)
			}
		}
		return &object.Array{Elements: result}

	// Membership
	case "contains":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.contains. got=%d, want=1", len(args))
		}

		searchKey := args[0].String()
		for _, elem := range arr.Elements {
			if elem.String() == searchKey {
				return object.TRUE
			}
		}
		return object.FALSE

	case "index_of":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.index_of. got=%d, want=1", len(args))
		}

		searchKey := args[0].String()
		for i, elem := range arr.Elements {
			if elem.String() == searchKey {
				return &object.Number{Value: float64(i)}
			}
		}
		return &object.Number{Value: -1}

	case "last_index_of":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.last_index_of. got=%d, want=1", len(args))
		}

		searchKey := args[0].String()
		lastIdx := -1
		for i, elem := range arr.Elements {
			if elem.String() == searchKey {
				lastIdx = i
			}
		}
		return &object.Number{Value: float64(lastIdx)}

	// Conversion
	case "join":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.join. got=%d, want=1", len(args))
		}

		separator, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to Array.join must be STRING, got %s", args[0].Type())
		}

		parts := make([]string, len(arr.Elements))
		for i, elem := range arr.Elements {
			parts[i] = elem.String()
		}
		return &object.String{Value: strings.Join(parts, separator.Value)}

	// Statistics
	case "sum":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.sum. got=%d, want=0", len(args))
		}

		sum := 0.0
		for _, elem := range arr.Elements {
			num, ok := elem.(*object.Number)
			if !ok {
				return object.NewError("Array.sum requires all elements to be NUMBER, got %s", elem.Type())
			}
			sum += num.Value
		}
		return &object.Number{Value: sum}

	case "average":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.average. got=%d, want=0", len(args))
		}

		if len(arr.Elements) == 0 {
			return object.NewError("cannot compute average of empty array")
		}

		sum := 0.0
		for _, elem := range arr.Elements {
			num, ok := elem.(*object.Number)
			if !ok {
				return object.NewError("Array.average requires all elements to be NUMBER, got %s", elem.Type())
			}
			sum += num.Value
		}
		return &object.Number{Value: sum / float64(len(arr.Elements))}

	case "min":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.min. got=%d, want=0", len(args))
		}

		if len(arr.Elements) == 0 {
			return object.NewError("cannot find min of empty array")
		}

		minNum, ok := arr.Elements[0].(*object.Number)
		if !ok {
			return object.NewError("Array.min requires all elements to be NUMBER")
		}

		minVal := minNum.Value
		for _, elem := range arr.Elements[1:] {
			num, ok := elem.(*object.Number)
			if !ok {
				return object.NewError("Array.min requires all elements to be NUMBER, got %s", elem.Type())
			}
			if num.Value < minVal {
				minVal = num.Value
			}
		}
		return &object.Number{Value: minVal}

	case "max":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.max. got=%d, want=0", len(args))
		}

		if len(arr.Elements) == 0 {
			return object.NewError("cannot find max of empty array")
		}

		maxNum, ok := arr.Elements[0].(*object.Number)
		if !ok {
			return object.NewError("Array.max requires all elements to be NUMBER")
		}

		maxVal := maxNum.Value
		for _, elem := range arr.Elements[1:] {
			num, ok := elem.(*object.Number)
			if !ok {
				return object.NewError("Array.max requires all elements to be NUMBER, got %s", elem.Type())
			}
			if num.Value > maxVal {
				maxVal = num.Value
			}
		}
		return &object.Number{Value: maxVal}

	// Grouping
	case "chunk":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.chunk. got=%d, want=1", len(args))
		}

		size, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to Array.chunk must be NUMBER, got %s", args[0].Type())
		}

		chunkSize := int(size.Value)
		if chunkSize <= 0 {
			return object.NewError("chunk size must be positive")
		}

		result := []object.Object{}
		for i := 0; i < len(arr.Elements); i += chunkSize {
			end := i + chunkSize
			if end > len(arr.Elements) {
				end = len(arr.Elements)
			}
			chunk := make([]object.Object, end-i)
			copy(chunk, arr.Elements[i:end])
			result = append(result, &object.Array{Elements: chunk})
		}

		return &object.Array{Elements: result}

	case "partition":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.partition. got=%d, want=1", len(args))
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return object.NewError("argument to Array.partition must be FUNCTION, got %s", args[0].Type())
		}

		trueGroup := []object.Object{}
		falseGroup := []object.Object{}

		for _, elem := range arr.Elements {
			condition := apply_function_from_method(fn, []object.Object{elem})
			if is_error(condition) {
				return condition
			}
			if is_truthy(condition) {
				trueGroup = append(trueGroup, elem)
			} else {
				falseGroup = append(falseGroup, elem)
			}
		}

		result := []object.Object{
			&object.Array{Elements: trueGroup},
			&object.Array{Elements: falseGroup},
		}
		return &object.Array{Elements: result}

	// Array operations
	case "zip":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Array.zip. got=%d, want=1", len(args))
		}

		other, ok := args[0].(*object.Array)
		if !ok {
			return object.NewError("argument to Array.zip must be ARRAY, got %s", args[0].Type())
		}

		length := len(arr.Elements)
		if len(other.Elements) < length {
			length = len(other.Elements)
		}

		result := make([]object.Object, length)
		for i := 0; i < length; i++ {
			pair := []object.Object{arr.Elements[i], other.Elements[i]}
			result[i] = &object.Array{Elements: pair}
		}

		return &object.Array{Elements: result}

	case "compact":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Array.compact. got=%d, want=0", len(args))
		}

		result := []object.Object{}
		for _, elem := range arr.Elements {
			if elem.Type() != object.NULL_OBJ {
				result = append(result, elem)
			}
		}

		return &object.Array{Elements: result}
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

// Time Methods

func call_time_method(t *object.Time, method_name string, args []object.Object) object.Object {
	switch method_name {
	// Formatting methods
	case "format":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.format. got=%d, want=1", len(args))
		}
		formatStr, ok := args[0].(*object.String)
		if !ok {
			return object.NewError("argument to Time.format must be STRING, got %s", args[0].Type())
		}
		// Convert Seda format to Go format
		goFormat := convertSedaFormatToGo(formatStr.Value)
		return &object.String{Value: t.Value.Format(goFormat)}

	case "to_string":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.to_string. got=%d, want=0", len(args))
		}
		return &object.String{Value: t.Value.Format(time.RFC3339)}

	case "unix":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.unix. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Unix())}

	case "unix_millis":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.unix_millis. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.UnixMilli())}

	// Component methods
	case "year":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.year. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Year())}

	case "month":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.month. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Month())}

	case "day":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.day. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Day())}

	case "hour":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.hour. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Hour())}

	case "minute":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.minute. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Minute())}

	case "second":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.second. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Second())}

	case "weekday":
		if len(args) != 0 {
			return object.NewError("wrong number of arguments for Time.weekday. got=%d, want=0", len(args))
		}
		return &object.Number{Value: float64(t.Value.Weekday())}

	// Arithmetic methods (return new Time)
	case "add_seconds":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.add_seconds. got=%d, want=1", len(args))
		}
		seconds, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to Time.add_seconds must be NUMBER, got %s", args[0].Type())
		}
		duration := time.Duration(seconds.Value) * time.Second
		return &object.Time{Value: t.Value.Add(duration)}

	case "add_minutes":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.add_minutes. got=%d, want=1", len(args))
		}
		minutes, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to Time.add_minutes must be NUMBER, got %s", args[0].Type())
		}
		duration := time.Duration(minutes.Value) * time.Minute
		return &object.Time{Value: t.Value.Add(duration)}

	case "add_hours":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.add_hours. got=%d, want=1", len(args))
		}
		hours, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to Time.add_hours must be NUMBER, got %s", args[0].Type())
		}
		duration := time.Duration(hours.Value) * time.Hour
		return &object.Time{Value: t.Value.Add(duration)}

	case "add_days":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.add_days. got=%d, want=1", len(args))
		}
		days, ok := args[0].(*object.Number)
		if !ok {
			return object.NewError("argument to Time.add_days must be NUMBER, got %s", args[0].Type())
		}
		duration := time.Duration(days.Value*24) * time.Hour
		return &object.Time{Value: t.Value.Add(duration)}

	// Comparison methods
	case "diff":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.diff. got=%d, want=1", len(args))
		}
		other, ok := args[0].(*object.Time)
		if !ok {
			return object.NewError("argument to Time.diff must be TIME, got %s", args[0].Type())
		}
		diff := t.Value.Sub(other.Value)
		return &object.Number{Value: diff.Seconds()}

	case "is_before":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.is_before. got=%d, want=1", len(args))
		}
		other, ok := args[0].(*object.Time)
		if !ok {
			return object.NewError("argument to Time.is_before must be TIME, got %s", args[0].Type())
		}
		if t.Value.Before(other.Value) {
			return object.TRUE
		}
		return object.FALSE

	case "is_after":
		if len(args) != 1 {
			return object.NewError("wrong number of arguments for Time.is_after. got=%d, want=1", len(args))
		}
		other, ok := args[0].(*object.Time)
		if !ok {
			return object.NewError("argument to Time.is_after must be TIME, got %s", args[0].Type())
		}
		if t.Value.After(other.Value) {
			return object.TRUE
		}
		return object.FALSE
	}

	return object.NewError("method '%s' not found on Time", method_name)
}

// convertSedaFormatToGo converts Seda time format strings to Go time format strings
func convertSedaFormatToGo(sedaFormat string) string {
	// Simple conversion for common patterns
	goFormat := sedaFormat
	goFormat = strings.ReplaceAll(goFormat, "YYYY", "2006")
	goFormat = strings.ReplaceAll(goFormat, "MM", "01")
	goFormat = strings.ReplaceAll(goFormat, "DD", "02")
	goFormat = strings.ReplaceAll(goFormat, "HH", "15")
	goFormat = strings.ReplaceAll(goFormat, "mm", "04")
	goFormat = strings.ReplaceAll(goFormat, "ss", "05")
	return goFormat
}