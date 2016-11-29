// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package openacctapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	jsh "github.com/derekdowling/go-json-spec-handler"
	jsc "github.com/derekdowling/go-json-spec-handler/client"
	"github.com/go-sql-driver/mysql"
	"github.com/sbosnick1/openacct/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildDsn(*testing.T) string {
	return "/openacct"
}

func resetDb(t *testing.T) {
	dsn := buildDsn(t)

	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// This creates a risk of a SQL injection attack, but we
	// are about to drop the database so there isn't much worse
	// things that could be done. Plus this is testing code, not
	// production.
	drop := fmt.Sprintf("drop database if exists %s", config.DBName)
	create := fmt.Sprintf("create database %s", config.DBName)

	_, err = db.Exec(drop)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(create)
	if err != nil {
		t.Fatal(err)
	}

	err = domain.CreateOrMigrate(dsn)
	if err != nil {
		t.Fatal(err)
	}
}

func addFundToDb(t *testing.T) {
	db, err := sql.Open("mysql", buildDsn(t))
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	_, err = db.Exec("insert into fund_impls(id, fund_currency, fund_name) value(1, 26, 'General')")
	if err != nil {
		t.Error(err)
	}
}

func buildHandler(t *testing.T) http.Handler {
	handler, err := BuildApiHandler(buildDsn(t))
	if err != nil {
		t.Error(err)
	}

	return handler
}

func getBaseURL(t *testing.T, base string) *url.URL {
	baseurl, err := url.Parse(base)
	if err != nil {
		t.Error(err)
	}

	baseurl, err = baseurl.Parse("/v1/")
	if err != nil {
		t.Error(err)
	}

	return baseurl
}

func jscListFunds(t *testing.T, baseurl string) *jsh.Document {
	assert := assert.New(t)
	require := require.New(t)

	doc, resp, err := jsc.List(getBaseURL(t, baseurl).String(), "fund")

	require.NoError(err, "Error listing funds.")
	require.NotNil(resp, "nil response.")
	require.NotNil(doc, "nil document.")
	assert.Equal(http.StatusOK, resp.StatusCode, "Unexpected status code.")
	assert.False(doc.HasErrors(), "Returned document has errors.")
	assert.True(doc.HasData(), "Returned document has no data.")

	return doc
}

func TestApiCanRetreiveFundRowsUsingJsonSpecClient(t *testing.T) {
	assert := assert.New(t)
	resetDb(t)
	addFundToDb(t)

	sut := httptest.NewServer(buildHandler(t))
	defer sut.Close()
	doc := jscListFunds(t, sut.URL)

	assert.JSONEq(`{"name":"General", "currency":"CAD"}`,
		string(doc.First().Attributes), "Unexpected attributes.")
}
