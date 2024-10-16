package js

import (
	"fmt"
	"testing"
)

func TestOr(t *testing.T) {

	res0 := Or(0, false, nil, "")
	res1 := Or(0, false, "", 666, "abc")

	require(t, res0 == nil)
	require(t, res1 == 666)
}

func TestAnd(t *testing.T) {

	res0 := And(0, false, "", 666, "abc")
	res1 := And(666, -0.123, true, "abc")

	require(t, res0 == false)
	require(t, res1 == true)
}

func require(t *testing.T, ok bool) {
	if !ok {
		t.Fail()
	}
}

func call(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	fn()
	return
}
