package models

type Account struct {
    name string
}

func (a *Account) Name() {
    return a.name
}
