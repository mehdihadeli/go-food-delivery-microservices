// Ref: https://github.com/erni27/mapper/
// https://github.com/alexsem80/go-mapper

package mapper

import (
	"emperror.dev/errors"
	"flag"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/logrous"
	reflectionHelper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/reflection_helper"
	"reflect"
)

var (
	// ErrNilFunction is the error returned by CreateCustomMap or CreateMapWith
	// if a nil function is passed to the method.
	ErrNilFunction = errors.New("mapper: nil function")
	// ErrMapNotExist is the error returned by the Map method
	// if a map for given types does not exist.
	ErrMapNotExist = errors.New("mapper: map does not exist")
	// ErrMapAlreadyExists is the error returned by one of the CreateMap method
	// if a given map already exists. Mapper does not allow to override MapFunc.
	ErrMapAlreadyExists = errors.New("mapper: map already exists")
	// ErrUnsupportedMap is the error returned by CreateMap or CreateMapWith
	// if a source - destination mapping is not supported. The mapping is supported only for
	// structs to structs with at least one exported field by a src side which corresponds to a dst field.
	ErrUnsupportedMap = errors.New("mapper: unsupported map")
)

const (
	SrcKeyIndex = iota
	DestKeyIndex
)

type mappingsEntry struct {
	SourceType      reflect.Type
	DestinationType reflect.Type
}

type typeMeta struct {
	keysToTags map[string]string
	tagsToKeys map[string]string
}

type MapFunc[TSrc any, TDst any] func(TSrc) TDst

var profiles = map[string][][2]string{}
var maps = map[mappingsEntry]interface{}{}
var mapperConfig *MapperConfig

func init() {
	mapperConfig = &MapperConfig{
		MapUnexportedFields: false,
	}
}

func Configure(config *MapperConfig) {
	mapperConfig = config
}

func CreateMap[TSrc any, TDst any]() error {
	var src TSrc
	var dst TDst
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&dst).Elem()

	if (srcType.Kind() != reflect.Struct && (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() != reflect.Struct)) || (desType.Kind() != reflect.Struct && (desType.Kind() == reflect.Ptr && desType.Elem().Kind() != reflect.Struct)) {
		return ErrUnsupportedMap
	}

	if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
		pointerStructTypeKey := mappingsEntry{SourceType: srcType, DestinationType: desType}
		nonePointerStructTypeKey := mappingsEntry{SourceType: srcType.Elem(), DestinationType: desType.Elem()}
		if _, ok := maps[nonePointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}
		if _, ok := maps[pointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}

		// add pointer struct map and none pointer struct map to registry
		maps[nonePointerStructTypeKey] = nil
		maps[pointerStructTypeKey] = nil
	} else {
		nonePointerStructTypeKey := mappingsEntry{SourceType: srcType, DestinationType: desType}
		pointerStructTypeKey := mappingsEntry{SourceType: reflect.New(srcType).Type(), DestinationType: reflect.New(desType).Type()}
		if _, ok := maps[nonePointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}
		if _, ok := maps[pointerStructTypeKey]; ok {
			return ErrMapAlreadyExists
		}

		// add pointer struct map and none pointer struct map to registry
		maps[nonePointerStructTypeKey] = nil
		maps[pointerStructTypeKey] = nil
	}

	if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
		srcType = srcType.Elem()
	}

	if desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Struct {
		desType = desType.Elem()
	}

	configProfile(srcType, desType)

	return nil
}

func CreateCustomMap[TSrc any, TDst any](fn MapFunc[TSrc, TDst]) error {
	if fn == nil {
		return ErrNilFunction
	}
	var src TSrc
	var dst TDst
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&dst).Elem()

	if (srcType.Kind() != reflect.Struct && (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() != reflect.Struct)) || (desType.Kind() != reflect.Struct && (desType.Kind() == reflect.Ptr && desType.Elem().Kind() != reflect.Struct)) {
		return ErrUnsupportedMap
	}

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	if _, ok := maps[k]; ok {
		return ErrMapAlreadyExists
	}
	maps[k] = fn

	//if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
	//	srcType = srcType.Elem()
	//}
	//
	//if desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Struct {
	//	desType = desType.Elem()
	//}

	return nil
}

