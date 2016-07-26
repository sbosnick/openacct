// Provides the JSON API service to expose the accounting domain
// business logic as a REST http service.
package apiservice

import (
	"github.com/sbosnick1/openacct/domain"
	"net/http"
)

func New(store domain.Store) (http.Handler, error) {
	return http.NotFoundHandler(), nil
}
