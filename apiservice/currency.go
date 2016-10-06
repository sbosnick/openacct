// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package apiservice

import (
	"encoding/json"

	"github.com/sbosnick1/openacct/domain"
)

type currency domain.Currency

func (c *currency) MarshalJSON() ([]byte, error) {
	value := domain.Currency(*c)
	return json.Marshal(value.String())
}

func (c *currency) UnmarshalJSON(data []byte) error {
	var value string
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	curr, err := domain.ParseCurrency(value)
	if err != nil {
		return err
	}

	*c = currency(curr)
	return nil
}
