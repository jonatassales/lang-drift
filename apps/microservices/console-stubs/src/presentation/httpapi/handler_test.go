package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/application"
	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/domain"
	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/infrastructure/seed"
)

type signalFake struct{ calls int }

func (fake *signalFake) Offer(context.Context, application.Description) (application.Description, error) {
	fake.calls++
	return application.Description{}, domain.ErrInvalid
}
func fixture(t *testing.T) (Handler, *seed.Store) {
	t.Helper()
	store, err := seed.New(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return Handler{Service: application.Service{Store: store, Now: time.Now}, Signaler: &signalFake{}}, store
}
func call(t *testing.T, handler Handler, method, path string, body any) (int, map[string]any) {
	t.Helper()
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	request := httptest.NewRequest(method, "http://127.0.0.1:3401"+path, bytes.NewReader(payload))
	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, request)
	var result map[string]any
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["source"] != "stub" {
		t.Fatalf("missing synthetic origin: %v", result)
	}
	return writer.Code, result
}

func TestPerRuntimeLifecycleFences(t *testing.T) {
	handler, store := fixture(t)
	first := store.Snapshot().Runtimes[0]
	detail := "/console/runtimes/" + first.ID
	for _, check := range []struct {
		method, path string
		body         any
		status       int
	}{
		{"PATCH", detail, first, 409}, {"GET", detail + "/stop", nil, 405}, {"POST", detail + "/stop", map[string]any{}, 200}, {"PATCH", detail, first, 200}, {"POST", detail + "/revoke", map[string]any{}, 200}, {"POST", detail + "/start", map[string]any{}, 409}, {"POST", "/seeds/empty", map[string]any{}, 404},
	} {
		status, _ := call(t, handler, check.method, check.path, check.body)
		if status != check.status {
			t.Fatalf("%s %s: status %d", check.method, check.path, status)
		}
	}
	_, home := call(t, handler, "GET", "/console/home", nil)
	if len(home["products"].([]any)) != 1 || len(home["running"].([]any)) != 3 {
		t.Fatal("an individual action changed unrelated fixtures")
	}
	_, created := call(t, handler, "POST", "/console/products", map[string]string{"path": "/seeds/new", "name": "Empty"})
	_, duplicate := call(t, handler, "POST", "/console/products", map[string]string{"path": "/seeds/./new"})
	if created["product"].(map[string]any)["id"] != duplicate["product"].(map[string]any)["id"] {
		t.Fatal("duplicate folder identity")
	}
	if len(store.Snapshot().Runtimes) != 4 {
		t.Fatal("registration created a runtime")
	}
	first.ProductID = created["product"].(map[string]any)["id"].(string)
	status, saved := call(t, handler, "POST", "/console/runtimes", first)
	if status != 200 || saved["runtime"].(map[string]any)["desiredState"] != "stopped" {
		t.Fatal("creation started a fixture")
	}
}

// createOtherProduct registers a second product through the API, so a test
// can prove isolation without depending on a seeded product.
func createOtherProduct(t *testing.T, handler Handler) string {
	t.Helper()
	status, created := call(t, handler, "POST", "/console/products", map[string]string{"path": "/seeds/other", "name": "Other"})
	if status != 200 {
		t.Fatalf("creating a product: status %d", status)
	}

	return created["product"].(map[string]any)["id"].(string)
}

func TestSeedsOneProjectAndItsRuntimes(t *testing.T) {
	handler, _ := fixture(t)
	_, listedProducts := call(t, handler, "GET", "/console/products", nil)
	products := listedProducts["products"].([]any)
	if len(products) != 1 || products[0].(map[string]any)["id"] != "seed-microruntimes" {
		t.Fatalf("seeded products: %v", products)
	}
	_, listedRuntimes := call(t, handler, "GET", "/console/runtimes?productId=seed-microruntimes", nil)
	if len(listedRuntimes["runtimes"].([]any)) != 4 {
		t.Fatalf("seeded runtimes: %v", listedRuntimes["runtimes"])
	}
	_, listedFolders := call(t, handler, "GET", "/console/local/folders", nil)
	folders := listedFolders["folders"].([]any)
	if len(folders) != 1 || folders[0].(map[string]any)["path"] != "/seeds/microruntimes" {
		t.Fatalf("offered folders: %v", folders)
	}
}

