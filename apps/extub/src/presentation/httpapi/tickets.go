package httpapi

import (
	"net/http"

	"github.com/jonatassales/tao/apps/extub/src/domain"
)

// tickets serves one generated ticket by id.
func (handler Handler) tickets(request *http.Request, parts []string) (any, error) {
	if request.Method != http.MethodGet {
		return nil, errMethod
	}
	if len(parts) != 1 {
		return nil, domain.ErrNotFound
	}
	ticket, err := handler.Service.Snapshot().Ticket(parts[0])
	if err != nil {
		return nil, err
	}
	return envelope("ticket", ticket), nil
}
