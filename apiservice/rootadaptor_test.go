// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package apiservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"goji.io"
	"golang.org/x/net/context"
)

func getRequestResponse(t *testing.T, urlStr string) (*http.Request, *httptest.ResponseRecorder) {
	request, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		t.Fatal("unable to create request.")
	}

	response := httptest.NewRecorder()

	return request, response
}

func TestZeroRootAdaptorRespondsNotFound(t *testing.T) {
	request, response := getRequestResponse(t, "/anything")

	var sut rootAdaptor
	sut.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Zero rootAdaptor responded with code '%s' but code '%s' was expected",
			http.StatusText(response.Code), http.StatusText(http.StatusNotFound))
	}
}

func TestRootAdaptorForwardsRequestAndResponse(t *testing.T) {
	request, response := getRequestResponse(t, "/anything")
	var actualRequest *http.Request
	var actualResponse http.ResponseWriter
	rootHandler := func(_ context.Context, resp http.ResponseWriter, req *http.Request) {
		actualRequest = req
		actualResponse = resp
	}

	sut := rootAdaptor{goji.HandlerFunc(rootHandler)}
	sut.ServeHTTP(response, request)

	if actualRequest != request {
		t.Errorf("rootAdaptor failed to forward the request")
	}
	if actualResponse != response {
		t.Errorf("rootAdaptor failed to forward the response writer")
	}
}

func TestRootAdaptorUsesBackgroundContest(t *testing.T) {
	request, response := getRequestResponse(t, "/anything")
	var actualContext context.Context
	rootHandler := func(ctx context.Context, _ http.ResponseWriter, _ *http.Request) {
		actualContext = ctx
	}

	sut := rootAdaptor{goji.HandlerFunc(rootHandler)}
	sut.ServeHTTP(response, request)

	if actualContext != context.Background() {
		t.Errorf("rootAdaptor failed to use context.Background() as the context")
	}
}
