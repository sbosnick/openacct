package domain

// A Fund is a named collection of accounts, all of which are
// demoniated in the same currency.
type Fund interface {
	Id() uint
	Currency() Currency
	Name() string
}

// The FundRepository is the means of accessing the Fund's in the store.
type FundRepository interface {
	GetAll() ([]Fund, error)
}
