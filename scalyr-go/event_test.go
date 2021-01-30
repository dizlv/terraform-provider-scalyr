package sdk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

const sendEventResponse = `{
  "status": "success"
}`

func TestSendEvent(t *testing.T) {
	ts := ServeAndValidate(t, "POST", "/api/addEvents", sendEventResponse, func(r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		cer := &CreateEventsRequest{}
		err := json.Unmarshal(data, cer)
		if err != nil {
			t.Errorf("Error reading request: %v", err)
		}
		if cer.Session != "a" {
			t.Errorf("Session should of been a instead of %v", cer.Session)
		}
		if cer.SessionInfo.ServerType != "b" {
			t.Errorf("ServerType should of been b instead of %v", cer.SessionInfo.ServerType)
		}
		if cer.SessionInfo.ServerID != "a" {
			t.Errorf("ServerID should of been a instead of %v", cer.SessionInfo.ServerID)
		}
		if (*cer.Events)[0].Attrs["message"] != "test" {
			t.Errorf("event[0].messageshould of been test instead of %v", (*cer.Events)[0].Attrs["message"])
		}
	})
	defer ts.Close()
	config := &ScalyrConfig{Endpoint: ts.URL}
	sc, _ := NewClient(config)

	event := &Event{Thread: "5", Sev: 3, Ts: "0", Attrs: map[string]interface{}{"message": "test", "meh": 1}}
	err := sc.SendEvent(event, &Thread{ID: 5, Name: "fred"}, "a", &SessionInfo{ServerID: "a", ServerType: "b"})
	if err != nil {
		t.Errorf("Unexpected Error got %v", err)
	}
}
