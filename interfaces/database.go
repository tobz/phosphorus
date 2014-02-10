package interfaces

import "database/sql"

type Database interface {
    Begin() (Transaction, error)
}

type Transaction interface {
    Insert(...interface{}) error
    Update(...interface{}) (int64, error)
    Delete(...interface{}) (int64, error)
    Get(interface{}, ...interface{}) (interface{}, error)
    Select(interface{}, string, ...interface{}) ([]interface{}, error)
    SelectOne(interface{}, string, ...interface{}) error
    Exec(string, ...interface{}) (sql.Result, error)
    Commit() error
    Rollback() error
}
