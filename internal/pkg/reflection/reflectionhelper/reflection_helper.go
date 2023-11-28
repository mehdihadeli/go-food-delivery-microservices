package reflectionHelper

// ref: https://gist.github.com/drewolson/4771479
// https://stackoverflow.com/a/60598827/581476
// https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

func GetAllFields(typ reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i))
	}
	return fields
}

func GetFieldValueByIndex[T any](object T, index int) interface{} {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.Field(index)
		if field.IsValid() == false {
			return nil
		}
		// for all exported fields (public)
		if field.CanInterface() {
			return field.Interface()
		} else {
			// for all unexported fields (private)
			return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
		}
	} else if v.Kind() == reflect.Struct {
		// for all exported fields (public)
		val := v
		field := val.Field(index)
		if field.IsValid() == false {
			return nil
		}
		if field.CanInterface() {
			return field.Interface()
		} else {
			// for all unexported fields (private)
			rs2 := reflect.New(val.Type()).Elem()
			rs2.Set(val)
			val = rs2.Field(index)
			val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()

			return val.Interface()
		}
	}
	return nil
}

func GetFieldValueByName[T any](object T, name string) interface{} {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.FieldByName(name)
		if field.IsValid() == false {
			return nil
		}
		// for all exported fields (public)
		if field.CanInterface() {
			return field.Interface()
		} else {
			// for all unexported fields (private)
			return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
		}
	} else if v.Kind() == reflect.Struct {
		// for all exported fields (public)
		val := v
		field := val.FieldByName(name)
		if field.IsValid() == false {
			return nil
		}
		if field.CanInterface() {
			return field.Interface()
		} else {
			// for all unexported fields (private)
			rs2 := reflect.New(val.Type()).Elem()
			rs2.Set(val)
			val = rs2.FieldByName(name)
			val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()

			return val.Interface()
		}
	}
	return nil
}

func SetFieldValueByIndex[T any](object T, index int, value interface{}) {
	v := reflect.ValueOf(&object).Elem()

	// https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.Field(index)
		if field.IsValid() == false {
			return
		}
		// for all exported fields (public)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			// for all unexported fields (private)
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
		}
	} else if v.Kind() == reflect.Struct {
		// for all exported fields (public)
		val := v
		field := val.Field(index)
		if field.IsValid() == false {
			return
		}
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
			object = val.Interface().(T)
		} else {
			// for all unexported fields (private)
			rs2 := reflect.New(val.Type()).Elem()
			rs2.Set(val)
			val = rs2.Field(index)
			val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()

			val.Set(reflect.ValueOf(value))
		}
	}
}

func SetFieldValueByName[T any](object T, name string, value interface{}) {
	v := reflect.ValueOf(&object).Elem()

	// https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
	if v.Kind() == reflect.Ptr {
		val := v.Elem()
		field := val.FieldByName(name)
		if field.IsValid() == false {
			return
		}
		// for all exported fields (public)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			// for all unexported fields (private)
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
		}
	} else if v.Kind() == reflect.Struct {
		// for all exported fields (public)
		val := v
		field := val.FieldByName(name)
		if field.IsValid() == false {
			return
		}
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
			object = val.Interface().(T)
		} else {
			// for all unexported fields (private)
			rs2 := reflect.New(val.Type()).Elem()
			rs2.Set(val)
			val = rs2.FieldByName(name)
			val = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()

			val.Set(reflect.ValueOf(value))
		}
	}
}

func GetFieldValue(field reflect.Value) reflect.Value {
	if field.CanInterface() {
		return field
	} else {
		// for all unexported fields (private)
		res := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		return res
	}
}

func SetFieldValue(field reflect.Value, value interface{}) {
	if field.CanInterface() && field.CanAddr() && field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		// for all unexported fields (private)
		reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
			Elem().
			Set(reflect.ValueOf(value))
	}
}

func GetFieldValueFromMethodAndObject[T interface{}](object T, name string) reflect.Value {
	v := reflect.ValueOf(&object).Elem()
	if v.Kind() == reflect.Ptr {
		val := v
		method := val.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		}
	} else if v.Kind() == reflect.Struct {
		method := v.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		} else {
			// https://www.geeksforgeeks.org/reflect-addr-function-in-golang-with-examples/
			pointerType := v.Addr()
			method := pointerType.MethodByName(name)
			res := method.Call(nil)
			return res[0]
		}
	}

	return *new(reflect.Value)
}

func GetFieldValueFromMethodAndReflectValue(val reflect.Value, name string) reflect.Value {
	if val.Kind() == reflect.Ptr {
		method := val.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		}
	} else if val.Kind() == reflect.Struct {
		method := val.MethodByName(name)
		if method.Kind() == reflect.Func {
			res := method.Call(nil)
			return res[0]
		} else {
			// https://www.geeksforgeeks.org/reflect-addr-function-in-golang-with-examples/
			pointerType := val.Addr()
			method := pointerType.MethodByName(name)
			res := method.Call(nil)
			return res[0]
		}
	}

	return *new(reflect.Value)
}

func SetValue[T any](data T, value interface{}) {
	var inputValue reflect.Value
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		inputValue = reflect.ValueOf(data).Elem()
	} else {
		inputValue = reflect.ValueOf(data)
	}

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Kind() == reflect.Ptr {
		inputValue.Set(valueReflect.Elem())
	} else {
		inputValue.Set(valueReflect)
	}
}

func ObjectTypePath(obj any) string {
	objType := reflect.TypeOf(obj).Elem()
	path := fmt.Sprintf("%s.%s", objType.PkgPath(), objType.Name())
	return path
}

func TypePath[T any]() string {
	var msg T
	return ObjectTypePath(msg)
}

func MethodPath(f interface{}) string {
	pointerName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	lastSlashIdx := strings.LastIndex(pointerName, "/")
	methodPath := pointerName[lastSlashIdx+1:]
	if methodPath[len(methodPath)-3:] == "-fm" {
		methodPath = methodPath[:len(methodPath)-3]
	}
	methodPath = strings.ReplaceAll(methodPath, ".", ":")
	methodPath = strings.ReplaceAll(methodPath, "(", "")
	methodPath = strings.ReplaceAll(methodPath, ")", "")
	methodPath = strings.ReplaceAll(methodPath, "*", "")
	return methodPath
}
