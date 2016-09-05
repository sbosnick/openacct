# Copyright Steven Bosnick 2016. All rights reserved.
# Use of this source code is governed by the GNU General Public License version 3.
# See the file COPYING for your rights under that license.

@cleandb
Feature: Create, List and Delete Funds
    Bookkeepers should be able to create new funds, and list and delete existing funds.

    A fund is a set of accounts, all of which are denomined in the same currency.
    Creating a fund requires specifying the currency for that fund. All transactions
    relating to accounts in the the same fund are demonicated using the currency
    of that fund.

    @wip
    Scenario: Bookkeeper lists empty funds
        When the bookkeeper has not added any funds
        Then the list of funds has 0 entries.

    @ignore
    Scenario: Bookkeeper adds a fund
        When the bookkeeper adds the "General" fund in "CDN" currency
        Then the list of funds has 1 entry
        And there is a "General" fund demonicated in "CDN" currency.

    @ignore
    Scenario: Bookkeeper adds a second fund
        Given that the bookkeeper has added the "General" fund in "CDN" currency
        When the bookkeeper adds the "USGeneral" fund in "USD" currency
        Then the list of funds has 2 entries
        And there is a "USGeneral" fund demonicated in "USD" currency.

    @ignore
    Scenario: Bookkeeper deletes one of two funds
        Given that the bookkeeper has added the following funds
            | fundname  | currency  |
            | General   | CDN       |
            | USGeneral | USD       |
        When the bookkeeper deletes the "USGeneral" fund
        Then the list of funds has 1 entry
        And there is not a "USGeneral" fund.
