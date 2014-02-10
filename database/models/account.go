package models

import "time"
import "database/sql"
import "github.com/tobz/phosphorus/database"
import "github.com/tobz/phosphorus/utils"
import "github.com/tobz/phosphorus/log"

func init() {
    database.RegisterTableSchema(Account{}, "account")
}

type Account struct {
    id uint64 `db:"account_id"`
    username string
    password string
    email sql.NullString
    created time.Time `db:"created_dt"`
    lastLogin time.Time `db:"last_login_dt"`
    status uint64
}

func NewAccount(username, password string) (*Account, error) {
    a := &Account{ username: username }

    // Set our password, which will encrypt it.
    err := a.SetPassword(password)
    if err != nil {
        return nil, err
    }

    return a, nil
}

func (a *Account) PasswordMatches(password string) bool {
    s, err := utils.LoadScryptFromHash(a.password)
    if err != nil {
        log.Server.Error("account", "caught error while trying to validate password for '%s': %s", a.Username, err)
        return false
    }

    matches, err := s.MatchesPlaintext(password)
    if err != nil {
        log.Server.Error("account", "caught error while trying to validate password for '%s': %s", a.Username, err)
        return false
    }

    return matches
}

func (a *Account) SetUsername(username string) {
    a.username = username
}

func (a *Account) Username() string {
    return a.username
}

func (a *Account) SetPassword(password string) error {
    // Encrypt our password.
    s, err := utils.NewScrypt()
    if err != nil {
        return err
    }

    passwordHash, err := s.HashPlaintext(password)
    if err != nil {
        return err
    }

    a.password = passwordHash

    return nil
}

func (a *Account) SetEmail(email string) {
    if email != "" {
        a.email.String = email
        a.email.Valid = true
    } else {
        a.email.Valid = false
    }
}

func (a *Account) Email() string {
    if a.email.Valid {
        return a.email.String
    }

    return ""
}

func (a *Account) SetLastLogin(lastLogin time.Time) {
    a.lastLogin = lastLogin
}

func (a *Account) LastLogin() time.Time {
    return a.lastLogin
}

func (a *Account) SetCreated(created time.Time) {
    a.created = created
}

func (a *Account) Created() time.Time {
    return a.created
}

