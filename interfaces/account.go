package interfaces

import "time"

type Account interface {
    SetUsername(string)
    Username() string

    SetPassword(string) error

    SetEmail(string)
    Email() string

    SetLastLogin(time.Time)
    LastLogin() time.Time

    SetCreated(time.Time)
    Created() time.Time
}
