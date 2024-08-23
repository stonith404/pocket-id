package dto

import (
	"errors"
	"reflect"
)

// MapStructList maps a list of source structs to a list of destination structs
func MapStructList[S any, D any](source []S, destination *[]D) error {
	for _, item := range source {
		var destItem D
		if err := MapStruct(item, &destItem); err != nil {
			return err
		}
		*destination = append(*destination, destItem)
	}
	return nil
}

// MapStruct maps a source struct to a destination struct
func MapStruct[S any, D any](source S, destination *D) error {
	// Ensure destination is a non-nil pointer
	destValue := reflect.ValueOf(destination)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return errors.New("destination must be a non-nil pointer to a struct")
	}

	// Ensure source is a struct
	sourceValue := reflect.ValueOf(source)
	if sourceValue.Kind() != reflect.Struct {
		return errors.New("source must be a struct")
	}

	return mapStructInternal(sourceValue, destValue.Elem())
}

func mapStructInternal(sourceVal reflect.Value, destVal reflect.Value) error {
	// Loop through the fields of the destination struct
	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Field(i)
		destFieldType := destVal.Type().Field(i)

		if destFieldType.Anonymous {
			// Recursively handle embedded structs
			if err := mapStructInternal(sourceVal, destField); err != nil {
				return err
			}
			continue
		}

		sourceField := sourceVal.FieldByName(destFieldType.Name)

		// If the source field is valid and can be assigned to the destination field
		if sourceField.IsValid() && destField.CanSet() {
			// Handle direct assignment for simple types
			if sourceField.Type() == destField.Type() {
				destField.Set(sourceField)
			} else if sourceField.Kind() == reflect.Slice && destField.Kind() == reflect.Slice {
				// Handle slices
				if sourceField.Type().Elem() == destField.Type().Elem() {
					newSlice := reflect.MakeSlice(destField.Type(), sourceField.Len(), sourceField.Cap())

					for j := 0; j < sourceField.Len(); j++ {
						newSlice.Index(j).Set(sourceField.Index(j))
					}

					destField.Set(newSlice)
				}
			} else if sourceField.Kind() == reflect.Struct && destField.Kind() == reflect.Struct {
				// Recursively map nested structs
				if err := mapStructInternal(sourceField, destField); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
