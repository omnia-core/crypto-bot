package log

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	maxLength = 4096
	log       = "log"
	required  = "required"
)

func logFields(v interface{}) string {
	val := reflect.ValueOf(v)
	typeName := val.Type().Name()

	// If v is a slice, process each item
	if val.Kind() == reflect.Slice {
		sliceLen := val.Len()
		sliceResults := make([]string, 0, sliceLen)
		for i := 0; i < sliceLen; i++ {
			sliceItem := val.Index(i)
			sliceResults = append(sliceResults, logFields(sliceItem.Interface()))
		}

		result := typeName + "[" + strings.Join(sliceResults, ", ") + "]"
		if len(result) > maxLength {
			overLen := len(result) - maxLength + 3 // +3 for "..."
			result = result[:maxLength-3] + fmt.Sprintf("...+%d", overLen)
		}
		return result
	}

	// If v is a struct, process its fields
	if val.Kind() == reflect.Struct {
		result := typeName + "{" + logFieldsRecursive(v) + "}"
		if len(result) > maxLength {
			overLen := len(result) - maxLength + 3 // +3 for "..."
			result = result[:maxLength-3] + fmt.Sprintf("...+%d}", overLen)
		}
		return result
	}

	// Otherwise, just convert v to a string
	result := fmt.Sprintf("%v", v)
	if len(result) > maxLength {
		overLen := len(result) - maxLength + 3 // +3 for "..."
		result = result[:maxLength-3] + fmt.Sprintf("...+%d", overLen)
	}
	return result
}

func logFieldsRecursive(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()

	if val.Kind() != reflect.Struct {
		return ""
	}

	results := make([]string, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get(log)

		// Check if field is a struct and recursively log its fields
		if fieldValue.Kind() == reflect.Struct && tag == required {
			innerLog := logFieldsRecursive(fieldValue.Interface())
			if innerLog != "" {
				results = append(results, fieldType.Name+":{"+innerLog+"}")
			}
			continue
		}

		// Handle slices and arrays of structs
		if (fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array) &&
			fieldValue.Type().Elem().Kind() == reflect.Struct && tag == required {
			sliceResults := make([]string, 0, fieldValue.Len())
			for j := 0; j < fieldValue.Len(); j++ {
				sliceItemLog := logFieldsRecursive(fieldValue.Index(j).Interface())
				if sliceItemLog != "" {
					sliceResults = append(sliceResults, "{"+sliceItemLog+"}")
				}
			}
			results = append(results, fieldType.Name+":["+strings.Join(sliceResults, ", ")+"]")
			continue
		}

		if tag == required {
			results = append(results, fmt.Sprintf("%s:%v", fieldType.Name, fieldValue.Interface()))
		}
	}

	return strings.Join(results, ", ")
}
