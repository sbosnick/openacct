// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package domain

import "github.com/jinzhu/gorm"

// A Store is an abstract factory for the repositories that provide
// persistance for the entities in the domain. It represents the
// collection of repositories that make up the domain.
type Store interface {
	FundRepository() FundRepository
}

type store struct {
	db *gorm.DB
}

func (s *store) FundRepository() FundRepository {
	return &fundRepository{s.db}
}