func TestHTTPBoundaryAndProductIsolation(t *testing.T) {
	handler, store := fixture(t)
	for _, check := range []struct{ host, origin string }{{"attacker.test", ""}, {"127.0.0.1:3401", "https://foreign.example"}} {
		request := httptest.NewRequest("GET", "http://127.0.0.1:3401/health", nil)
		request.Host = check.host
		request.Header.Set("Origin", check.origin)
		writer := httptest.NewRecorder()
		handler.ServeHTTP(writer, request)
		if writer.Code != 403 {
			t.Fatal("accepted foreign request")
		}
	}
	otherProductID := createOtherProduct(t, handler)
	_, list := call(t, handler, "GET", "/console/runtimes?productId="+otherProductID, nil)
	if len(list["widgets"].([]any)) != 0 {
		t.Fatal("cross-product widgets")
	}
	_, secret := call(t, handler, "POST", "/console/products/seed-microruntimes/secrets", map[string]string{"name": "Fixture"})
	id := secret["credential"].(map[string]any)["id"].(string)
	if status, _ := call(t, handler, "DELETE", "/console/products/"+otherProductID+"/secrets/"+id, nil); status != 404 {
		t.Fatal("cross-product revocation")
	}
	if status, _ := call(t, handler, "POST", "/console/runtimes", map[string]any{}); status != 400 {
		t.Fatal("invalid runtime accepted")
	}
	before := store.Snapshot()
	before.Runtimes[0].Name = "mutated"
	if store.Snapshot().Runtimes[0].Name == "mutated" {
		t.Fatal("read aliases store state")
	}
	for i := 0; i < 20; i++ {
		_, _ = handler.Service.CreateProduct("/seeds/"+time.Now().String(), "Bounded")
	}
	if len(store.Snapshot().Products) > 16 {
		t.Fatal("unbounded products")
	}
}

// requireAPIHeaders fails unless an answer carries every header that keeps
// its JSON data for the console alone: never stored, sniffed, rendered,
// framed, sent as a referrer or read by another origin.
func requireAPIHeaders(t *testing.T, label string, header http.Header) {
	t.Helper()
	for name, value := range map[string]string{
		"Content-Type":                 "application/json",
		"Cache-Control":                "no-store",
		"X-Content-Type-Options":       "nosniff",
		"Content-Security-Policy":      "default-src 'none'; frame-ancestors 'none'",
		"X-Frame-Options":              "DENY",
		"Referrer-Policy":              "no-referrer",
		"Cross-Origin-Resource-Policy": "same-origin",
	} {
		if got := header.Get(name); got != value {
			t.Fatalf("%s: %s is %q, want %q", label, name, got, value)
		}
	}
}

// localServe sends one request to console-stubs on loopback and decodes its answer,
// when it has one.
func localServe(t *testing.T, handler Handler, request *http.Request) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
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

// oversized is a product body past the 64 KiB limit.
var oversized = `{"path":"/` + strings.Repeat("a", 65536) + `"}`

// padded is a valid product body of exactly size bytes.
func padded(size int) string {
	body := `{"path":"/seeds/padded","name":"Padded"}`
	return body + strings.Repeat(" ", size-len(body))
}

// chunked hides a body's length, as a chunked request does.
func chunked(body string) io.Reader {
	return struct{ io.Reader }{strings.NewReader(body)}
}

func TestEveryLocalAnswerCarriesTheAPIHeaders(t *testing.T) {
	handler, store := fixture(t)
	runtime := store.Snapshot().Runtimes[0]
	for _, check := range []struct {
		method, path, origin, body string
		status                     int
	}{
		{"GET", "/health", "", "", 200},
		{"GET", "/probe", "", "", 204},
		{"GET", "/console/products", "https://foreign.example", "", 403},
		{"GET", "/console/missing", "", "", 404},
		{"GET", "/console/runtimes/" + runtime.ID + "/stop", "", "", 405},
		{"POST", "/console/products", "", "{", 400},
		{"POST", "/console/products", "", oversized, 413},
	} {
		request := httptest.NewRequest(check.method, "http://127.0.0.1:3401"+check.path, strings.NewReader(check.body))
		if check.origin != "" {
			request.Header.Set("Origin", check.origin)
		}
		label := fmt.Sprintf("%s %s", check.method, check.path)
		writer, _ := localServe(t, handler, request)
		if writer.Code != check.status {
			t.Fatalf("%s answered %d, want %d", label, writer.Code, check.status)
		}
		requireAPIHeaders(t, label, writer.Header())
	}
}

func TestLocalModeRefusesCrossSiteCalls(t *testing.T) {
	handler, _ := fixture(t)
	for _, check := range []struct {
		host, origin, site string
		status             int
	}{
		{"127.0.0.1:3401", "", "cross-site", 403},
		{"127.0.0.1:3401", "https://foreign.example", "", 403},
		{"127.0.0.1:3401", "http://localhost:3401", "", 403},
		{"127.0.0.1:3401", "null", "", 403},
		{"127.0.0.1:3401", "http://127.0.0.1:3401", "same-origin", 200},
		{"localhost:3401", "http://localhost:3401", "same-origin", 200},
		{"127.0.0.1:3401", "", "", 200},
	} {
		request := httptest.NewRequest("GET", "http://"+check.host+"/console/products", nil)
		if check.origin != "" {
			request.Header.Set("Origin", check.origin)
		}
		if check.site != "" {
			request.Header.Set("Sec-Fetch-Site", check.site)
		}
		writer, result := localServe(t, handler, request)
		if writer.Code != check.status {
			t.Fatalf("%+v answered %d %v", check, writer.Code, result)
		}
		if check.status == 403 && errorCode(result) != "local_only" {
			t.Fatalf("%+v refused with %v", check, result)
		}
	}
}

