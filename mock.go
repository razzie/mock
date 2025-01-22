package mock

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/stretchr/testify/mock"
)

const (
	// Anything is used when the argument being tested
	// shouldn't be taken into consideration.
	Anything = mock.Anything

	ptrSize = 4 << (^uintptr(0) >> 63)
)

var (
	mu    sync.Mutex
	cache map[reflect.Type]func() (any, *mock.Mock)
)

// Mock takes a T interface type and returns a runtime generated mock implementation for it
// plus a testify mock.Mock object to control its behavior
func Mock[T any]() (T, *mock.Mock) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if typ.Kind() != reflect.Interface {
		panic("not an interface type: " + typ.Name())
	}

	mu.Lock()
	defer mu.Unlock()
	if cache == nil {
		cache = make(map[reflect.Type]func() (any, *mock.Mock))
	}
	cached, ok := cache[typ]
	if ok {
		obj, m := cached()
		return obj.(T), m
	}

	mockType := reflect.StructOf([]reflect.StructField{
		{
			Name: "MockMethodContexts__" + typ.Name(),
			Type: reflect.TypeFor[unsafe.Pointer](),
		},
		{
			Name: "Mock__" + typ.Name(),
			Type: reflect.TypeFor[*mock.Mock](),
		},
		{
			Name:      typ.Name(),
			Type:      typ,
			Anonymous: true,
		},
	})
	contexts := make([]unsafe.Pointer, typ.NumMethod())
	abiMethods := getAbiMethods(mockType)
	for i := range abiMethods {
		fn := mockMethod(typ.Method(i))
		contexts[i] = unpackIfaceData(fn)
		abiMethods[i].setFn(mockMethodWrappers[i])
	}

	newMock := func() (any, *mock.Mock) {
		contextsPtr := unsafe.Pointer(unsafe.SliceData(contexts))
		m := new(mock.Mock)
		obj := reflect.New(mockType).Elem()
		obj.Field(0).Set(reflect.ValueOf(contextsPtr))
		obj.Field(1).Set(reflect.ValueOf(m))
		return obj.Interface(), m
	}
	cache[typ] = newMock

	obj, m := newMock()
	return obj.(T), m
}

func mockMethod(method reflect.Method) any {
	inTypes := getInTypes(method)
	outTypes := getOutTypes(method)
	fnType := reflect.FuncOf(inTypes, outTypes, method.Type.IsVariadic())
	fnImpl := func(inValues []reflect.Value) (outValues []reflect.Value) {
		recv := inValues[0]
		in := make([]any, len(inValues)-1)
		for i, val := range inValues[1:] {
			switch {
			case val.IsZero():
				in[i] = nil
			default:
				in[i] = val.Interface()
			}
		}
		m := *(**mock.Mock)(unsafe.Add(recv.UnsafePointer(), ptrSize))
		out := m.MethodCalled(method.Name, in...)
		if len(out) != len(outTypes) {
			panic("mock: invalid number of outputs")
		}
		outValues = make([]reflect.Value, len(out))
		for i, val := range out {
			switch {
			case val == nil:
				outValues[i] = reflect.Zero(outTypes[i])
			case reflect.TypeOf(val) != outTypes[i]:
				outValues[i] = reflect.ValueOf(val).Convert(outTypes[i])
			default:
				outValues[i] = reflect.ValueOf(val)
			}
		}
		return
	}
	fnVal := reflect.MakeFunc(fnType, fnImpl)
	return fnVal.Interface()
}

func getInTypes(method reflect.Method) []reflect.Type {
	recv := reflect.TypeFor[unsafe.Pointer]()
	types := make([]reflect.Type, method.Type.NumIn()+1)
	types[0] = recv
	for i := range types[1:] {
		types[i+1] = method.Type.In(i)
	}
	return types
}

func getOutTypes(method reflect.Method) []reflect.Type {
	types := make([]reflect.Type, method.Type.NumOut())
	for i := range types {
		types[i] = method.Type.Out(i)
	}
	return types
}

func unpackIfaceData(iface any) unsafe.Pointer {
	return (*(*[2]unsafe.Pointer)(unsafe.Pointer(&iface)))[1]
}

var mockMethodWrappers = []func(){
	m0, m1, m2, m3, m4, m5, m6, m7,
	m8, m9, m10, m11, m12, m13, m14, m15,
	m16, m17, m18, m19, m20, m21, m22, m23,
	m24, m25, m26, m27, m28, m29, m30, m31,
	m32, m33, m34, m35, m36, m37, m38, m39,
	m40, m41, m42, m43, m44, m45, m46, m47,
	m48, m49, m50, m51, m52, m53, m54, m55,
	m56, m57, m58, m59, m60, m61, m62, m63,
}

func m0()
func m1()
func m2()
func m3()
func m4()
func m5()
func m6()
func m7()
func m8()
func m9()
func m10()
func m11()
func m12()
func m13()
func m14()
func m15()
func m16()
func m17()
func m18()
func m19()
func m20()
func m21()
func m22()
func m23()
func m24()
func m25()
func m26()
func m27()
func m28()
func m29()
func m30()
func m31()
func m32()
func m33()
func m34()
func m35()
func m36()
func m37()
func m38()
func m39()
func m40()
func m41()
func m42()
func m43()
func m44()
func m45()
func m46()
func m47()
func m48()
func m49()
func m50()
func m51()
func m52()
func m53()
func m54()
func m55()
func m56()
func m57()
func m58()
func m59()
func m60()
func m61()
func m62()
func m63()
