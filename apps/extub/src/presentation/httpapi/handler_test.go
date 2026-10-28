package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonatassales/tao/apps/extub/src/application"
	"github.com/jonatassales/tao/apps/extub/src/domain"
	"github.com/jonatassales/tao/apps/extub/src/infrastructure/seed"
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
