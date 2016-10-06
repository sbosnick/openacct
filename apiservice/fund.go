// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package apiservice

import (
	"strconv"

	jsh "github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/jsh-api"
	"github.com/sbosnick1/openacct/domain"
	"golang.org/x/net/context"
)

const (
	fundResourceType = "fund"
)

func newFundResource(repository domain.FundRepository) *jshapi.Resource {
	return jshapi.NewCRUDResource(fundResourceType, &fundStore{repository})
}

type fundAttributes struct {
	Name     string `json:"name,omitempty" valid:"required,utfletternum"`
	Currency string `json:"currency,omitempty" valid:"required,currency"`
}

// A fundStore is a store for the fund resorce type. It adapts a
// domain.FundRepository to a json api spec. resource.
type fundStore struct {
	repository domain.FundRepository
}

func (f *fundStore) Save(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	if f.repository == nil {
		return nil, jsh.ISE("fundStore requires a FundRepository")
	}

	var attributes fundAttributes
	jsherrs := object.Unmarshal(fundResourceType, &attributes)
	if jsherrs != nil {
		return nil, jsherrs
	}

	currency, err := domain.ParseCurrency(attributes.Currency)
	if err != nil {
		// the validation on fundAttributes should have ensured this
		// does not happen
		return nil, jsh.ISE(err.Error())
	}

	fund, err := f.repository.Create(attributes.Name, currency)
	if err != nil {
		return nil, jsh.ISE(err.Error())
	}

	obj, jsherr := createFundObject(fund)
	if jsherr != nil {
		return nil, jsherr
	}

	return obj, nil
}

func (f *fundStore) Get(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType) {
	panic("not implemented")
}

func (f *fundStore) List(ctx context.Context) (jsh.List, jsh.ErrorType) {
	if f.repository == nil {
		return nil, jsh.ISE("fundStore requires a FundRepository")
	}

	fund, err := f.repository.GetAll()
	if err != nil {
		return nil, jsh.ISE(err.Error())
	}

	list := make(jsh.List, 0)
	for _, f := range fund {
		obj, err := createFundObject(f)
		if err != nil {
			return nil, err
		}

		list = append(list, obj)
	}

	return list, nil
}

func (f *fundStore) Update(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	panic("not implemented")
}

func (f *fundStore) Delete(ctx context.Context, id string) jsh.ErrorType {
	panic("not implemented")
}

func createFundObject(fund domain.Fund) (*jsh.Object, *jsh.Error) {
	id := strconv.FormatUint(uint64(fund.Id()), 10)

	obj, err := jsh.NewObject(id, fundResourceType,
		fundAttributes{Name: fund.Name(), Currency: fund.Currency().String()})
	if err != nil {
		return nil, err
	}

	return obj, nil
}
