package openacctapi

import (
	"net/http"
)

func BuildApiHandler() (http.Handler, error) {
	return http.NotFoundHandler(), nil
}
