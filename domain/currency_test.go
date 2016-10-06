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

func TestParseCurrencyGiveExpectedCurrencies(t *testing.T) {
	expected := []currencyCodePair{
		{AED, "AED"}, {ZWL, "ZWL"},
		{CAD, "CAD"}, {USD, "USD"},
		{EUR, "EUR"}, {GBP, "GBP"},
		{XXX, "XXX"}, {XTS, "XTS"},
		{CAD, "cad"}, {USD, "usd"},
		{EUR, "eur"}, {GBP, "gbp"},
		{XXX, "xxx"}, {XTS, "xts"},
	}

	for _, pair := range expected {
		actual, err := ParseCurrency(pair.code)
		assert.NoError(t, err, "ParseCurrency() returned an unexpected error.")
		assert.Equal(t, pair.currency, actual)
	}
}

func TestParseCurrencyWithShortOrLongInputIsError(t *testing.T) {
	badinput := []string{"", "a", "A", "AA", "aa", "abcd", "ABCD"}

	for _, input := range badinput {
		_, err := ParseCurrency(input)
		assert.Error(t, err, "ParseCurrency() failed to return an expected error.")
	}
}

func TestParseCurrencyWithBadInputIsError(t *testing.T) {
	// if this test fails on the "uuu" or "UUU" input check if a new currency
	// code "UUU" has been defined. This tests assumes that "UUU" is not a currency
	// code.
	badinput := []string{"123", "ae3", "*j*", "uuu", "UUU"}

	for _, input := range badinput {
		_, err := ParseCurrency(input)
		assert.Error(t, err, "ParseCurrency() failed to return an expected error.")
	}
}

func TestIsCurrencyIsTrueForValidCurrencies(t *testing.T) {
	valid := []string{"AED", "ZWL", "CAD", "USD", "EUR", "GBP", "XXX", "XTS",
		"cad", "usd", "eur", "gbp", "xxx", "xts"}

	for _, value := range valid {
		assert.True(t, IsCurrency(value))
	}
}

func TestIsCurrencyIsFalseWithShortOrLongInput(t *testing.T) {
	badinput := []string{"", "a", "A", "AA", "aa", "abcd", "ABCD"}

	for _, value := range badinput {
		assert.False(t, IsCurrency(value))
	}
}

func TestIsCurrencyIsFalseWithBadInput(t *testing.T) {
	// if this test fails on the "uuu" or "UUU" input check if a new currency
	// code "UUU" has been defined. This tests assumes that "UUU" is not a currency
	// code.
	badinput := []string{"123", "ae3", "*j*", "uuu", "UUU"}

	for _, value := range badinput {
		assert.False(t, IsCurrency(value))
	}
}
