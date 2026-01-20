package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBridge_Http_Fetch(t *testing.T) {
	// 1. Start Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from Go Server")
	}))
	defer ts.Close()

	harness := NewHarness(t)

	// 2. JS Code to fetch from local server
	// We inject the server URL into the JS environment
	_ = harness.Engine.GlobalSet("TEST_URL", ts.URL)

	harness.Run(t, `
		const http = __go_http__;
		const url = TEST_URL;

		const resp = await http.Fetch(url);
		
		if (resp.StatusCode !== 200) {
			throw new Error("Expected status 200, got " + resp.StatusCode);
		}
		if (resp.Body !== "Hello from Go Server") {
			throw new Error("Unexpected body: " + resp.Body);
		}
	`)
}
