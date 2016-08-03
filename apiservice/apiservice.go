// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

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
