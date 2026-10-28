package httpapi

import (
	"github.com/jonatassales/tao/apps/extub/src/domain"
	"net/http"
)

func (handler Handler) products(request *http.Request, parts []string) (any, error) {
	method := request.Method
	if len(parts) == 0 {
		if method == http.MethodGet {
			catalog := handler.Service.Snapshot()
			runtimeProducts := map[string]string{}
			for _, runtime := range catalog.Runtimes {
				runtimeProducts[runtime.ID] = runtime.ProductID
			}
			return map[string]any{"source": "stub", "products": catalog.Products, "runtimeProducts": runtimeProducts}, nil
		}
		if method != http.MethodPost {
			return nil, errMethod
		}
		var input struct {
			Path string `json:"path"`
			Name string `json:"name"`
		}
		if err := decode(request, &input); err != nil {
			return nil, err
		}
		product, err := handler.Service.CreateProduct(input.Path, input.Name)
		return envelope("product", product), err
	}
	catalog := handler.Service.Snapshot()
	product, err := catalog.Product(parts[0])
	if err != nil {
		return nil, err
	}
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}
	switch action {
	case "executions", "history":
		if method != http.MethodGet || len(parts) != 2 {
			return nil, errMethod
		}
		return catalog.Activity(product.ID, action, request.URL.Query().Get("cursor"), handler.Service.Now())
	case "secrets":
		return handler.secrets(request, parts)
	case "episodes":
		if method != http.MethodGet {
			return nil, errMethod
		}
		return envelope("episodes", []any{}), nil
	case "recommendations":
		if method != http.MethodPost {
			return nil, errMethod
		}
		return recommendation(), nil
	case "prepare", "discover":
		if method != http.MethodPost {
			return nil, errMethod
		}
	case "":
		if method != http.MethodGet {
			return nil, errMethod
		}
	default:
		return nil, domain.ErrNotFound
	}
	unknown := 0
	for _, widget := range catalog.Widgets(handler.Service.Now(), 0) {
		if widget.Runtime.ProductID != product.ID {
			continue
		}
		for _, execution := range widget.Executions {
			if execution.State == "unknown" {
				unknown++
			}
		}
	}
	return map[string]any{"source": "stub", "product": product, "agent": map[string]any{"state": "connected", "lastSeenAt": handler.Service.Now().UTC()}, "unknownCount": unknown}, nil
}
func (handler Handler) secrets(request *http.Request, parts []string) (any, error) {
	product := parts[0]
	switch request.Method {
	case http.MethodGet:
		own := []domain.Secret{}
		for _, secret := range handler.Service.Snapshot().Secrets {
			if secret.ProductID == product {
				own = append(own, secret)
			}
		}
		return envelope("secrets", own), nil
	case http.MethodPost:
		if len(parts) != 2 {
			return nil, domain.ErrNotFound
		}
		var input struct {
			Name string `json:"name"`
		}
		if err := decode(request, &input); err != nil {
			return nil, err
		}
		secret, err := handler.Service.CreateSecret(product, input.Name)
		return map[string]any{"source": "stub", "credential": secret, "secret": "synthetic-" + secret.ID + "-not-a-credential"}, err
	case http.MethodDelete:
		if len(parts) != 3 {
			return nil, domain.ErrNotFound
		}
		return envelope("revoked", true), handler.Service.RevokeSecret(product, parts[2])
	default:
		return nil, errMethod
	}
}
func recommendation() any {
	return envelope("recommendation", map[string]any{"status": "suggested", "values": map[string]any{"cpu": 2, "memoryMiB": 2048, "model": "seed/model", "effort": "high", "llmBudgetUsd": 10}, "episodeIds": []string{}, "reasonCode": "stub_suggestion", "limitationCodes": []string{"future_workload_unverified"}})
}
