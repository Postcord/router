package main

import (
	"github.com/Postcord/router"
	"testing"
)

func Test_add(t *testing.T) {
	_, b := builder()
	router.TestComponent(t, b, "/set/:number")
}
