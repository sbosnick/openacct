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

func insertFunds(t *testing.T, db *gorm.DB, funds []fund) {
	for _, fund := range funds {
		db.Create(&fund)
		require.NoError(t, db.Error, "Unable to create a fund.")
	}
}

func TestFundRepositoryGetAllRetreivesAllFunds(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db := getEmptyDb(t)
	expected := []fund{{1, CAD, "General"}, {2, USD, "Special"}}
	insertFunds(t, db, expected)

	sut := fundRepository{db}
	actual, err := sut.GetAll()

	require.NoError(err, "Unable to get all funds.")
	require.NotNil(actual, "GetAll() returned nil funds list.")
	assert.Equal(len(expected), len(actual), "Unexpected number of funds returned from GetAll().")
	for _, f := range actual {
		assert.Contains(expected, *f.(*fund), "Unexpected fund returned from GetAll().")
	}
}
