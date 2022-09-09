package typeMapper

//https://stackoverflow.com/a/34722791/581476
//https://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-run-time-in-go
//https://www.reddit.com/r/golang/comments/38u4j4/how_to_create_an_object_with_reflection/

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var types map[string]reflect.Type
var packages map[string]map[string]reflect.Type

// discoverTypes initializes types and packages
func init() {
	types = make(map[string]reflect.Type)
	packages = make(map[string]map[string]reflect.Type)

	discoverTypes()
}

func discoverTypes() {
	typ := reflect.TypeOf(0)
	sections, offset := typelinks2()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			emptyInterface := (*emptyInterface)(unsafe.Pointer(&typ))
			emptyInterface.data = resolveTypeOff(rodata, off)
			if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {

				// by defaultLogger just discover pointer types, but we also register this pointer type actual struct type to the registry
				loadedTypePtr := typ
				loadedType := typ.Elem()

				pkgTypes := packages[loadedType.PkgPath()]
				pkgTypesPtr := packages[loadedTypePtr.PkgPath()]

				if pkgTypes == nil {
					pkgTypes = map[string]reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				if pkgTypesPtr == nil {
					pkgTypesPtr = map[string]reflect.Type{}
					packages[loadedTypePtr.PkgPath()] = pkgTypesPtr
				}
				f := strings.Contains(loadedType.String(), "Test")
				if f {
					fmt.Println(loadedType.String())
				}

				types[loadedType.String()] = loadedType
				types[loadedTypePtr.String()] = loadedTypePtr
				pkgTypes[loadedType.Name()] = loadedType
				pkgTypesPtr[loadedTypePtr.Name()] = loadedTypePtr
			}
		}
	}
}

// TypeByName return the type by its name
func TypeByName(typeName string) reflect.Type {
	if typ, ok := types[typeName]; ok {
		return typ
	}
	return nil
}

func GetTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	return t.String()
}

// TypeByPackageName return the type by its package and name
func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name]
	}
	return nil
}

// InstanceByTypeName return an empty instance of the type by its name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return an empty struct
func InstanceByTypeName(name string) interface{} {
	typ := TypeByName(name)

	return getInstanceFromType(typ)
}

// InstancePointerByTypeName return an empty pointer instance of the type by its name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return a pointer to the struct
func InstancePointerByTypeName(name string) interface{} {
	typ := TypeByName(name)
	if typ.Kind() == reflect.Ptr {
		var res = reflect.New(typ.Elem()).Interface()
		return res
	}

	return reflect.New(typ).Interface()
}

// InstanceByPackageName return an empty instance of the type by its name and package name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return an empty struct
func InstanceByPackageName(pkgPath string, name string) interface{} {
	typ := TypeByPackageName(pkgPath, name)

	return getInstanceFromType(typ)
}

func getInstanceFromType(typ reflect.Type) interface{} {
	if typ.Kind() == reflect.Ptr {
		var res = reflect.New(typ.Elem()).Interface()
		return res
	}

	return reflect.Zero(typ).Interface()
	// return reflect.New(typ).Elem().Interface()
}

// GenericInstanceByTypeName return an empty instance of the generic type by its name
func GenericInstanceByTypeName[T any](typeName string) T {
	var res = InstanceByTypeName(typeName).(T)

	return res
}
