// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package domain

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	mysqlDialect string = "mysql"
)

func New(dsn string) (Store, error) {
	db, err := gorm.Open(mysqlDialect, dsn)
	if err != nil {
		return nil, err
	}

	return &store{db}, nil
}

func CreateOrMigrate(dsn string) error {
	db, err := gorm.Open(mysqlDialect, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.AutoMigrate(&fund{}).Error
	if err != nil {
		return err
	}

	return nil
}
