package main

import (
	//"errors"
	"fmt"
	"reflect"
	//"strconv"
)

func SetField(obj interface{}, name string, value interface{}) error {
	//var err error
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
		//switch structFieldType.Kind() {
		//	case reflect.Int:
		//		val, err = strconv.Atoi(val.String())
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}
		//return errors.New("Provided value type didn't match obj field type")
		structFieldValue.Set(val.Convert(structFieldType))
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
	switch dataType.Kind() {
	case reflect.Slice:
		val := reflect.ValueOf(data)
		valuePtr := reflect.ValueOf(out)
		value := valuePtr.Elem()
		elementType := reflect.TypeOf(out).Elem()
		value.Set(reflect.MakeSlice(elementType, val.Len(), val.Len()))
		for i := 0; i < val.Len(); i++ {
			//item := reflect.New(val.Index(i).Type()).Elem()
			//item := reflect.New(elementType)//.Elem()
			//err := i2s(val.Index(i).Elem().Interface(), item.Interface())
			out_new := ptr(value.Index(i)).Interface()
			err := i2s(val.Index(i).Elem().Interface(), out_new)
			if err != nil {
				return err
			}
			fmt.Printf("After set %v\n", out_new)
			//result = append(result, &item)
			value.Index(i).Set(reflect.ValueOf(out_new).Elem())
		}
		//result := reflect.MakeSlice(dataType, val.Len(), val.Len())
		//for i := 0; i < val.Len(); i++ {
		//	item := reflect.New(val.Index(i).Type())
		//	err := i2s(val.Index(i), item)
		//	if err != nil {
		//		return err
		//	}
		//	result = append(result, &item)
		//}
		//out = result.Interface()
		//reflect.ValueOf(out).Set(result)
	case reflect.Map:
		source := make(map[string]interface{})
		iter := reflect.ValueOf(data).MapRange()
		for iter.Next() {
			key := iter.Key().String()
			value := iter.Value().Interface()
			source[key] = value
			err := SetField(out, key, value)
			if err != nil {
				return err
			}
			//fmt.Print(key)
			//fmt.Print(value)
		}
		//return nil
	case reflect.Struct:
		val := reflect.ValueOf(data)
		for i := 0; i < val.NumField(); i++ {
			err := SetField(out, val.Type().Field(i).Name, val.Field(i).Elem())
			if err != nil {
				return err
			}
		}
	}
	//fmt.Println(source)
	//val := reflect.ValueOf(out)
	//for i := 0; i < val.NumField(); i++ {
	//	switch v := val.Field(i).(type) {
	//	case int64:
	//		val.Field(i).SetInt(strconv.Atoi(s[val.Field(i).Name]))
	//	}
	//}
	//s := reflect.TypeOf(out)//reflect.ValueOf(&out).Elem()
	////typeOfT := s.Type()
	//for i := 0; i < s.NumField(); i++ {
	//	f := s.Field(i)
	//	fmt.Printf("%v\n", f)
	//	//fmt.Printf("%d: %s %s = %v\n", i,
	//	//	typeOfT.Field(i).Name, f.Type(), f.Interface())
	//}
	//e := reflect.ValueOf(data)
	//
	//for i := 0; i < e.NumField(); i++ {
	//	varName := e.Type().Field(i).Name
	//	varType := e.Type().Field(i).Type
	//	varValue := e.Field(i).Interface()
	//	fmt.Printf("%v %v %v\n", varName,varType,varValue)
	//}
	//t := reflect.TypeOf(data).Elem()
	//typeof := reflect.TypeOf(data).Elem()
	//fmt.Println(typeof)
	//out = reflect.New(typeof).Elem()//.Interface()
	//fmt.Printf("Type of result %s\n", reflect.TypeOf(out).Name())

	//st := reflect.TypeOf(out)
	////
	//for i := 0; i < st.NumField(); i++ {
	//	field := st.Field(i)
	////
	////	json := field.Tag.Get("json")
	////	name := field.Name
	////
	////	fmt.Printf("json [%s] field [%s]\n", json, name)
	////
	//	val, _ := t.FieldByName(field.Name)
	//	fmt.Printf("field name %v\n", val)
	////	//val.Set(reflect.ValueOf(data.(map[string]interface{})[json]))
	////
	//}
	//out = new(interface{})
	return nil
}
