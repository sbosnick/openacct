// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

//go:generate go run mkcurrency.go

package domain

type Currency uint

type currencyInfo struct {
	symbol     string
	name       string
	number     int
	minorUnits int
}
