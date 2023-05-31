package util

import (
	"fmt"
	"reflect"
)

func GetFieldValue(iface any, fieldName string) any {
	val := reflect.ValueOf(iface)
	if val.Kind() == reflect.Map {
		v := val.MapIndex(reflect.ValueOf(fieldName))
		if v.Kind() == reflect.Invalid || v.IsNil() {
			return nil
		}
		return v.Interface()
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	fieldType := val.Type()
	for i := 0; i < fieldType.NumField(); i++ {
		field := fieldType.Field(i)
		if field.Name == fieldName {
			fieldVal := val.Field(i)
			return fieldVal.Interface()
		}
	}
	return nil // field not found
}

func SetInterfaceField(obj any, fieldName string, newValue any) error {
	// get the type and value of the input object
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)

	if objType.Kind() == reflect.Map {
		objValue.SetMapIndex(reflect.ValueOf(fieldName), reflect.ValueOf(newValue))
		return nil
	}

	// check if the object is a pointer
	if objType.Kind() != reflect.Ptr {
		return fmt.Errorf("input object must be a pointer")
	}

	// get the underlying value of the pointer and check if it is a struct
	s := objValue.Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("input object must point to a struct")
	}

	// get the field by name
	field := s.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}

	// check if the field is an interface
	if field.Type().Kind() != reflect.Interface {
		return fmt.Errorf("field %s is not an interface", fieldName)
	}

	// set the field to the new value
	field.Set(reflect.ValueOf(newValue))

	// return no error on successful assignment
	return nil
}