func TestBodiesOverTheLimitAnswer413(t *testing.T) {
	handler, store := fixture(t)
	runtime := store.Snapshot().Runtimes[0]
	for _, check := range []struct {
		label, path string
		body        io.Reader
		status      int
	}{
		{"a JSON body over the limit", "/console/products", strings.NewReader(oversized), 413},
		{"a body over the limit that is not JSON", "/console/products", strings.NewReader(strings.Repeat("a", 70<<10)), 413},
		{"a body over the limit on a route that reads none", "/console/runtimes/" + runtime.ID + "/stop", strings.NewReader(oversized), 413},
		{"a chunked product over the limit", "/console/products", chunked(oversized), 413},
		{"a chunked runtime over the limit", "/console/runtimes", chunked(oversized), 413},
		{"a chunked load test over the limit", "/console/runtimes/" + runtime.ID + "/load-test", chunked(oversized), 413},
		{"a body one byte over the limit", "/console/products", strings.NewReader(padded(65537)), 413},
		{"a body of exactly 64 KiB", "/console/products", strings.NewReader(padded(65536)), 200},
	} {
		request := httptest.NewRequest("POST", "http://127.0.0.1:3401"+check.path, check.body)
		writer, result := localServe(t, handler, request)
		if writer.Code != check.status {
			t.Fatalf("%s answered %d %v", check.label, writer.Code, result)
		}
		if check.status == 413 && errorCode(result) != "request_too_large" {
			t.Fatalf("%s refused with %v", check.label, result)
		}
	}
}

func TestSeedCapsAnswer429(t *testing.T) {
	handler, store := fixture(t)
	catalog := store.Snapshot()
	product := catalog.Products[0].ID
	runtime := catalog.Runtimes[1]
	for _, check := range []struct {
		label, path string
		room        int
		body        func(int) any
	}{
		{"product", "/console/products", 16 - len(catalog.Products), func(index int) any {
			return map[string]string{"path": fmt.Sprintf("/seeds/capped-%d", index)}
		}},
		{"runtime", "/console/runtimes", 16 - len(catalog.Runtimes), func(int) any { return runtime }},
		{"secret", "/console/products/" + product + "/secrets", 32 - len(catalog.Secrets), func(int) any { return map[string]string{} }},
	} {
		for index := 0; index < check.room; index++ {
			if status, result := call(t, handler, "POST", check.path, check.body(index)); status != 200 {
				t.Fatalf("%s %d answered %d %v", check.label, index, status, result)
			}
		}
		status, result := call(t, handler, "POST", check.path, check.body(check.room))
		if status != 429 || errorCode(result) != "seed_limit" {
			t.Fatalf("the %s past the cap answered %d %v", check.label, status, result)
		}
	}
}

func TestFolderAndRuntimeBoundsAnswer400(t *testing.T) {
	handler, store := fixture(t)
	for _, folder := range []string{"seeds", "%00", "/seeds%00/microruntimes", "/" + strings.Repeat("a", 1024)} {
		status, result := call(t, handler, "GET", "/console/local/folders?path="+folder, nil)
		if status != 400 || errorCode(result) != "invalid_input" {
			t.Fatalf("folder %.40q answered %d %v", folder, status, result)
		}
	}
	longest := "/" + strings.Repeat("a", 1023)
	if status, listed := call(t, handler, "GET", "/console/local/folders?path="+longest, nil); status != 200 || listed["path"] != longest {
		t.Fatalf("a 1024-byte folder answered %d", status)
	}

	for label, change := range map[string]func(*domain.Component){
		"a negative port":                 func(component *domain.Component) { component.Port = -1 },
		"a port over 65535":               func(component *domain.Component) { component.Port = 65536 },
		"an image over 256 bytes":         func(component *domain.Component) { component.Image = strings.Repeat("i", 257) },
		"a relative readiness path":       func(component *domain.Component) { component.ReadinessPath = "health" },
		"a readiness path over 256 bytes": func(component *domain.Component) { component.ReadinessPath = "/" + strings.Repeat("r", 256) },
	} {
		runtime := store.Snapshot().Runtimes[1]
		change(&runtime.Components[0])
		for _, route := range []struct{ method, path string }{
			{"POST", "/console/runtimes"},
			{"PATCH", "/console/runtimes/" + runtime.ID},
		} {
			status, result := call(t, handler, route.method, route.path, runtime)
			if status != 400 || errorCode(result) != "invalid_input" {
				t.Fatalf("%s %s answered %d %v", route.method, label, status, result)
			}
		}
	}
}
