package openacctapi

import (
	"github.com/sbosnick1/openacct/apiservice"
	"github.com/sbosnick1/openacct/domain"
	"net/http"
)

func BuildApiHandler(dsn string) (http.Handler, error) {
	store, err := domain.New(dsn)
	if err != nil {
		return nil, err
	}

	handler, err := apiservice.New(store)
	if err != nil {
		return nil, err
	}

	return handler, nil
}
