package interfaces

import "database/sql"

type Database interface {
	Begin() (Transaction, error)

	Insert(...interface{}) error
	Update(...interface{}) (int64, error)
	Delete(...interface{}) (int64, error)
	Get(interface{}, ...interface{}) (interface{}, error)
	Select(interface{}, string, ...interface{}) ([]interface{}, error)
	SelectOne(interface{}, string, ...interface{}) error
	Exec(string, ...interface{}) (sql.Result, error)

	SelectFloat(string, ...interface{}) (float64, error)
	SelectInt(string, ...interface{}) (int64, error)
	SelectStr(string, ...interface{}) (string, error)
}

type Transaction interface {
	Insert(...interface{}) error
	Update(...interface{}) (int64, error)
	Delete(...interface{}) (int64, error)
	Get(interface{}, ...interface{}) (interface{}, error)
	Select(interface{}, string, ...interface{}) ([]interface{}, error)
	SelectOne(interface{}, string, ...interface{}) error
	Exec(string, ...interface{}) (sql.Result, error)

	SelectFloat(string, ...interface{}) (float64, error)
	SelectInt(string, ...interface{}) (int64, error)
	SelectStr(string, ...interface{}) (string, error)

	Commit() error
	Rollback() error
}

type DatabaseTableMap interface {
	SetKeys(bool, ...string)
	SetUniqueCompoundKey(...string)
}
