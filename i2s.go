package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		dataType := reflect.TypeOf(value)
		switch dataType.Kind() {
		case reflect.Map, reflect.Slice:
			out_new := ptr(structFieldValue).Interface()
			err := i2s(val.Interface(), out_new)
			if err != nil {
				return err
			}
			structFieldValue.Set(reflect.ValueOf(out_new).Elem())
		default:
			var convertedValue interface{}
			var err error
			switch structFieldType.Kind() {
			case reflect.Bool:
				convertedValue, err = strconv.ParseBool(value.(string))
			case reflect.Int, reflect.Int32, reflect.Int64:
				switch reflect.TypeOf(value).Kind() {
				case reflect.Float32, reflect.Float64:
					convertedValue, err = strconv.Atoi(fmt.Sprintf("%.0f", value))
				case reflect.String:
					return fmt.Errorf("Expected int")
				}
			case reflect.String:
				switch reflect.TypeOf(value).Kind() {
				case reflect.String:
					convertedValue = value.(string)
				default:
					return fmt.Errorf("Expected string")
				}
			default:
				return fmt.Errorf("Unknown format")
			}
			if err != nil {
				return err
			}
			structFieldValue.Set(reflect.ValueOf(convertedValue))
		}
		return nil
	}

	structFieldValue.Set(val)
	return nil
}

// ptr wraps the given value with pointer: V => *V, *V => **V, etc.
func ptr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.
	return pv
}

func i2s(data interface{}, out interface{}) error {
	dataType := reflect.TypeOf(data)
	if dataType == nil {
		return fmt.Errorf("Unknown type")
	}
	switch dataType.Kind() {
	case reflect.Slice:
		val := reflect.ValueOf(data)
		valuePtr := reflect.ValueOf(out)
		value := valuePtr.Elem()
		elementType := reflect.TypeOf(out).Elem()
		if elementType.Kind() != reflect.Slice {
			return fmt.Errorf("Slice is not expected")
		}
		value.Set(reflect.MakeSlice(elementType, val.Len(), val.Len()))
		for i := 0; i < val.Len(); i++ {
			out_new := ptr(value.Index(i)).Interface()
			err := i2s(val.Index(i).Elem().Interface(), out_new)
			if err != nil {
				return err
			}
			fmt.Printf("After set %v\n", out_new)
			value.Index(i).Set(reflect.ValueOf(out_new).Elem())
		}
	case reflect.Map:
		elementType := reflect.TypeOf(out)
		if elementType.Kind() != reflect.Ptr {
			return fmt.Errorf("No pointer type expected")
		}
		source := make(map[string]interface{})
		iter := reflect.ValueOf(data).MapRange()
		for iter.Next() {
			key := iter.Key().String()
			value := iter.Value().Interface()
			source[key] = value
			if value != nil {
				err := SetField(out, key, value)
				if err != nil {
					return err
				}
			}
		}
		if len(source) == 0 {
			return fmt.Errorf("Empty data")
		}
	case reflect.Struct:
		val := reflect.ValueOf(data)
		for i := 0; i < val.NumField(); i++ {
			err := SetField(out, val.Type().Field(i).Name, val.Field(i).Elem())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
