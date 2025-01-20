package mock

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	type UncommonType struct {
		Str    string
		StrPtr *string
		Slice  []int
	}
	type TestInterface interface {
		Ints(a, b int) int
		Variadic(s ...string) string
		NullableInput(a any, b *string)
		NullableOutput() error
		ManyArgs(a, b, c int16, d, e, f int32, g, h, i int64, j, k, l uintptr, x float32) (r1, r2 int, r3, r4 float32)
		UncommonArgs(UncommonType, *UncommonType) (UncommonType, *UncommonType, reflect.Type)
		StackArg([32]int)
		ConvertReturnValue() (int16, <-chan int)
	}
	i, m := Mock[TestInterface]()

	t.Run("ints", func(t *testing.T) {
		m.On("Ints", 1, 2).Return(3)
		assert.Equal(t, 3, i.Ints(1, 2))
		m.AssertCalled(t, "Ints", 1, 2)
	})

	t.Run("variadic", func(t *testing.T) {
		m.On("Variadic", []string{"a", "b"}).Once().Return("ab")
		assert.Equal(t, "ab", i.Variadic("a", "b"))
		m.AssertCalled(t, "Variadic", []string{"a", "b"})

		m.On("Variadic", nil).Once().Return("")
		assert.Equal(t, "", i.Variadic())
	})

	t.Run("nullable input", func(t *testing.T) {
		m.On("NullableInput", nil, nil).Return()
		i.NullableInput(nil, nil)
		m.AssertCalled(t, "NullableInput", nil, nil)
	})

	t.Run("nullable output", func(t *testing.T) {
		m.On("NullableOutput").Once().Return(assert.AnError)
		assert.Equal(t, assert.AnError, i.NullableOutput())
		m.AssertCalled(t, "NullableOutput")

		m.On("NullableOutput").Once().Return(nil)
		assert.NoError(t, i.NullableOutput())
	})

	t.Run("many args", func(t *testing.T) {
		m.On("ManyArgs",
			int16(1), int16(2), int16(3),
			int32(4), int32(5), int32(6),
			int64(7), int64(8), int64(9),
			uintptr(10), uintptr(11), uintptr(12),
			float32(0.5)).
			Return(1, 2, float32(3.0), float32(4.0))
		r1, r2, r3, r4 := i.ManyArgs(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0.5)
		assert.Equal(t, 1, r1)
		assert.Equal(t, 2, r2)
		assert.Equal(t, float32(3.0), r3)
		assert.Equal(t, float32(4.0), r4)
	})

	t.Run("uncommon args", func(t *testing.T) {
		typ := reflect.TypeFor[int]()
		hello := "hello"
		world := "world"
		uIn := UncommonType{
			Str:    hello,
			StrPtr: &hello,
			Slice:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		}
		uOut := UncommonType{
			Str:    world,
			StrPtr: &world,
			Slice:  []int{-1},
		}
		m.On("UncommonArgs", uIn, &uIn).Return(uOut, &uOut, typ)
		u_, uPtr_, typ_ := i.UncommonArgs(uIn, &uIn)
		assert.Equal(t, uOut, u_)
		assert.NotNil(t, uPtr_)
		assert.Equal(t, &uOut, uPtr_)
		assert.NotNil(t, typ_)
		assert.Equal(t, typ.Name(), typ_.Name())
		m.AssertCalled(t, "UncommonArgs", uIn, &uIn)
	})

	t.Run("stack arg", func(t *testing.T) {
		var arg [32]int
		for i := range arg {
			arg[i] = i
		}
		m.On("StackArg", arg).Return()
		i.StackArg(arg)
		m.AssertCalled(t, "StackArg", arg)
	})

	t.Run("convert return value", func(t *testing.T) {
		ch := make(chan int)
		m.On("ConvertReturnValue").Return(1, ch)
		_, _ = i.ConvertReturnValue()
		m.AssertCalled(t, "ConvertReturnValue")
	})

	t.Run("multiple instances", func(t *testing.T) {
		type Fooer interface {
			Foo() int
		}
		f1, m1 := Mock[Fooer]()
		f2, m2 := Mock[Fooer]()
		m1.On("Foo").Return(1)
		m2.On("Foo").Return(2)
		assert.Equal(t, 1, f1.Foo())
		assert.Equal(t, 2, f2.Foo())
		m1.AssertExpectations(t)
		m2.AssertExpectations(t)
	})
}
