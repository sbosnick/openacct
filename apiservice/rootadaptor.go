package apiservice

import (
	"net/http"

	"goji.io"
	"golang.org/x/net/context"
)

// rootAdaptor adapts a goji.Handler as an http.Handler using context.Background()
// as the context. rootAdaptor should be used to adapt the root handler of a hierarchry.
// The zero valued rootAdaptor responds to all requests with "status not found".
type rootAdaptor struct {
	root goji.Handler
}

func (r *rootAdaptor) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if r.root == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	r.root.ServeHTTPC(context.Background(), response, request)
}
