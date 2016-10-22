// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package apiservice

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	jsh "github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api"
	"github.com/sbosnick1/openacct/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

const (
	StatusUnprocessableEntity int = 422
)

type fakeFund struct {
	id       uint
	currency domain.Currency
	name     string
}

func (f *fakeFund) Currency() domain.Currency {
	return f.currency
}

func (f *fakeFund) Name() string {
	return f.name
}

func (f *fakeFund) Id() uint {
	return f.id
}

type fakeFundRepository struct {
	getAllCalled     bool
	createFundCalled bool
	nextID           uint
	funds            []domain.Fund
}

func (f *fakeFundRepository) GetAll() ([]domain.Fund, error) {
	f.getAllCalled = true
	return f.funds, nil
}

func (f *fakeFundRepository) Create(name string, currency domain.Currency) (domain.Fund, error) {
	f.createFundCalled = true
	f.nextID++
	fund := fakeFund{f.nextID, currency, name}
	f.funds = append(f.funds, &fund)
	return &fund, nil
}

func newFakeFundRepository(funds []fakeFund) *fakeFundRepository {
	var realfunds []domain.Fund
	for _, f := range funds {
		realfunds = append(realfunds, &f)
	}
	return &fakeFundRepository{funds: realfunds}
}

func newFundObject(t *testing.T, id string, attributes map[string]string) *jsh.Object {
	obj, jsherr := jsh.NewObject(id, "fund", attributes)
	if jsherr != nil {
		t.Fatal(jsherr)
	}
	return obj
}

func parseResponseBody(t *testing.T, payload *httptest.ResponseRecorder, mode jsh.DocumentMode) *jsh.Document {
	parser := jsh.Parser{
		Method:  "",
		Headers: payload.Header(),
	}

	doc, jsherr := parser.Document(ioutil.NopCloser(payload.Body), mode)
	require.Nil(t, jsherr, "jsh.Parser failed to parse the recorded response")
	return doc
}

func TestZeroFundStoreListsWithISE(t *testing.T) {
	var sut fundStore
	_, err := sut.List(context.Background())

	if err.StatusCode() != http.StatusInternalServerError {
		t.Errorf("zero fundStore gave error status %s when %s was expected",
			http.StatusText(err.StatusCode()), http.StatusText(http.StatusInternalServerError))
	}
}

func TestFundStoreListGetsAllFundsFromDomain(t *testing.T) {
	assert := assert.New(t)
	var rep fakeFundRepository

	sut := fundStore{&rep}
	_, err := sut.List(context.Background())

	assert.NoError(err, "Unexpected error when listing funds.")
	assert.True(rep.getAllCalled, "GetAll() was not called in the fund repository.")
}

func TestFundStoreListReturnsFundsFromDomain(t *testing.T) {
	assert := assert.New(t)
	rep := newFakeFundRepository([]fakeFund{{1, domain.CAD, "General"},
		{2, domain.USD, "Special"}})

	sut := fundStore{rep}
	list, err := sut.List(context.Background())

	require.NoError(t, err, "Unexpected error when listing funds.")
	require.Len(t, list, 2, "Unexpected number of funds returned.")
	assert.Equal("fund", list[0].Type, "Unexpected type of object returned.")
	assert.Equal("fund", list[1].Type, "Unexpected type of object returned.")

	for _, obj := range list {
		switch obj.ID {
		case "1":
			assert.JSONEq(`{"currency" : "CAD", "name" : "General"}`,
				string(obj.Attributes),
				"Unexpected attributes on a returned object.")
		case "2":
			assert.JSONEq(`{"currency" : "USD", "name" : "Special"}`,
				string(obj.Attributes),
				"Unexpected attributes on a returned object.")
		}
	}
}

func TestZeroFundStoreSavesWithISE(t *testing.T) {
	obj := newFundObject(t, "", map[string]string{})

	var sut fundStore
	_, jsherr := sut.Save(context.Background(), obj)

	assert.Equal(t, http.StatusInternalServerError, jsherr.StatusCode(),
		"zero fundStore gave unexpected status on Save()")
}

func TestFundStoreSaveWithoutNameIsError(t *testing.T) {
	rep := newFakeFundRepository([]fakeFund{})
	obj := newFundObject(t, "", map[string]string{"currency": "CAD"})

	sut := fundStore{rep}
	_, jsherr := sut.Save(context.Background(), obj)

	assert.Equal(t, StatusUnprocessableEntity, jsherr.StatusCode(),
		"fundStore gave unexpect status on Save()")
}

func TestFundStoreSaveWithInvalidNameIsError(t *testing.T) {
	badnames := []string{"#$%", ""}

	for _, badname := range badnames {
		rep := newFakeFundRepository([]fakeFund{})
		obj := newFundObject(t, "", map[string]string{"currency": "CAD", "name": badname})

		sut := fundStore{rep}
		_, jsherr := sut.Save(context.Background(), obj)

		assert.Equal(t, StatusUnprocessableEntity, jsherr.StatusCode(),
			"fundStore gave unexpect status on Save()")
	}
}

