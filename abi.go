package mock

import (
	"reflect"
	"unsafe"
)

//go:linkname addReflectOff reflect.addReflectOff
func addReflectOff(ptr unsafe.Pointer) int32

type abiMethod struct {
	Name int32 // name of method
	Mtyp int32 // method type (without receiver)
	Ifn  int32 // fn used in interface call (one-word receiver)
	Tfn  int32 // fn used for normal method call
}

func (m *abiMethod) setFn(fn any) {
	off := addReflectOff(reflect.ValueOf(fn).UnsafePointer())
	m.Ifn = off
}

func getAbiMethods(typ reflect.Type) []abiMethod {
	if typ.Kind() != reflect.Struct {
		panic("not struct kind: " + typ.Kind().String())
	}
	type Type struct {
		Size_       uintptr
		PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
		Hash        uint32  // hash of type; avoids computation in hash tables
		TFlag       uint8   // extra type information flags
		Align_      uint8   // alignment of variable with this type
		FieldAlign_ uint8   // alignment of struct field with this type
		Kind_       uint8   // enumeration for C
		Equal       func(unsafe.Pointer, unsafe.Pointer) bool
		GCData      *byte
		Str         int32 // string form
		PtrToThis   int32 // type for pointer to this type, may be zero
	}
	type UncommonType struct {
		PkgPath int32  // import path; empty for built-in types like int, string
		Mcount  uint16 // number of methods
		Xcount  uint16 // number of exported methods
		Moff    uint32 // offset from this uncommontype to [mcount]Method
		_       uint32 // unused
	}
	type StructField struct {
		Name   *byte   // name is always non-empty
		Typ    *Type   // type of field
		Offset uintptr // byte offset of field
	}
	type StructType struct {
		Type
		PkgPath *byte
		Fields  []StructField
	}
	type StructTypeUncommon struct {
		StructType
		u UncommonType
	}
	t := unpackIfaceData(typ)
	u := &(*StructTypeUncommon)(t).u
	if u.Mcount == 0 {
		return nil
	}
	m := unsafe.Add(unsafe.Pointer(u), u.Moff)
	return unsafe.Slice((*abiMethod)(m), u.Mcount)
}
