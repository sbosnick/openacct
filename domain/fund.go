// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package domain

import "github.com/jinzhu/gorm"

// A Fund is a named collection of accounts, all of which are
// demoniated in the same currency.
type Fund interface {
	Id() uint
	Currency() Currency
	Name() string
}

type fundImpl struct {
	ID           uint
	FundCurrency Currency
	FundName     string `sql:"size:255;unique;index`
}

func (f *fundImpl) Id() uint {
	return f.ID
}

func (f *fundImpl) Currency() Currency {
	return f.FundCurrency
}

func (f *fundImpl) Name() string {
	return f.FundName
}

// The FundRepository is the means of accessing the Fund's in the store.
type FundRepository interface {
	GetAll() ([]Fund, error)
	Create(name string, currency Currency) (Fund, error)
}

type fundRepository struct {
	db *gorm.DB
}

func (f *fundRepository) GetAll() ([]Fund, error) {
	var funds []fundImpl

	err := f.db.Find(&funds).Error
	if err != nil {
		return nil, err
	}

	var ret []Fund
	for _, f := range funds {
		ret = append(ret, &f)
	}

	return ret, nil
}

func (f *fundRepository) Create(name string, currency Currency) (Fund, error) {
	fund := fundImpl{FundName: name, FundCurrency: currency}

	err := f.db.Create(&fund).Error
	if err != nil {
		return nil, err
	}

	return &fund, nil
}
