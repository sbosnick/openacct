// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package domain

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func insertFunds(t *testing.T, db *gorm.DB, funds []fundImpl) {
	for _, fund := range funds {
		db.Create(&fund)
		require.NoError(t, db.Error, "Unable to create a fund.")
	}
}

func TestFundRepositoryGetAllRetreivesAllFunds(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db := getEmptyDb(t)
	expected := []fundImpl{{1, CAD, "General"}, {2, USD, "Special"}}
	insertFunds(t, db, expected)

	sut := fundRepository{db}
	actual, err := sut.GetAll()

	require.NoError(err, "Unable to get all funds.")
	require.NotNil(actual, "GetAll() returned nil funds list.")
	assert.Equal(len(expected), len(actual), "Unexpected number of funds returned from GetAll().")
	for _, f := range actual {
		assert.Contains(expected, *f.(*fundImpl), "Unexpected fund returned from GetAll().")
	}
}

func TestFundRepositoryCreateAddsFund(t *testing.T) {
	db := getEmptyDb(t)

	sut := fundRepository{db}
	_, err := sut.Create("General", CAD)

	require.NoError(t, err, "Unable to create new fund")
	var f fundImpl
	err = db.Where("fund_name = ?", "General").First(&f).Error
	require.NoError(t, err, "Unable to query expected fund")
	assert.Equal(t, "General", f.FundName, "Unexpected fund name")
	assert.Equal(t, CAD, f.FundCurrency, "Unexpected fund currency")
}

func TestFundRepositoryCreateReturnsNewFund(t *testing.T) {
	db := getEmptyDb(t)

	sut := fundRepository{db}
	actual, err := sut.Create("General", CAD)

	require.NoError(t, err, "Unable to create new fund")
	require.NotNil(t, actual, "Returned fund was nil")
	assert.NotZero(t, actual.Id(), "Id of the returned fund was zero")
	assert.Equal(t, "General", actual.Name())
	assert.Equal(t, CAD, actual.Currency())
}
