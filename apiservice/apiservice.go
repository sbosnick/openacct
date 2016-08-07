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

func New(store domain.Store) (http.Handler, error) {
	return http.NotFoundHandler(), nil
}

func newApi(store domain.Store) *jshapi.API {
	api := jshapi.New(apiV1Prefix)
	api.Add(newFundResource(store.FundRepository()))
	return api
}
