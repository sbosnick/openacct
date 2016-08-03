// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package api

import (
	jsh "github.com/derekdowling/go-json-spec-handler"
	jsc "github.com/derekdowling/go-json-spec-handler/client"
	. "github.com/gucumber/gucumber"
	"github.com/sbosnick1/openacct/cmd/openacctapi"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var worldServerKey = "server"

func setServer(srv *httptest.Server) {
	World[worldServerKey] = srv
}

func getServer() *httptest.Server {
	return World[worldServerKey].(*httptest.Server)
}

func closeServer() {
	srv := getServer()
	if srv == nil {
		log.Panic("httptest Server is nil in World map.")
	}
	srv.Close()
	setServer(nil)
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

	apiurl, err := url.Parse("api/v1/")
	if err != nil {
		log.Panic(err)
	}

	return root.ResolveReference(apiurl).String()
}

func init() {
	BeforeAll(func() {
		handler, err := openacctapi.BuildApiHandler()
		if err != nil {
			log.Fatal(err)
		}

		setServer(httptest.NewServer(handler))
	})

	AfterAll(func() {
		closeServer()
	})

	When(`^the bookkeeper has not added any funds$`, func() {
		// Do nothing. This should be the default state when starting from
		// an empty database.
	})

	Then(`^the list of funds has (\d+) entries.$`, func(count int) {
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

		if !doc.HasData() {
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
	})

	When(`^the bookkeeper adds the "(.+?)" fund in "(.+?)" currency$`, func(s1 string, s2 string) {
		T.Skip() // pending
	})

	Then(`^the list of funds has (\d+) entry$`, func(i1 int) {
		T.Skip() // pending
	})

	And(`^there is a "(.+?)" fund demonicated in "(.+?)" currency.$`, func(s1 string, s2 string) {
		T.Skip() // pending
	})

	Given(`^that the bookkeeper has added the "(.+?)" fund in "(.+?)" currency$`, func(s1 string, s2 string) {
		T.Skip() // pending
	})

	Then(`^the list of funds has (\d+) entries$`, func(i1 int) {
		T.Skip() // pending
	})

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
