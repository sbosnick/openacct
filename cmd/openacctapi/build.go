// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

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

	handler := apiservice.New(store)

	return handler, nil
}
