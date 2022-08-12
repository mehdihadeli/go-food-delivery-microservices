package reflectionHelper

// ref: https://gist.github.com/drewolson/4771479
// https://stackoverflow.com/a/60598827/581476
// https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field

import (
	"reflect"
	"unsafe"
)

func GetFieldValueByIndex[T any](object T, index int) interface{} {
	if reflect.ValueOf(object).Kind() == reflect.Ptr {
		val := reflect.ValueOf(object).Elem()
		field := val.Field(index)
		// for all exported fields (public)
		if field.CanInterface() {
			return field.Interface()
		} else {
			// for all unexported fields (private)
			return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
		}
	} else if reflect.ValueOf(object).Kind() == reflect.Struct {
		// for all exported fields (public)
		val := reflect.ValueOf(object)
		field := val.Field(index)
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
	if reflect.ValueOf(object).Kind() == reflect.Ptr {
		val := reflect.ValueOf(object).Elem()
		field := val.FieldByName(name)
		// for all exported fields (public)
		if field.CanInterface() {
			return field.Interface()
		} else {
			// for all unexported fields (private)
			return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
		}
	} else if reflect.ValueOf(object).Kind() == reflect.Struct {
		// for all exported fields (public)
		val := reflect.ValueOf(object)
		field := val.FieldByName(name)
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
	//https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
	if reflect.ValueOf(object).Kind() == reflect.Ptr {
		val := reflect.ValueOf(object).Elem()
		field := val.Field(index)
		// for all exported fields (public)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			// for all unexported fields (private)
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
		}
	} else if reflect.ValueOf(object).Kind() == reflect.Struct {
		// for all exported fields (public)
		val := reflect.ValueOf(&object).Elem()
		field := val.Field(index)
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
	//https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
	if reflect.ValueOf(object).Kind() == reflect.Ptr {
		val := reflect.ValueOf(object).Elem()
		field := val.FieldByName(name)
		// for all exported fields (public)
		if field.CanInterface() && field.CanAddr() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		} else {
			// for all unexported fields (private)
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
		}
	} else if reflect.ValueOf(object).Kind() == reflect.Struct {
		// for all exported fields (public)
		val := reflect.ValueOf(&object).Elem()
		field := val.FieldByName(name)
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
