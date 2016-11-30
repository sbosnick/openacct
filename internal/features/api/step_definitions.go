// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"

	jsh "github.com/derekdowling/go-json-spec-handler"
	jsc "github.com/derekdowling/go-json-spec-handler/client"
	"github.com/go-sql-driver/mysql"
	. "github.com/gucumber/gucumber"
	"github.com/sbosnick1/openacct/cmd/openacctapi"
	"github.com/sbosnick1/openacct/domain"
)

var worldServerKey = "server"

func setServer(srv *httptest.Server) {
	World[worldServerKey] = srv
}

func getServer() *httptest.Server {
	return World[worldServerKey].(*httptest.Server)
}

func getRootURL() string {
	srv := getServer()
	if srv == nil {
		log.Panic("httptest Server is nil in World map.")
	}

	return srv.URL
}

func getBaseURL() string {
	root, err := url.Parse(getRootURL())
	if err != nil {
		log.Panic(err)
	}

	apiurl, err := url.Parse("v1/")
	if err != nil {
		log.Panic(err)
	}

	return root.ResolveReference(apiurl).String()
}

func getDsn() string {
	return "/openacct"
}

func openServer() {
	handler, err := openacctapi.BuildApiHandler(getDsn())
	if err != nil {
		log.Fatal(err)
	}

	setServer(httptest.NewServer(handler))
}

func closeServer() {
	srv := getServer()
	if srv == nil {
		log.Panic("httptest Server is nil in World map.")
	}
	srv.Close()
	setServer(nil)
}

func cleanDb() {
	dsn := getDsn()

	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	_, err = db.Exec(create)
	if err != nil {
		log.Fatal(err)
	}

	err = domain.CreateOrMigrate(dsn)
	if err != nil {
		log.Fatal(err)
	}
}

// This function is used to represent the default state when starting from
// an empty database.
func doNothing() {
}

func checkFundCount(count int) {
	doc, resp, err := jsc.List(getBaseURL(), "fund")
	if err != nil {
		T.Errorf(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		T.Errorf("Unexpected status code in the response: %s", resp.Status)
		return
	}

	if doc == nil {
		T.Errorf("No jsh document returned.")
		return
	}

	if doc.HasErrors() {
		T.Errorf("The returned jsh document had the following errors %s", doc.Error())
		return
	}

	if doc.Data == nil {
		T.Errorf("The returned jsh document did not have any data.")
		return
	}

	if doc.Mode != jsh.ListMode {
		T.Errorf("The returned jsh document was not in list mode.")
		return
	}

	if len(doc.Data) != count {
		T.Errorf("The returned list had %d entries but %d were expected", len(doc.Data), count)
		return
	}
}

func addFund(fundName string, currency string) {
	object, jsherr := jsh.NewObject("", "fund",
		map[string]string{"name": fundName, "currency": currency})
	if jsherr != nil {
		T.Errorf(jsherr.Error())
		return
	}

	_, resp, err := jsc.Post(getBaseURL(), object)
	if err != nil {
		T.Errorf(err.Error())
		return
	}

	if resp.StatusCode != http.StatusCreated {
		T.Errorf("Unexpected status code in the response: %s", resp.Status)
		return
	}
}

func checkForFund(fundName string, currency string) {
	doc, resp, err := jsc.List(getBaseURL(), "fund")
	if err != nil {
		T.Errorf(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		T.Errorf("Unexpected status code in the response: %s", resp.Status)
		return
	}

	if doc.HasErrors() {
		T.Errorf("Returned document has errors: %s", doc.Error())
		return
	}

	if !doc.HasData() {
		T.Errorf("Returned document has no data.")
		return
	}

	for _, object := range doc.Data {
		var attributes = make(map[string]string)
		jsherr := object.Unmarshal("fund", &attributes)
		if jsherr != nil {
			T.Errorf(jsherr.Error())
		}

		if attributes["name"] == fundName && attributes["currency"] == currency {
			return
		}
	}

	T.Errorf("The returned list of funds did not include one with a name %s and currency %s.",
		fundName, currency)
}

func init() {
	BeforeAll(openServer)

	AfterAll(closeServer)

	Before("@cleandb", cleanDb)

	When(`^the bookkeeper has not added any funds$`, doNothing)

	Then(`^the list of funds has (\d+) entr(y|ies).?$`, func(count int, _ string) {
		checkFundCount(count)
	})

	When(`^the bookkeeper adds the "(.+?)" fund in "(.+?)" currency$`, addFund)

	And(`^there is a "(.+?)" fund demonicated in "(.+?)" currency.$`, checkForFund)

	Given(`^that the bookkeeper has added the "(.+?)" fund in "(.+?)" currency$`, addFund)

	Given(`^that the bookkeeper has added the following funds$`, func(table [][]string) {
		T.Skip() // pending
	})

	When(`^the bookkeeper deletes the "(.+?)" fund$`, func(s1 string) {
		T.Skip() // pending
	})

	And(`^there is not a "(.+?)" fund.$`, func(s1 string) {
		T.Skip() // pending
	})

}
