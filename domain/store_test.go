// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package domain

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestStoreFundRepositoryFowardsDb(t *testing.T) {
	assert := assert.New(t)
	fakedb := &gorm.DB{}

	sut := store{fakedb}
	repo := sut.FundRepository()

	realrepo := repo.(*fundRepository)
	assert.Equal(fakedb, realrepo.db)
}
