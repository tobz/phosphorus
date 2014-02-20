package database

import "database/sql"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/coopernurse/gorp"

type Database struct {
	dbMap *gorp.DbMap
}

func NewDatabase(dbMap *gorp.DbMap) *Database {
	return &Database{dbMap: dbMap}
}

func (d *Database) Begin() (interfaces.Transaction, error) {
	return d.dbMap.Begin()
}

func (d *Database) Insert(list ...interface{}) error {
	return d.dbMap.Insert(list...)
}

func (d *Database) Update(list ...interface{}) (int64, error) {
	return d.dbMap.Update(list...)
}

func (d *Database) Delete(list ...interface{}) (int64, error) {
	return d.dbMap.Delete(list...)
}

func (d *Database) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return d.dbMap.Get(i, keys...)
}

func (d *Database) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return d.dbMap.Select(i, query, args...)
}

func (d *Database) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return d.dbMap.SelectOne(holder, query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.dbMap.Exec(query, args...)
}

func (d *Database) SelectFloat(query string, args ...interface{}) (float64, error) {
	return d.dbMap.SelectFloat(query, args...)
}

func (d *Database) SelectInt(query string, args ...interface{}) (int64, error) {
	return d.dbMap.SelectInt(query, args...)
}

func (d *Database) SelectStr(query string, args ...interface{}) (string, error) {
	return d.dbMap.SelectStr(query, args...)
}

type DatabaseTableMap struct {
	tableMap *gorp.TableMap
}

func NewDatabaseTableMap(tableMap *gorp.TableMap) *DatabaseTableMap {
	return &DatabaseTableMap{tableMap: tableMap}
}

func (dtm *DatabaseTableMap) SetKeys(autoIncrement bool, columns ...string) {
	dtm.tableMap.SetKeys(autoIncrement, columns...)
}

func (dtm *DatabaseTableMap) SetUniqueCompoundKey(columns ...string) {
	dtm.tableMap.SetUniqueTogether(columns...)
}
