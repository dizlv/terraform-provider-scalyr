package sdk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func ServeAndValidate(t *testing.T, method string, path string, contents string, validate func(r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("Wrong Method, %v", r.Method)
		}
		if r.URL.Path != path {
			t.Errorf("Wrong Path, %v", r.URL.Path)
		}
		fmt.Fprintln(w, contents)
		validate(r)
	}))
	return ts
}

func TestAPIResponseSuccess(t *testing.T) {
	success := &APIResponse{Status: "success"}
	if validateAPIResponse(success, "") != nil {
		t.Errorf("APIResponse.Status == success should not be an error")
	}
}

func TestAPIResponseFailed(t *testing.T) {
	success := &APIResponse{Status: "meh"}
	if validateAPIResponse(success, "") == nil {
		t.Errorf("APIResponse.Status != success should produce an error")
	}
}

func TestAPIResponseFailedNoMessage(t *testing.T) {
	success := &APIResponse{Status: "success", Message: "meh"}
	if validateAPIResponse(success, "") != nil {
		t.Errorf("APIResponse.Status == success should not be an error")
	}
}

type TestResponse struct {
	Hello string `json:"hello"`
}

func TestBasicJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{ \"hello\": \"yo\" }")
	}))
	defer ts.Close()
	config := &ScalyrConfig{Endpoint: ts.URL}
	tr := &TestResponse{}
	err := NewRequest("GET", "/meh", config).jsonResponse(tr)
	if tr.Hello != "yo" {
		t.Errorf("Basic JSON Response fail - %v", err)
	}
}

func TestMissingAuthJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{ \"hello\": \"yo\" }")
	}))
	defer ts.Close()
	config := &ScalyrConfig{Endpoint: ts.URL}
	tr := &TestResponse{}
	r := NewRequest("GET", "/meh", config).withWriteConfig().withReadConfig().withReadLog().withWriteLog()
	err := r.jsonResponse(tr)
	if err == nil {
		t.Errorf("Should of gotten an error about missing authentication")
	}
	expectedAuthMethods := []string{"WriteConfig", "ReadConfig", "ReadLog", "WriteLog"}
	if !reflect.DeepEqual(r.supportedKeys, expectedAuthMethods) {
		t.Errorf("Should of gotten %v but got %v", r.supportedKeys, expectedAuthMethods)
	}
}

func TestAuthOrderJSONResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{ \"hello\": \"yo\" }")
	}))
	defer ts.Close()
	config := &ScalyrConfig{Endpoint: ts.URL, Tokens: ScalyrTokens{WriteLog: "writeLog", ReadLog: "readLog", WriteConfig: "writeConfig", ReadConfig: "readConfig"}}
	tr := &TestResponse{}
	r := NewRequest("GET", "/meh", config).withWriteConfig().withReadConfig().withReadLog().withWriteLog()
	err := r.jsonResponse(tr)
	if err != nil {
		t.Errorf("Should not have gotten an error about missing authentication")
	}

	if r.apiKey != "writeConfig" {
		t.Errorf("WriteConfig API Key should have been used")
	}

	r = NewRequest("GET", "/meh", config).withReadConfig().withReadLog().withWriteLog()
	err = r.jsonResponse(tr)
	if r.apiKey != "readConfig" {
		t.Errorf("ReadConfig API Key should have been used")
	}

	r = NewRequest("GET", "/meh", config).withReadLog().withWriteLog()
	err = r.jsonResponse(tr)
	if r.apiKey != "readLog" {
		t.Errorf("ReadLog API Key should have been used")
	}

	r = NewRequest("GET", "/meh", config).withWriteLog()
	err = r.jsonResponse(tr)
	if r.apiKey != "writeLog" {
		t.Errorf("WriteLog API Key should have been used")
	}
}
