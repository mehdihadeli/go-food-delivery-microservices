package typeMapper

import "unsafe"

//go:linkname typelinks2 reflect.typelinks
func typelinks2() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

type emptyInterface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}
