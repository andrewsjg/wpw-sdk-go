package configuration

import (
	"strings"
	"testing"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/wpwerrors"
)

func TestReadInt(t *testing.T) {
	var item Item
	var value int
	var err error

	item.Key = "a"
	item.Value = "b"

	_, err = item.ReadInt()
	if err == nil {
		t.Error("ReadInt() should fail on not integer value.")
		t.FailNow()
	}
	if value != 0 {
		t.Error("ReadInt() should return 0 for error case.")
		t.FailNow()
	}
	if !strings.Contains(err.Error(), wpwerrors.GetError(wpwerrors.CONVERT_VALUE).Error()) {
		t.Error("Wrong error. Expected error is CONVERT_VALUE, got: " + err.Error())
		t.FailNow()
	}

	item.Value = "102"
	value, err = item.ReadInt()
	if err != nil {
		t.Error("ReadInt() failed on reading integer.")
		t.FailNow()
	}
	if value != 102 {
		t.Error("ReadInt() should return correct integer value.")
		t.FailNow()
	}
}

func TestReadBool(t *testing.T) {
	var item Item
	var value bool
	var err error

	item.Key = "a"
	item.Value = "b"

	_, err = item.ReadBool()
	if err == nil {
		t.Error("ReadBool() should fail on not boolean value.")
		t.FailNow()
	}
	if value != false {
		t.Error("ReadBool() should return false for error case.")
		t.FailNow()
	}

	item.Value = "true"
	value, err = item.ReadBool()
	if err != nil {
		t.Error("ReadBool() failed on reading boolean.")
		t.FailNow()
	}
	if value != true {
		t.Error("ReadBool() should return correct boolean value.")
		t.FailNow()
	}

	item.Value = "false"
	value, err = item.ReadBool()
	if err != nil {
		t.Error("ReadBool() failed on reading boolean.")
		t.FailNow()
	}
	if value != false {
		t.Error("ReadBool() should return correct boolean value.")
		t.FailNow()
	}
}
