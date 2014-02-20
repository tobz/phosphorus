package models

import "time"
import "database/sql"
import "github.com/tobz/phosphorus/utils"

type Account struct {
	AccountID uint64         `db:"account_id"`
	Username  string         `db:"username"`
	Password  string         `db:"password"`
	Realm     uint8          `db:"realm"`
	Email     sql.NullString `db:"email"`
	Created   time.Time      `db:"created_dt"`
	LastLogin time.Time      `db:"last_login_dt"`
	Status    uint64         `db:"status"`
}

func NewAccount(username, password string) (*Account, error) {
	a := &Account{Username: username, Created: time.Now(), Status: 0}

	// Set our password, which will encrypt it.
	err := a.setPassword(password)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Account) PasswordMatches(password string) (bool, error) {
	s, err := utils.LoadScryptFromHash(a.Password)
	if err != nil {
		return false, err
	}

	matches, err := s.MatchesPlaintext(password)
	if err != nil {
		return false, err
	}

	return matches, nil
}

func (a *Account) setPassword(password string) error {
	// Encrypt our password.
	s, err := utils.NewScrypt()
	if err != nil {
		return err
	}

	passwordHash, err := s.HashPlaintext(password)
	if err != nil {
		return err
	}

	a.Password = passwordHash

	return nil
}
