// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package domain

import (
	"database/sql"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeDsn() string {
	database := os.Getenv("OPENACCT_DB_DATABASE")
	if len(database) == 0 {
		database = "openacct"
	}

	host := os.Getenv("OPENACCT_DB_HOST")
	var net string
	if len(host) > 0 {
		net = "tcp"
	}

	config := &mysql.Config{
		User:   os.Getenv("OPENACCT_DB_USER"),
		Passwd: os.Getenv("OPENACCT_DB_PASSWD"),
		Net:    net,
		Addr:   host,
		DBName: database,
	}

	return config.FormatDSN()
}

func deleteAllTables(t *testing.T, dsn string) {
	require := require.New(t)

	config, err := mysql.ParseDSN(dsn)
	require.NoError(err, "Unable to parse dsn")

	db, err := sql.Open("mysql", dsn)
	require.NoError(err, "Unable to open database to delete tables.")
	defer db.Close()

	rows, err := db.Query(
		"select concat('drop table if exists ', table_name, ';') "+
			"from information_schema.tables "+
			"where table_schema=?", config.DBName)
	require.NoError(err, "Unable to query the list of tables.")
	defer rows.Close()

	var droppers []string
	for rows.Next() {
		var dropcmd string
		err = rows.Scan(&dropcmd)
		require.NoError(err, "Unable to scan the tablename from the rows")

		droppers = append(droppers, dropcmd)
	}

	for _, dropper := range droppers {
		_, err := db.Exec(dropper)
		require.NoError(err, "Unable to drop tables.")
	}
}

func getEmptyDb(t *testing.T) *gorm.DB {
	dsn := makeDsn()

	deleteAllTables(t, dsn)
	err := CreateOrMigrate(dsn)
	require.NoError(t, err, "Unable to create the database schema.")

	db, err := gorm.Open("mysql", dsn)
	require.NoError(t, err, "Unable to open the database.")

	return db
}

func TestCreateOrMigrate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	dsn := makeDsn()
	deleteAllTables(t, dsn)

	err := CreateOrMigrate(dsn)

	require.NoError(err, "CreateOrMigrate() failed.")
	db, err := gorm.Open("mysql", dsn)
	require.NoError(err, "gorm.Open() failed.")
	defer db.Close()
	assert.True(db.HasTable(&fund{}))
}
