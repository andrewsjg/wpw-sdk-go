package configuration

import (
	"os"
	"strings"
	"testing"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/wpwerrors"
)

var testFileName = os.TempDir() + string(os.PathSeparator) + "configload_test.json"
var goodConfigFileContent = "[{ \"key\": \"wsLogEnable\", \"value\": \"true\" }, { \"key\": \"wsLogPort\", \"value\": \"8182\" }, { \"key\": \"wsLogLevel\", \"value\": \"debug,error,fatal\" }]"
var badConfigFileContent = " } not json file { "

func setupConfigFile(t *testing.T, fileContent string) {

	var file, err = os.OpenFile(testFileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		t.Error("Cannot open new file: " + testFileName + ", " + err.Error())
		t.FailNow()
	}

	defer file.Close()

	_, err = file.WriteString(fileContent)

	if err != nil {
		t.Error("Cannot write to file: " + testFileName + ", " + err.Error())
		t.FailNow()
	}
	err = file.Sync()
	if err != nil {
		t.Error("Cannot sync file, " + err.Error())
		t.FailNow()
	}
}

func tearDownConfigFile(t *testing.T) {
	var err = os.Remove(testFileName)
	if err != nil {
		t.Error("Cannot remove file, " + err.Error())
		t.FailNow()
	}
}

func testBadFile(t *testing.T) {

	setupConfigFile(t, badConfigFileContent)
	defer tearDownConfigFile(t)

	_, err := Load(testFileName)

	if err == nil {
		t.Error("Load should file, the content of the file is not json.")
		t.FailNow()
	}
	if !strings.Contains(err.Error(), wpwerrors.GetError(wpwerrors.DECODE_JSON).Error()) {
		t.Error("Wrong error. Expected error is DECODE_JSON, got: " + err.Error())
		t.FailNow()
	}
}

func testGoodFile(t *testing.T) Configuration {

	setupConfigFile(t, goodConfigFileContent)
	defer tearDownConfigFile(t)

	ret, err := Load(testFileName)
	if err != nil {
		t.Error("Not expected error for load existing file.")
		t.FailNow()
	}

	if ret.items == nil {
		t.Error("Items cannot be nil after load configuration of existing file.")
		t.FailNow()
	}
	return ret
}

func TestLoad(t *testing.T) {

	var err error

	_, err = Load("")
	if err == nil {
		t.Error("Expected error for null string.")
		t.FailNow()
	}

	if !strings.Contains(err.Error(), wpwerrors.GetError(wpwerrors.WRONG_CONFIG_PATH).Error()) {
		t.Error("Wrong error. Expected error is WRONG_CONFIG_PATH, got: " + err.Error())
		t.FailNow()
	}

	_, err = Load("some not existing file")
	if err == nil {
		t.Error("Expected error for load not existing file.")
		t.FailNow()
	}

	testBadFile(t)
	testGoodFile(t)

}

func TestGetValue(t *testing.T) {
	conf := testGoodFile(t)

	item := conf.GetValue("notExistentItem")
	if item.Value != "" {
		t.Error("Expected item should be nil.")
		t.FailNow()
	}

	item = conf.GetValue("wsLogEnable")
	if item.Value != "true" {
		t.Error("Expected item is not true.")
		t.FailNow()
	}

	item = conf.GetValue("wsLogPort")
	if item.Value != "8182" {
		t.Error("wsLogPort should be 8182.")
		t.FailNow()
	}

	item = conf.GetValue("wsLogLevel")
	if item.Value != "debug,error,fatal" {
		t.Error("wsLogPort should be debug,error,fatal.")
		t.FailNow()
	}
}

func TestGetItems(t *testing.T) {
	conf := testGoodFile(t)

	items := conf.GetItems()
	if items == nil {
		t.Error("Items is nil.")
		t.FailNow()
	}

	if items["wsLogEnable"].Value != "true" {
		t.Error("wsLogEnable should be true.")
		t.FailNow()
	}

	if items["wsLogPort"].Value != "8182" {
		t.Error("wsLogPort should be 8182.")
		t.FailNow()
	}

	if items["wsLogLevel"].Value != "debug,error,fatal" {
		t.Error("wsLogPort should be debug,error,fatal.")
		t.FailNow()
	}
}
