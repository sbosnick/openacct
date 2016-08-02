package apiservice

import (
	"net/http"
	"testing"

	"github.com/sbosnick1/openacct/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
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
	getAllCalled bool
	funds        []domain.Fund
}

func (f *fakeFundRepository) GetAll() ([]domain.Fund, error) {
	f.getAllCalled = true
	return f.funds, nil
}

func newFakeFundRepository(funds []fakeFund) *fakeFundRepository {
	var realfunds []domain.Fund
	for _, f := range funds {
		realfunds = append(realfunds, &f)
	}
	return &fakeFundRepository{false, realfunds}
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
