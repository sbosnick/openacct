// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package apiservice

import (
	"encoding/json"
	"testing"

	"github.com/sbosnick1/openacct/domain"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	C currency
}

func TestCurrencyMarshalsToExpectedString(t *testing.T) {
	testcases := map[domain.Currency]string{
		domain.CAD: `"CAD"`,
		domain.USD: `"USD"`,
		domain.EUR: `"EUR"`,
	}

	for value, expected := range testcases {
		var sut currency = currency(value)
		actual, err := json.Marshal(&sut)

		assert.NoError(t, err, "Marshalling currency returned unexpected error.")
		assert.Equal(t, []byte(expected), actual,
			"Mashalling currency return unexpected bytes.")
	}
}

func TestCurrencyUnmarshalsFromValidStrings(t *testing.T) {
	testcases := map[string]domain.Currency{
		`"CAD"`: domain.CAD,
		`"cad"`: domain.CAD,
		`"CaD"`: domain.CAD,
		`"USD"`: domain.USD,
		`"usd"`: domain.USD,
		`"USd"`: domain.USD,
		`"EUR"`: domain.EUR,
	}

	for value, expected := range testcases {
		var sut currency
		err := json.Unmarshal([]byte(value), &sut)

		assert.NoError(t, err, "Unmarshalling currency returned unexpected error.")
		assert.Equal(t, sut, currency(expected),
			"Unmarshaling currency returned unexpected currency value.")
	}
}

func TestCurrencyInStructMarshalsToExpectedString(t *testing.T) {
	expected := `{"C":"CAD"}`

	sut := testStruct{currency(domain.CAD)}
	actual, err := json.Marshal(&sut)

	assert.NoError(t, err, "Marshalling struct with currency returned unexpected error.")
	assert.Equal(t, expected, string(actual),
		"Marshalling struct with currency returned unexpected bytes.")
}
