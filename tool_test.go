package coapmq_test

import (
	"testing"

	. "github.com/kkdai/coapmq"
)

func TestSliceElementRemvoe(t *testing.T) {
	var slice []string
	ret := RemoveStringFromSlice(slice, "")
	if len(ret) != 0 {
		t.Error("Remove empty failed")
	}

	slice = append(slice, "a")
	ret = RemoveStringFromSlice(slice, "b")
	if len(ret) != 1 {
		t.Error("Remove only one failed")
	}

	slice = append(slice, "b")
	ret = RemoveStringFromSlice(slice, "a")
	if len(ret) != 1 || ret[0] != "b" {
		t.Error("Remove failed")
	}

	ret = RemoveStringFromSlice(ret, "b")
	if len(ret) != 0 {
		t.Error("Remove remain item failed")
	}
}