func Map[TDes any, TSrc any](src TSrc) (TDes, error) {
	if reflect.ValueOf(src).IsZero() {
		return *new(TDes), nil
	}

	var des TDes
	srcType := reflect.TypeOf(src)
	desType := reflect.TypeOf(des)
	desIsArray := false
	srcIsArray := false

	if srcType.Kind() == reflect.Array || (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Array) || srcType.Kind() == reflect.Slice || (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Slice) {
		srcType = srcType.Elem()
		srcIsArray = true
	}

	if desType.Kind() == reflect.Array || (desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Array) || desType.Kind() == reflect.Slice || (desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Slice) {
		desType = desType.Elem()
		desIsArray = true
	}

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	fn, ok := maps[k]
	if !ok {
		return *new(TDes), ErrMapNotExist
	}
	if fn != nil {
		fnReflect := reflect.ValueOf(fn)

		if desIsArray && srcIsArray {
			linq.From(src).Select(func(x interface{}) interface{} {
				return fnReflect.Call([]reflect.Value{reflect.ValueOf(x)})[0].Interface()
			}).ToSlice(&des)

			return des, nil
		} else {
			return fnReflect.Call([]reflect.Value{reflect.ValueOf(src)})[0].Interface().(TDes), nil
		}
	}

	desTypeValue := reflect.ValueOf(&des).Elem()
	err := processValues[TDes, TSrc](reflect.ValueOf(src), desTypeValue)
	if err != nil {
		return *new(TDes), err
	}

	return des, nil
}

func configProfile(srcType reflect.Type, destType reflect.Type) {
	// parse logger flags
	flag.Parse()

	// check for provided types kind.
	// if not struct - skip.
	if srcType.Kind() != reflect.Struct {
		logrous.DefaultLogger.Errorf("expected reflect.Struct kind for type %s, but got %s", srcType.String(), srcType.Kind().String())
	}

	if destType.Kind() != reflect.Struct {
		logrous.DefaultLogger.Errorf("expected reflect.Struct kind for type %s, but got %s", destType.String(), destType.Kind().String())
	}

	// profile is slice of src and dest structs fields names
	var profile [][2]string

	// get types metadata
	srcMeta := getTypeMeta(srcType)
	destMeta := getTypeMeta(destType)
	srcMethods := getTypeMethods(srcType)

	for srcKey, srcTag := range srcMeta.keysToTags {
		if _, ok := destMeta.keysToTags[strcase.ToCamel(srcKey)]; ok {
			profile = append(profile, [2]string{srcKey, strcase.ToCamel(srcKey)})
		}

		// case src key equals dest key
		if _, ok := destMeta.keysToTags[srcKey]; ok {
			profile = append(profile, [2]string{srcKey, srcKey})
			continue
		}

		// case src key equals dest tag
		if destKey, ok := destMeta.tagsToKeys[srcKey]; ok {
			profile = append(profile, [2]string{srcKey, destKey})
			continue
		}

		// case src tag equals dest key
		if _, ok := destMeta.keysToTags[srcTag]; ok {
			profile = append(profile, [2]string{srcKey, srcTag})
			continue
		}

		// case src tag equals dest tag
		if destKey, ok := destMeta.tagsToKeys[srcTag]; ok {
			profile = append(profile, [2]string{srcKey, destKey})
			continue
		}
	}

	for _, method := range srcMethods {
		if _, ok := destMeta.keysToTags[method]; ok {
			profile = append(profile, [2]string{method, method})
			continue
		}
	}

	// save profile with unique srcKey for provided types
	profiles[getProfileKey(srcType, destType)] = profile
}

func getProfileKey(srcType reflect.Type, destType reflect.Type) string {
	return fmt.Sprintf("%s_%s", srcType.Name(), destType.Name())
}

func getTypeMeta(val reflect.Type) typeMeta {
	fieldsNum := val.NumField()

	keysToTags := make(map[string]string)
	tagsToKeys := make(map[string]string)

	for i := 0; i < fieldsNum; i++ {
		field := val.Field(i)
		fieldName := field.Name
		fieldTag := field.Tag.Get("mapper")

		keysToTags[fieldName] = fieldTag
		if fieldTag != "" {
			tagsToKeys[fieldTag] = fieldName
		}
	}

	return typeMeta{
		keysToTags: keysToTags,
		tagsToKeys: tagsToKeys,
	}
}

func getTypeMethods(val reflect.Type) []string {
	methodsNum := val.NumMethod()
	var keys []string

	for i := 0; i < methodsNum; i++ {

		methodName := val.Method(i).Name
		keys = append(keys, methodName)
	}

	return keys
}

