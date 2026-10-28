package httpapi

import (
	"github.com/jonatassales/tao/apps/extub/src/domain"
	"net/http"
)

func (handler Handler) runtimes(request *http.Request, parts []string) (any, error) {
	method := request.Method
	if len(parts) == 0 {
		if method == http.MethodPost {
			var input domain.Runtime
			if err := decode(request, &input); err != nil {
				return nil, err
			}
			runtime, err := handler.Service.SaveRuntime("", input)
			return envelope("runtime", runtime), err
		}
		if method != http.MethodGet {
			return nil, errMethod
		}
		widgets := []domain.Widget{}
		runtimes := []domain.Runtime{}
		product := request.URL.Query().Get("productId")
		for _, widget := range handler.Service.Widgets(0) {
			if product != "" && widget.Runtime.ProductID != product {
				continue
			}
			widgets = append(widgets, widget)
			runtimes = append(runtimes, widget.Runtime)
		}
		return map[string]any{"source": "stub", "widgets": widgets, "runtimes": runtimes}, nil
	}
	id := parts[0]
	catalog := handler.Service.Snapshot()
	if catalog.RuntimeIndex(id) < 0 {
		return nil, domain.ErrNotFound
	}
	if len(parts) > 1 {
		if len(parts) != 2 {
			return nil, domain.ErrNotFound
		}
		action := parts[1]
		if action == "load-test" && method == http.MethodPost {
			var input map[string]any
			if err := decode(request, &input); err != nil || len(input) != 0 {
				return nil, domain.ErrInvalid
			}
			result, err := handler.Service.RunLoadTest(request.Context(), id)
			return envelope("loadTest", result), err
		}
		if action == "history" && method == http.MethodGet {
			events := []domain.Event{}
			for index := len(catalog.EventRecords) - 1; index >= 0; index-- {
				event := catalog.EventRecords[index].Event
				if event.RuntimeID == id {
					events = append(events, event)
				}
			}
			return map[string]any{"source": "stub", "events": events, "episodes": []any{}}, nil
		}
		if method != http.MethodPost {
			return nil, errMethod
		}
		if action == "recommendations" {
			return recommendation(), nil
		}
		if action != "start" && action != "stop" && action != "revoke" {
			return nil, domain.ErrNotFound
		}
		if err := handler.Service.Operate(id, action); err != nil {
			return nil, err
		}
	} else {
		switch method {
		case http.MethodGet:
		case http.MethodDelete:
			if err := handler.Service.Operate(id, "delete"); err != nil {
				return nil, err
			}
		case http.MethodPatch:
			var input domain.Runtime
			if err := decode(request, &input); err != nil {
				return nil, err
			}
			if _, err := handler.Service.SaveRuntime(id, input); err != nil {
				return nil, err
			}
		default:
			return nil, errMethod
		}
	}
	for _, widget := range handler.Service.Widgets(0) {
		if widget.Runtime.ID == id {
			if len(parts) == 1 && method == http.MethodGet {
				return map[string]any{"source": "stub", "runtime": widget.Runtime, "environments": widget.Environments, "executions": widget.Executions, "widget": widget}, nil
			}
			return widget, nil
		}
	}
	return nil, domain.ErrNotFound
}
