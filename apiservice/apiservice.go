// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

// Provides the JSON API service to expose the accounting domain
// business logic as a REST http service.
package apiservice

import (
	"net/http"

	"github.com/derekdowling/jsh-api"
	"github.com/sbosnick1/openacct/domain"
)

const apiV1Prefix = "/v1"

// New() is a factory for the api service to expose the provided
// domain.Store. The returned handler will service request for resources
// on a JSON API at URL's prefixed with "/v1".
func New(store domain.Store) http.Handler {
	return &rootAdaptor{newApi(store)}
}

func newApi(store domain.Store) *jshapi.API {
	api := jshapi.New(apiV1Prefix)
	api.Add(newFundResource(store.FundRepository()))
	return api
}
