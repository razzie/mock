package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Fooer interface {
	Foo() error
}

func TestMock(t *testing.T) {
	f, m := Mock[Fooer]()
	m.On("Foo").Return(assert.AnError)
	err := f.Foo()
	assert.Equal(t, assert.AnError, err)
}
