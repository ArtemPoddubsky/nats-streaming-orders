package app_test

import (
	"net/http"
	"testing"
)

func TestApp_RenderPage(t *testing.T) {

	tests := []struct {
		about    string
		reqID    string
		respCode int
	}{
		{
			about:    "TEST 1: Example orderUID.",
			reqID:    "b563feb7b2b84b6test",
			respCode: 200,
		},
		{
			about:    "TEST 2: Non-existing orderUID.",
			reqID:    "aaaaaaaaaaaaaa",
			respCode: 404,
		},
	}

	for idx, val := range tests {
		resp, err := http.Get("http://localhost:5000/" + val.reqID)

		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != val.respCode {
			t.Error("Test", idx+1, "FAIL", "\nExpected:", val.respCode, "Got:", resp.StatusCode)
		}

		resp.Body.Close()
	}
}
