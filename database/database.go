package database

import "github.com/tobz/phosphorus/interfaces"
import "github.com/coopernurse/gorp"

type Database struct {
    dbMap *gorp.DbMap
}

func NewDatabase(dbMap *gorp.DbMap) *Database {
    return &Database{ dbMap: dbMap }
}

func (d *Database) Begin() (interfaces.Transaction, error) {
    return d.dbMap.Begin()
}
