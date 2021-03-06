package wpwerrors

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
		err := GetError(ErrorID(idx))
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

	errDummy := fmt.Errorf("dummy in golang format")
	err = GetError(UNKNOWN, errDummy, err)
	if err == nil {
		t.Error("nil not expected from GetError()")
		t.Fail()
	}
	// error should contain brackets, because includes other (errDummy) error
	if strings.Contains(err.Error(), openBracket) == false {
		t.Error("err.Error() does not contain open bracket:" + err.Error())
		t.Fail()
	}
	if strings.Contains(err.Error(), closeBracket) == false {
		t.Error("err.Error() does not contain close bracket: " + err.Error())
		t.Fail()
	}

	// not all type of variables are supported yet
	err = GetError(UNKNOWN, 5)
	if err == nil {
		t.Error("nil not expected from GetError()")
		t.Fail()
	}
	if strings.Contains(err.Error(), unsupportedTypeError) == false {
		t.Error("err.Error() does not contain " + unsupportedTypeError + ": " + err.Error())
		t.Fail()
	}

	err = GetError(UNKNOWN, "dummy string", 0.5)
	if err == nil {
		t.Error("nil not expected from GetError()")
		t.Fail()
	}
	if strings.Contains(err.Error(), unsupportedTypeError) == false {
		t.Error("err.Error() does not contain " + unsupportedTypeError + ": " + err.Error())
		t.Fail()
	}
}
