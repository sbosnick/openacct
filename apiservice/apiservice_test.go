// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package apiservice

import (
	"net/http"
	"testing"

	"golang.org/x/net/context"

	"github.com/sbosnick1/openacct/domain"
	"github.com/stretchr/testify/assert"
)

type fakeStore struct {
	fundRepository domain.FundRepository
}

func (f *fakeStore) FundRepository() domain.FundRepository {
	return f.fundRepository
}

func TestNewApiListsAllFunds(t *testing.T) {
	assert := assert.New(t)
	request, responsewriter := getRequestResponse(t, "/v1/fund")
	fakerepository := newFakeFundRepository([]fakeFund{{1, domain.CAD, "General"},
		{2, domain.USD, "Special"}})
	fakestore := &fakeStore{fakerepository}

	sut := newApi(fakestore)
	sut.ServeHTTPC(context.Background(), responsewriter, request)

	assert.Equal(http.StatusOK, responsewriter.Code, "Unexpected status code.")
	assert.True(fakerepository.getAllCalled, "Listing the funds failed to call GetAll()")
}
