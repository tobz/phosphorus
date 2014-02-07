package models

type Account struct {
	name string
}

func (a *Account) Name() string {
	return a.name
}
