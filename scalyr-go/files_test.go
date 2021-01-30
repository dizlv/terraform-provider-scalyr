package sdk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const getFileResponse = `{
  "path": "/test/test",
  "content": "meh",
  "version": 55,
  "createDate": 1000,
  "modDate": 2000
}`

const putFileResponse = `{
  "path": "/test/test"
}`

func TestGetFile(t *testing.T) {
	ts := ServeAndValidate(t, "POST", "/api/getFile", getFileResponse, func(r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		gfr := &GetFileRequest{}
		err := json.Unmarshal(data, gfr)
		if err != nil {
			t.Errorf("Error reading request: %v", err)
		}
		if gfr.Path != "/test/test" {
			t.Errorf("Incorrect Path. Expected /meh, got %v", gfr.Path)
		}
	})
	defer ts.Close()
	config := &ScalyrConfig{Endpoint: ts.URL}
	sc, _ := NewClient(config)
	gfr, err := sc.GetFile("/test/test")
	if err != nil {
		t.Errorf("Unexpected Error got %v", err)
	}
	if gfr.Content != "meh" {
		t.Errorf("File Contents Invalid - expected meh, got %v", gfr.Content)
	}
	if gfr.Version != 55 {
		t.Errorf("File Version should be 55, got %v", gfr.Version)
	}
	if gfr.Path != "/test/test" {
		t.Errorf("File Path should be /meh, got %v", gfr.Path)
	}
	if !time.Unix(1, 0).Equal(time.Time(gfr.CreateDate)) {
		t.Errorf("CreateDate Wrong, got %v", gfr.CreateDate)
	}
	if !time.Unix(2, 0).Equal(time.Time(gfr.ModDate)) {
		t.Errorf("ModDate Wrong, got %v", gfr.ModDate)
	}
}

func TestPutFile(t *testing.T) {
	ts := ServeAndValidate(t, "POST", "/api/putFile", putFileResponse, func(r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		pfr := &PutFileRequest{}
		err := json.Unmarshal(data, pfr)
		if err != nil {
			t.Errorf("Error reading request: %v", err)
		}
		if pfr.Path != "/test/test" {
			t.Errorf("Incorrect Path. Expected /test/test, got %v", pfr.Path)
		}
		if pfr.Content != "rar" {
			t.Errorf("Wrong. Expected rar, got %v", pfr.Content)
		}
	})
	defer ts.Close()
	config := &ScalyrConfig{Endpoint: ts.URL}
	sc, _ := NewClient(config)
	pfr, err := sc.PutFile("/test/test", "rar")
	if err != nil {
		t.Errorf("Unexpected Error in PutFile - %v", err)
	}
	if pfr.Path != "/test/test" {
		t.Errorf("Unexpected Path in PutFileResponse - %v", pfr.Path)
	}

}
