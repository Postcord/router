package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModalRouter_prep(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		m := ModalRouter{}
		m.prep()
		assert.NotNil(t, m.routes)
	})
	t.Run("map exists", func(t *testing.T) {
		m := ModalRouter{routes: map[string]*ModalContent{"a": nil}}
		m.prep()
		assert.Equal(t, map[string]*ModalContent{"a": nil}, m.routes)
	})
}
