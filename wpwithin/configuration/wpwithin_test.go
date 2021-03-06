package configuration

import (
	"testing"
)

func TestParseConfig_sunny(t *testing.T) {

	conf := Configuration{}
	conf.items = map[string]Item{
		"wsLogEnable": Item{Key: "wsLogEnable", Value: "true"},
		"wsLogPort":   Item{Key: "wsLogPort", Value: "1234"},
		"wsLogLevel":  Item{Key: "wsLogLevel", Value: "fatal,error"},
	}

	// sunny day scenarios
	var wpw WPWithin
	wpw.ParseConfig(conf)
	if wpw.WSLogEnable == false {
		t.Error("WSLogEnable should be true")
		t.FailNow()
	}
	if wpw.WSLogPort != 1234 {
		t.Error("WSLogPort should be 1234")
		t.FailNow()
	}
	if wpw.WSLogLevel != "fatal,error" {
		t.Error("WSLogLevel should be info, but is: " + wpw.WSLogLevel)
		t.FailNow()
	}

	conf.items["wsLogHost"] = Item{Key: "wsLogHost", Value: "ALL"}
	wpw.ParseConfig(conf)
	if wpw.WSLogHost != "ALL" {
		t.Error("WSLogHost should be ALL, but is: " + wpw.WSLogHost)
		t.FailNow()
	}
}

func TestParseConfig_rainy1(t *testing.T) {
	// rainy day scenarios 1
	conf := Configuration{}
	conf.items = map[string]Item{
		"wsLogEnable": Item{Key: "wsLogEnable", Value: "bad value"},
		"wsLogPort":   Item{Key: "wsLogPort", Value: "bad value"},
		"wsLogLevel":  Item{Key: "wsLogLevel", Value: "info"},
	}
	var wpw2 WPWithin
	wpw2.ParseConfig(conf)

	if wpw2.WSLogEnable == true {
		t.Error("WSLogEnable should be false")
		t.FailNow()
	}
}
func TestParseConfig_rainy2(t *testing.T) {
	// rainy day scenarios 1
	conf := Configuration{}
	conf.items = map[string]Item{
		"wsLogEnable": Item{Key: "wsLogEnable", Value: "bad value"},
		"wsLogPort":   Item{Key: "wsLogPort", Value: "bad value"},
		"wsLogLevel":  Item{Key: "wsLogLevel", Value: "info"},
	}
	var wpw2 WPWithin
	wpw2.ParseConfig(conf)

	if wpw2.WSLogPort != 0 {
		t.Error("WSLogPort should be 5678")
		t.FailNow()
	}
	if wpw2.WSLogLevel != "info" {
		t.Error("WSLogLevel should be info, but is: " + wpw2.WSLogLevel)
	}
}
