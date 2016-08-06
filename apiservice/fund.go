// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package apiservice

import (
	"strconv"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/sbosnick1/openacct/domain"
	"golang.org/x/net/context"
)

const (
	fundResourceType = "fund"
)

type fundAttributes struct {
	Name     string `json:"name,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// A fundStore is a store for the fund resorce type. It adapts a
// domain.FundRepository to a json api spec. resource.
type fundStore struct {
	repository domain.FundRepository
}

func (f *fundStore) Save(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	panic("not implemented")
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

	var list jsh.List
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
