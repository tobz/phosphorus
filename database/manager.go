package database

import "fmt"
import "database/sql"
import "github.com/coopernurse/gorp"
import _ "github.com/go-sql-driver/mysql"
import _ "github.com/lib/pq"
import _ "github.com/mattn/go-sqlite3"
import "github.com/tobz/phosphorus/log"
import "github.com/tobz/phosphorus/interfaces"

type tableSchemaHandler struct {
	tableObject      interface{}
	tableMapModifier func(interfaces.DatabaseTableMap)
}

var tableMap = make(map[string]*tableSchemaHandler)

func RegisterTableSchema(tableObject interface{}, tableName string, tableMapModifier func(interfaces.DatabaseTableMap)) {
	if _, ok := tableMap[tableName]; ok {
		panic(fmt.Sprintf("Table '%s' is already defined!", tableName))
	}

	tableMap[tableName] = &tableSchemaHandler{tableObject: tableObject, tableMapModifier: tableMapModifier}

	log.Server.Info("database", "Registered schema for '%s' table...", tableName)
}

func NewDatabaseConnection(databaseType, databaseDsn string) (interfaces.Database, error) {
	// Set up our driver and dialect before handing over to modl.
	dialect, err := getDatabaseDialect(databaseType)
	if err != nil {
		return nil, err
	}

	// Now set up our connection.
	databaseConnection, err := sql.Open(databaseType, databaseDsn)
	if err != nil {
		return nil, err
	}

	dbMap := &gorp.DbMap{Db: databaseConnection, Dialect: dialect}

	// Now run our registered schema objects against the database map to bring them in.
	for tableName, tableHandler := range tableMap {
		tableMap := dbMap.AddTableWithName(tableHandler.tableObject, tableName)
		if tableHandler.tableMapModifier != nil {
			tableHandler.tableMapModifier(NewDatabaseTableMap(tableMap))
		}
	}

	// Make sure our schema exists.
	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		return nil, fmt.Errorf("failed to ensure schema exists: %s", err)
	}

	log.Server.Info("database", "New database connection created; using '%s'", databaseType)

	return NewDatabase(dbMap), nil
}

func getDatabaseDialect(databaseType string) (gorp.Dialect, error) {
	switch databaseType {
	case "mysql":
		return gorp.MySQLDialect{"InnoDB", "UTF-8"}, nil
	case "postgres":
		return gorp.PostgresDialect{}, nil
	case "sqlite3":
		return gorp.SqliteDialect{}, nil
	}

	return nil, fmt.Errorf("couldn't find matching dialect")
}
