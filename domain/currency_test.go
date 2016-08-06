// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type currencyCodePair struct {
	currency Currency
	code     string
}

func TestCurrencyStringGivesExpectedCurrencyCodes(t *testing.T) {
	assert := assert.New(t)
	expected := []currencyCodePair{
		{CAD, "CAD"}, {USD, "USD"},
		{EUR, "EUR"}, {GBP, "GBP"},
		{XXX, "XXX"}, {XTS, "XTS"},
		{ZWL, "ZWL"}, {ZWL + 1, "XXX"},
		{ZWL + 2, "XXX"},
	}

	for _, pair := range expected {
		assert.Equal(pair.code, pair.currency.String())
	}
}
