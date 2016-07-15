package api

import (
	. "github.com/gucumber/gucumber"
	//"net/http/httptest"
)

func init() {
	BeforeAll(func() {
		//World["server"] = httptest.NewServer(handler)
	})

	AfterAll(func() {
		//server := World["server"].(httptest.Server)
		//server.Close()
	})

	When(`^the bookkeeper lists the funds$`, func() {
		T.Skip() // pending
	})

	Then(`^the list of funds has (\d+) entries.$`, func(i1 int) {
		T.Skip() // pending
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
