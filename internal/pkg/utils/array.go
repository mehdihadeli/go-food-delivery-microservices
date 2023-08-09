package utils

import "github.com/goccy/go-reflect"

func Contains[T any](arr []T, x T) bool {
	for _, v := range arr {
		if reflect.ValueOf(v) == reflect.ValueOf(x) {
			return true
		}
	}
	return false
}

func ContainsFunc[T any](arr []T, predicate func(T) bool) bool {
	for _, v := range arr {
		if predicate(v) == true {
			return true
		}
	}
	return false
}
