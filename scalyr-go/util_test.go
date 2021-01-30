package sdk

import (
	"encoding/json"
	"testing"
	"time"
)

type APITimeTest struct {
	TestTime APITime
}

func TestAPITime(t *testing.T) {
	now := time.Now()
	test := &APITimeTest{TestTime: APITime(now)}
	if !time.Time(test.TestTime).Equal(now) {
		t.Errorf("Expected our times to be the same! %v != %v", now, time.Time(test.TestTime))
	}
	err := json.Unmarshal([]byte("{ \"TestTime\": 0 }"), test)
	if err != nil {
		t.Errorf("Failed to Unmarshal a decent JSON payload - %v", err)
	}
	if !time.Time(test.TestTime).Equal(time.Unix(0, 0)) {
		t.Errorf("Expected Epoch value, got - %v", test.TestTime)
	}
	err = json.Unmarshal([]byte("{ \"TestTime\": 1000 }"), test)
	if !time.Time(test.TestTime).Equal(time.Unix(1, 0)) {
		t.Errorf("Expected Epoch value + 1000ms, got - %v", test.TestTime)
	}
}
