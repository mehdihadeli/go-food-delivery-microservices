package typeMapper

// https://stackoverflow.com/a/34722791/581476
// https://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-run-time-in-go
// https://www.reddit.com/r/golang/comments/38u4j4/how_to_create_an_object_with_reflection/

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var (
	types    map[string][]reflect.Type
	packages map[string]map[string][]reflect.Type
)

// discoverTypes initializes types and packages
func init() {
	types = make(map[string][]reflect.Type)
	packages = make(map[string]map[string][]reflect.Type)

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
				// just discover pointer types, but we also register this pointer type actual struct type to the registry
				loadedTypePtr := typ
				loadedType := typ.Elem()

				pkgTypes := packages[loadedType.PkgPath()]
				pkgTypesPtr := packages[loadedTypePtr.PkgPath()]

				if pkgTypes == nil {
					pkgTypes = map[string][]reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				if pkgTypesPtr == nil {
					pkgTypesPtr = map[string][]reflect.Type{}
					packages[loadedTypePtr.PkgPath()] = pkgTypesPtr
				}

				if strings.Contains(loadedType.String(), "Test") {
					n := GetFullTypeNameByType(loadedType)
					n2 := GetTypeNameByType(loadedType)
					fmt.Println(n)
					fmt.Println(n2)
				}

				types[GetFullTypeNameByType(loadedType)] = append(
					types[GetFullTypeNameByType(loadedType)],
					loadedType,
				)
				types[GetFullTypeNameByType(loadedTypePtr)] = append(
					types[GetFullTypeNameByType(loadedTypePtr)],
					loadedTypePtr,
				)
				types[GetTypeNameByType(loadedType)] = append(
					types[GetTypeNameByType(loadedType)],
					loadedType,
				)
				types[GetTypeNameByType(loadedTypePtr)] = append(
					types[GetTypeNameByType(loadedTypePtr)],
					loadedTypePtr,
				)
			}
		}
	}
}

func RegisterType(typ reflect.Type) {
	types[GetFullTypeName(typ)] = append(types[GetFullTypeName(typ)], typ)
	types[GetTypeName(typ)] = append(types[GetTypeName(typ)], typ)
}

func RegisterTypeWithKey(key string, typ reflect.Type) {
	types[key] = append(types[key], typ)
}

func GetAllRegisteredTypes() map[string][]reflect.Type {
	return types
}

// TypeByName return the type by its name
func TypeByName(typeName string) reflect.Type {
	if typ, ok := types[typeName]; ok {
		return typ[0]
	}
	return nil
}

func TypesByName(typeName string) []reflect.Type {
	if types, ok := types[typeName]; ok {
		return types
	}
	return nil
}

func TypeByNameAndImplementedInterface[TInterface interface{}](typeName string) reflect.Type {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	implementedInterface := GetTypeFromGeneric[TInterface]()
	if types, ok := types[typeName]; ok {
		for _, t := range types {
			if t.Implements(implementedInterface) {
				return t
			}
		}
	}
	return nil
}

func TypesImplementedInterfaceWithFilterTypes[TInterface interface{}](
	types []reflect.Type,
) []reflect.Type {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	implementedInterface := GetTypeFromGeneric[TInterface]()

	var res []reflect.Type
	for _, t := range types {
		if t.Implements(implementedInterface) {
			res = append(res, t)
		}
	}

	return res
}

func TypesImplementedInterface[TInterface interface{}]() []reflect.Type {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	implementedInterface := GetTypeFromGeneric[TInterface]()

	var res []reflect.Type
	for _, t := range types {
		for _, v := range t {
			if v.Implements(implementedInterface) {
				res = append(res, v)
			}
		}
	}

	return res
}

// GetFullTypeName returns the full name of the type by its package name
func GetFullTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	return t.String()
}

func GetFullTypeNameByType(typ reflect.Type) string {
	return typ.String()
}

// GetTypeName returns the name of the type without its package name
func GetTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return fmt.Sprintf("*%s", t.Elem().Name())
}

func GetTypeNameByT[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return fmt.Sprintf("*%s", t.Elem().Name())
}

func GetNonPointerTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return t.Elem().Name()
}

func GetTypeNameByType(typ reflect.Type) string {
	if typ.Kind() != reflect.Ptr {
		return typ.Name()
	}

	return fmt.Sprintf("*%s", typ.Elem().Name())
}

// TypeByPackageName return the type by its package and name
func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name][0]
	}
	return nil
}

func TypesByPackageName(pkgPath string, name string) []reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes[name]
	}
	return nil
}

func GetTypeFromGeneric[T interface{}]() reflect.Type {
	res := reflect.TypeOf((*T)(nil)).Elem()
	return res
}

func GetBaseType(value interface{}) interface{} {
	if reflect.ValueOf(value).Kind() == reflect.Pointer {
		return reflect.ValueOf(value).Elem().Interface()
	}

	return value
}

func GetReflectType(value interface{}) reflect.Type {
	if reflect.TypeOf(value).Kind() == reflect.Pointer &&
		reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
		return reflect.TypeOf(value).Elem()
	}

	res := reflect.TypeOf(value)
	return res
}

func GetBaseReflectType(value interface{}) reflect.Type {
	if reflect.ValueOf(value).Kind() == reflect.Pointer {
		return reflect.TypeOf(reflect.ValueOf(value).Elem().Interface())
	}

	return reflect.TypeOf(value)
}

func GenericInstanceByT[T any]() T {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	typ := GetTypeFromGeneric[T]()
	return getInstanceFromType(typ).(T)
}

func InstanceByType(typ reflect.Type) interface{} {
	return getInstanceFromType(typ)
}

// InstanceByTypeName return an empty instance of the type by its name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return an empty struct
func InstanceByTypeName(name string) interface{} {
	typ := TypeByName(name)

	return getInstanceFromType(typ)
}

func InstanceByTypeNameAndImplementedInterface[TInterface interface{}](name string) interface{} {
	typ := TypeByNameAndImplementedInterface[TInterface](name)

	return getInstanceFromType(typ)
}

// InstancePointerByTypeName return an empty pointer instance of the type by its name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return a pointer to the struct
func InstancePointerByTypeName(name string) interface{} {
	typ := TypeByName(name)
	if typ.Kind() == reflect.Ptr {
		res := reflect.New(typ.Elem()).Interface()
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
		res := reflect.New(typ.Elem()).Interface()
		return res
	}

	return reflect.Zero(typ).Interface()
	// return reflect.New(typ).Elem().Interface()
}

func GetImplementInterfaceTypes[T any]() map[string][]reflect.Type {
	result := make(map[string][]reflect.Type)

	// Get the interface type
	interfaceType := reflect.TypeOf((*T)(nil)).Elem()

	// Iterate over the types in the map
	for groupName, typeList := range types {
		var implementingTypes []reflect.Type

		// Check each type in the list
		for _, t := range typeList {
			// Check if the type implements the interface
			if t.Implements(interfaceType) {
				implementingTypes = append(implementingTypes, t)
			}
		}

		if len(implementingTypes) > 0 {
			result[groupName] = implementingTypes
		}
	}

	return result
}
