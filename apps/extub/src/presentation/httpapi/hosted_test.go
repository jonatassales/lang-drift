package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

const hostedToken = "a-long-shared-preview-token"

func hostedCall(t *testing.T, handler Handler, method, path, token string) (int, map[string]any) {
	t.Helper()
	request := httptest.NewRequest(method, "https://langdrift-stubs.vercel.app"+path, bytes.NewReader([]byte("{}")))
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, request)
	var result map[string]any
	if writer.Body.Len() > 0 {
		if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
	}
	return writer.Code, result
}

func errorCode(result map[string]any) string {
	failure, _ := result["error"].(map[string]any)
	code, _ := failure["code"].(string)
	return code
}

func TestLocalModeStillRefusesForeignHosts(t *testing.T) {
	handler, _ := fixture(t)
	code, result := hostedCall(t, handler, "GET", "/console/products", hostedToken)
	if code != 403 || errorCode(result) != "local_only" {
		t.Fatalf("local mode answered %d %v", code, result)
	}
}

func TestHostedModeAdmitsOnlyTheToken(t *testing.T) {
	handler, _ := fixture(t)
	handler.Token = hostedToken

	for _, token := range []string{"", "wrong-token", hostedToken + "x"} {
		code, result := hostedCall(t, handler, "GET", "/console/products", token)
		if code != 401 || errorCode(result) != "token_required" {
			t.Fatalf("token %q answered %d %v", token, code, result)
		}
	}

	code, result := hostedCall(t, handler, "GET", "/console/products", hostedToken)
	if code != 200 || len(result["products"].([]any)) == 0 {
		t.Fatalf("hosted read answered %d %v", code, result)
	}

	code, result = hostedCall(t, handler, "GET", "/health", "")
	if code != 200 || result["status"] != "ready" {
		t.Fatalf("health answered %d %v", code, result)
	}
	code, result = hostedCall(t, handler, "POST", "/health", "")
	if code != 401 || errorCode(result) != "token_required" {
		t.Fatalf("posting to health answered %d %v", code, result)
	}
}

func TestHostedModeIsReadOnlyWithoutAStream(t *testing.T) {
	handler, store := fixture(t)
	handler.Token = hostedToken
	runtime := store.Snapshot().Runtimes[0]

	for _, check := range []struct{ method, path string }{
		{"POST", "/console/products"},
		{"POST", "/console/runtimes/" + runtime.ID + "/start"},
		{"DELETE", "/console/runtimes/" + runtime.ID},
		{"POST", "/console/auth/login"},
		{"HEAD", "/console/products"},
		{"PATCH", "/console/runtimes/" + runtime.ID},
		{"OPTIONS", "/console/products"},
		{"PUT", "/webrtc/offer/x"},
	} {
		code, result := hostedCall(t, handler, check.method, check.path, hostedToken)
		if code != 403 || errorCode(result) != "read_only_preview" {
			t.Fatalf("%s %s answered %d %v", check.method, check.path, code, result)
		}
	}
	if store.Snapshot().Runtimes[0].DesiredState != runtime.DesiredState {
		t.Fatal("a hosted write changed the seed")
	}

	handler.Signaler = nil
	for _, path := range []string{"/webrtc/offer", "/webrtc/offer/"} {
		code, result := hostedCall(t, handler, "POST", path, hostedToken)
		if code != 503 || errorCode(result) != "stream_unavailable" {
			t.Fatalf("offer at %s answered %d %v", path, code, result)
		}
	}
}