func TestFundStoreSaveWithoutCurrencyIsError(t *testing.T) {
	rep := newFakeFundRepository([]fakeFund{})
	obj := newFundObject(t, "", map[string]string{"name": "General"})

	sut := fundStore{rep}
	_, jsherr := sut.Save(context.Background(), obj)

	assert.Equal(t, StatusUnprocessableEntity, jsherr.StatusCode(),
		"fundStore gave unexpect status on Save()")
}

func TestFundStoreSaveWithInvalidCurrencyIsError(t *testing.T) {
	// if this test starts failing on the "UUU" check if a new
	// currency with code UUU has been added.
	badcurrencies := []string{"#$%", "BADDD", "UUU"}

	for _, badcurrency := range badcurrencies {
		rep := newFakeFundRepository([]fakeFund{})
		obj := newFundObject(t, "",
			map[string]string{"currency": badcurrency, "name": "General"})

		sut := fundStore{rep}
		_, jsherr := sut.Save(context.Background(), obj)

		assert.Equal(t, StatusUnprocessableEntity, jsherr.StatusCode(),
			"fundStore gave unexpect status on Save()")
	}
}

func TestFundStoreSaveCreatesFundInDomain(t *testing.T) {
	rep := newFakeFundRepository([]fakeFund{})
	obj := newFundObject(t, "", map[string]string{"currency": "CAD", "name": "General"})

	sut := fundStore{rep}
	_, jsherr := sut.Save(context.Background(), obj)

	assert.Nil(t, jsherr, "fundStore gave unexpected error on Save()")
	assert.True(t, rep.createFundCalled, "fundStore did not call CreateFund()")
	require.Len(t, rep.funds, 1, "unexpected number of funds in the fake repository")
	assert.Equal(t, "General", rep.funds[0].Name(), "unexpected fund name in the repository")
	assert.Equal(t, domain.CAD, rep.funds[0].Currency(), "unexpected currency in the repository")
}

func TestFundStoreSaveReturnsCreatedFund(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	rep := newFakeFundRepository([]fakeFund{})
	obj := newFundObject(t, "", map[string]string{"currency": "CAD", "name": "General"})

	sut := fundStore{rep}
	actual, jsherr := sut.Save(context.Background(), obj)

	require.Nil(jsherr, "fundStore gave unexpted error on Save()")
	require.NotNil(actual, "fundStore returned nil fund on Save()")
	assert.Equal("1", actual.ID, "returned fund had unexpected ID")
	assert.JSONEq(`{"currency" : "CAD", "name" : "General"}`, string(actual.Attributes),
		"Unexpected attributes on the returned fund")
}

func TestNewFundResource(t *testing.T) {
	assert := assert.New(t)
	fakerepository := newFakeFundRepository([]fakeFund{{1, domain.CAD, "General"},
		{2, domain.USD, "Special"}})
	request, respsonsewriter := getRequestResponse(t, "/fund")

	sut := jshapi.New("/")
	sut.Add(newFundResource(fakerepository))
	sut.ServeHTTPC(context.Background(), respsonsewriter, request)

	assert.Equal(http.StatusOK, respsonsewriter.Code, "Unexpected status code.")
	assert.True(fakerepository.getAllCalled, "Listing the funds failed to call GetAll()")
}

func TestNewFundResourceIncludesAttributes(t *testing.T) {
	assert := assert.New(t)
	fakerepository := newFakeFundRepository([]fakeFund{{1, domain.CAD, "General"},
		{2, domain.USD, "Special"}})
	request, respsonsewriter := getRequestResponse(t, "/fund")

	sut := jshapi.New("/")
	sut.Add(newFundResource(fakerepository))
	sut.ServeHTTPC(context.Background(), respsonsewriter, request)

	doc := parseResponseBody(t, respsonsewriter, jsh.ListMode)
	assert.True(doc.HasData(), "Returned document unexpected has no data.")
	for _, obj := range doc.Data {
		switch obj.ID {
		case "1":
			assert.JSONEq(`{"currency" : "CAD", "name" : "General"}`,
				string(obj.Attributes),
				"Unexpected attributes on an object in the returned document.")
		case "2":
			assert.JSONEq(`{"currency" : "USD", "name" : "Special"}`,
				string(obj.Attributes),
				"Unexpected attributes on an object in the returned document.")
		}

	}
}

func TestNewFundResorceWithEmptyReporsitory(t *testing.T) {
	assert := assert.New(t)
	var fakerepository fakeFundRepository
	request, respsonsewriter := getRequestResponse(t, "/fund")

	sut := jshapi.New("/")
	sut.Add(newFundResource(&fakerepository))
	sut.ServeHTTPC(context.Background(), respsonsewriter, request)

	assert.Equal(http.StatusOK, respsonsewriter.Code, "Unexpected status code.")
	assert.True(fakerepository.getAllCalled, "Listing the funds failed to call GetAll()")
}
