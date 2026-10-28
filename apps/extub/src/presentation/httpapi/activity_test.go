package httpapi

import (
	"fmt"
	"testing"
)

func TestActivityPagesAndGeneratedSecrets(t *testing.T) {
	handler, store := fixture(t)
	runtime := store.Snapshot().Runtimes[0]
	for i := 0; i < 60; i++ {
		if err := handler.Service.Operate(runtime.ID, "stop"); err != nil {
			t.Fatal(err)
		}
		if err := handler.Service.Operate(runtime.ID, "start"); err != nil {
			t.Fatal(err)
		}
	}
	for _, kind := range []string{"executions", "history"} {
		route := "/console/products/" + runtime.ProductID + "/" + kind
		code, first := call(t, handler, "GET", route, nil)
		if code != 200 || len(first["items"].([]any)) != 50 || first["nextCursor"] == nil {
			t.Fatalf("invalid first page: %v", first)
		}
		cursor := first["nextCursor"].(string)
		if err := handler.Service.Operate(runtime.ID, "stop"); err != nil {
			t.Fatal(err)
		}
		_, second := call(t, handler, "GET", route+"?cursor="+cursor, nil)
		seen := map[string]bool{}
		for _, page := range []map[string]any{first, second} {
			for _, raw := range page["items"].([]any) {
				item := raw.(map[string]any)
				key := fmt.Sprint(item["position"])
				if seen[key] {
					t.Fatal("duplicate cursor row")
				}
				seen[key] = true
			}
		}
		code, _ = call(t, handler, "GET", route+"?cursor=invalid", nil)
		if code != 400 {
			t.Fatal("invalid cursor accepted")
		}
	}
	code, secret := call(t, handler, "POST", "/console/products/"+runtime.ProductID+"/secrets", map[string]any{})
	if code != 200 || secret["secret"] == "" {
		t.Fatal("secret not generated")
	}
	other, err := handler.Service.CreateProduct("/seeds/other", "Other")
	if err != nil {
		t.Fatal(err)
	}
	_, listed := call(t, handler, "GET", "/console/products/"+other.ID+"/history", nil)
	if len(listed["items"].([]any)) != 0 {
		t.Fatal("product data leaked")
	}
}
