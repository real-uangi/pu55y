package test

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	var a *string
	b := "hello"
	a = &b
	fmt.Printf("%s %d %d %s\n", *a, a, &b, b)

}
