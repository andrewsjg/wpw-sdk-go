package errors

import (
	"fmt"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	var err = errors[UNKNOWN]
	if err.ID != "UNKNOWN" {
		t.Error("Wrong id " + err.ID)
		t.Fail()
	}
}

func TestGetError(t *testing.T) {
	for idx, val := range errors {
		err := GetError(ErrorId(idx))
		if err == nil {
			t.Error("GetError() returns nil: ", val)
			t.Fail()
		} else {
			fmt.Println(err)
		}
	}
	err := GetError(UNKNOWN)
	if err == nil {
		t.Error("nil not expected from GetError()")
		t.Fail()
	}

	additionalErrorData := [...]string{"some", "additional", "data"}
	err = GetError(UNKNOWN, additionalErrorData[0], additionalErrorData[1], additionalErrorData[2])
	if err == nil {
		t.Error("nil not expected from GetError()")
		t.Fail()
	}

	for _, errData := range additionalErrorData {
		if strings.Contains(err.Error(), errData) == false {
			t.Error("err.Error() does not contain additional data: " + errData)
			t.Fail()
		}
	}
}
