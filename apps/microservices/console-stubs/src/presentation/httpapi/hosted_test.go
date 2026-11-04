package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

const hostedToken = "a-long-shared-preview-token"

// hostedServe sends one request to the hosted preview with the given
// Authorization header, when there is one.
func hostedServe(t *testing.T, handler Handler, method, path, authorization string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	request := httptest.NewRequest(method, "https://langdrift-stubs.vercel.app"+path, bytes.NewReader([]byte("{}")))
	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}
	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, request)
	var result map[string]any
	if writer.Body.Len() > 0 {
		if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
	}
	return writer, result
}

func hostedCall(t *testing.T, handler Handler, method, path, token string) (int, map[string]any) {
	t.Helper()
	authorization := ""
	if token != "" {
		authorization = "Bearer " + token
	}
	writer, result := hostedServe(t, handler, method, path, authorization)
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

func TestHostedModeChallengesForABearerToken(t *testing.T) {
	handler, _ := fixture(t)
	handler.Token = hostedToken

	for _, authorization := range []string{
		"",
		"Bearer wrong-token",
		"Bearer " + hostedToken[:len(hostedToken)-1],
		"Bearer " + hostedToken + "x",
		"Basic " + hostedToken,
		"Bearer",
		"Bearer" + hostedToken,
		hostedToken,
	} {
		writer, result := hostedServe(t, handler, "GET", "/console/products", authorization)
		if writer.Code != 401 || errorCode(result) != "token_required" {
			t.Fatalf("%q answered %d %v", authorization, writer.Code, result)
		}
		if challenge := writer.Header().Get("WWW-Authenticate"); challenge != `Bearer realm="extub"` {
			t.Fatalf("%q challenged with %q", authorization, challenge)
		}
	}

	for _, authorization := range []string{
		"Bearer " + hostedToken,
		"bearer " + hostedToken,
		"BEARER " + hostedToken,
		"Bearer  " + hostedToken,
	} {
		writer, result := hostedServe(t, handler, "GET", "/console/products", authorization)
		if writer.Code != 200 || len(result["products"].([]any)) == 0 {
			t.Fatalf("%q answered %d %v", authorization, writer.Code, result)
		}
	}
}

func TestEveryHostedAnswerCarriesTheAPIHeaders(t *testing.T) {
	handler, _ := fixture(t)
	handler.Token = hostedToken
	authorized := "Bearer " + hostedToken
	for _, check := range []struct {
		method, path, authorization string
		status                      int
	}{
		{"GET", "/health", "", 200},
		{"GET", "/console/products", "", 401},
		{"GET", "/console/products", authorized, 200},
		{"GET", "/console/missing", authorized, 404},
		{"POST", "/console/products", authorized, 403},
		{"POST", "/webrtc/offer", authorized, 503},
	} {
		label := check.method + " " + check.path
		writer, _ := hostedServe(t, handler, check.method, check.path, check.authorization)
		if writer.Code != check.status {
			t.Fatalf("%s answered %d, want %d", label, writer.Code, check.status)
		}
		requireAPIHeaders(t, label, writer.Header())
	}
}
