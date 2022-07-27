package typeMapper

import (
	"reflect"
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
				loadedType := typ.Elem()
				pkgTypes := packages[loadedType.PkgPath()]
				if pkgTypes == nil {
					pkgTypes = map[string]reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				types[loadedType.String()] = loadedType
				pkgTypes[loadedType.Name()] = loadedType
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
	if t := reflect.TypeOf(input); t.Kind() == reflect.Ptr {
		return t.Elem().String()
	} else {
		return t.String()
	}
}

// TypeByPackageName return the type by its package and name
func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name]
	}
	return nil
}

//https://stackoverflow.com/a/34722791/581476

// TypeInstanceByName return an empty instance of the type by its name
func TypeInstanceByName(name string) interface{} {
	//https://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-run-time-in-go
	//https://www.reddit.com/r/golang/comments/38u4j4/how_to_create_an_object_with_reflection/
	return reflect.New(TypeByName(name)).Elem().Interface()
}

// TypePointerInstanceByName return an empty pointer instance of the type by its name
func TypePointerInstanceByName(name string) interface{} {
	//https://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-run-time-in-go
	//https://www.reddit.com/r/golang/comments/38u4j4/how_to_create_an_object_with_reflection/
	return reflect.New(TypeByName(name)).Interface()
}

// GenericInstanceTypeByName return an empty instance of the generic type by its name
func GenericInstanceTypeByName[T any](typeName string) T {
	return TypeInstanceByName(typeName).(T)
}

// TypeInstanceByPackageName return an empty instance of the type by its name and package name
func TypeInstanceByPackageName(pkgPath string, name string) interface{} {
	return reflect.New(TypeByPackageName(pkgPath, name)).Elem().Interface()
}