// mapStructs func perform structs casts.
func mapStructs[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	// get values types
	// if types or their slices were not registered - abort
	profile, ok := profiles[getProfileKey(src.Type(), dest.Type())]
	if !ok {
		logrous.DefaultLogger.Errorf("no conversion specified for types %s and %s", src.Type().String(), dest.Type().String())
		return
	}

	// iterate over struct fields and map values
	for _, keys := range profile {
		destinationField := dest.FieldByName(keys[DestKeyIndex])
		sourceField := src.FieldByName(keys[SrcKeyIndex])
		var sourceFiledValue reflect.Value

		if sourceField.Kind() != reflect.Invalid {
			//var destinationFieldValue reflect.Value
			if !sourceField.CanInterface() {
				if mapperConfig.MapUnexportedFields {
					sourceFiledValue = reflectionHelper.GetFieldValue(sourceField)
				} else {
					// for getting pointer for non-pointer struct we can use reflect.Addr() for calling pointer receivers properties
					sourceFiledValue = reflectionHelper.GetFieldValueFromMethodAndReflectValue(src.Addr(), strcase.ToCamel(keys[SrcKeyIndex]))
				}
			} else {
				if mapperConfig.MapUnexportedFields {
					sourceFiledValue = reflectionHelper.GetFieldValue(sourceField)
				} else {
					sourceFiledValue = sourceField
				}
			}
		} else {
			// there is no field corresponding to destination filed, so we search on source methods (properties) for getting src field value for example `Id()` property
			sourceFiledValue = reflectionHelper.GetFieldValueFromMethodAndReflectValue(src.Addr(), strcase.ToCamel(keys[SrcKeyIndex]))
		}

		processValues[TDes, TSrc](sourceFiledValue, destinationField)
	}
}

// mapSlices func perform slices casts.
func mapSlices[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	// Make dest slice
	dest.Set(reflect.MakeSlice(dest.Type(), src.Len(), src.Cap()))

	// Get each element of slice
	// process values mapping
	for i := 0; i < src.Len(); i++ {
		srcVal := src.Index(i)
		destVal := dest.Index(i)

		processValues[TDes, TSrc](srcVal, destVal)
	}
}

// mapPointers func perform pointers casts.
func mapPointers[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	// create new struct from provided dest type
	val := reflect.New(dest.Type().Elem()).Elem()

	processValues[TDes, TSrc](src.Elem(), val)

	// assign address of instantiated struct to destination
	dest.Set(val.Addr())
}

// mapMaps func perform maps casts.
func mapMaps[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	// Make dest map
	dest.Set(reflect.MakeMapWithSize(dest.Type(), src.Len()))

	// Get each element of map as key-values
	// process keys and values mapping and update dest map
	srcMapIter := src.MapRange()
	destMapIter := dest.MapRange()

	for destMapIter.Next() && srcMapIter.Next() {
		destKey := reflect.New(destMapIter.Key().Type()).Elem()
		destValue := reflect.New(destMapIter.Value().Type()).Elem()
		processValues[TDes, TSrc](srcMapIter.Key(), destKey)
		processValues[TDes, TSrc](srcMapIter.Value(), destValue)

		dest.SetMapIndex(destKey, destValue)
	}
}

func processValues[TDes any, TSrc any](src reflect.Value, dest reflect.Value) error {
	// if src of dest is an interface - get underlying type
	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}

	if dest.Kind() == reflect.Interface {
		dest = dest.Elem()
	}

	// get provided values' kinds
	srcKind := src.Kind()
	destKind := dest.Kind()

	// skip invalid kinds
	if srcKind == reflect.Invalid || destKind == reflect.Invalid {
		return nil
	}

	// check if kinds are equal
	if srcKind != destKind {
		// TODO dynamic cast, m.b. with Mapper extensions
		return nil
	}

	// if types are equal set dest value
	if src.Type() == dest.Type() {
		reflectionHelper.SetFieldValue(dest, src.Interface())
		//dest.Set(src)
		return nil
	}

	// resolve kind and choose mapping function
	// or set dest value
	switch src.Kind() {
	case reflect.Struct:
		mapStructs[TDes, TSrc](src, dest)
	case reflect.Slice:
		mapSlices[TDes, TSrc](src, dest)
	case reflect.Map:
		mapMaps[TDes, TSrc](src, dest)
	case reflect.Ptr:
		mapPointers[TDes, TSrc](src, dest)
	default:
		dest.Set(src)
	}

	return nil
}
