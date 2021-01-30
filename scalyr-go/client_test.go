package sdk

import (
	"os"
	"testing"
)

func init() {
	os.Clearenv()
}

func TestHasTeam(t *testing.T) {
	sc := &ScalyrConfig{Team: "meh"}
	if !sc.hasTeam() {
		t.Errorf("Team should have been true")
	}
	sc = &ScalyrConfig{}
	if sc.hasTeam() {
		t.Errorf("Team should not have been true")
	}
}

func TestNewClient(t *testing.T) {

	sc, err := NewClient(&ScalyrConfig{})
	if err != nil {
		t.Errorf("Null Config should of worked")
	}
	if sc.Region != "us" {
		t.Errorf("Null Config should default to US")
	}
	if sc.Endpoint != "https://app.scalyr.com/" {
		t.Errorf("Null Config should default to https://app.scalyr.com/ found %v", sc.Endpoint)
	}

	os.Setenv("SCALYR_SERVER", "http://test")
	sc, err = NewClient(&ScalyrConfig{})
	if sc.Endpoint != "http://test" {
		t.Errorf("Setting SCALYR_SERVER should work, expected http://test but got %v", sc.Endpoint)
	}

	os.Setenv("SCALYR_SERVER", "test")
	sc, err = NewClient(&ScalyrConfig{})
	if sc.Endpoint != "https://test" {
		t.Errorf("Setting SCALYR_SERVER should work, expected https://test but got %v", sc.Endpoint)
	}

	os.Setenv("SCALYR_WRITELOG_TOKEN", "writelog")
	os.Setenv("SCALYR_READLOG_TOKEN", "readlog")
	os.Setenv("SCALYR_WRITECONFIG_TOKEN", "writeconfig")
	os.Setenv("SCALYR_READCONFIG_TOKEN", "readconfig")
	sc, err = NewClient(&ScalyrConfig{})
	if sc.Tokens.ReadLog != "readlog" {
		t.Errorf("ReadLog token should be readlog got %v", sc.Tokens.ReadLog)
	}
	if sc.Tokens.WriteLog != "writelog" {
		t.Errorf("WriteLog token should be writelog got %v", sc.Tokens.WriteLog)
	}
	if sc.Tokens.ReadConfig != "readconfig" {
		t.Errorf("ReadConfig token should be readconfig got %v", sc.Tokens.ReadConfig)
	}
	if sc.Tokens.WriteConfig != "writeconfig" {
		t.Errorf("WriteConfig token should be writeconfig got %v", sc.Tokens.WriteConfig)
	}
}
