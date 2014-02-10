package database

import "fmt"
import "database/sql"
import "github.com/coopernurse/gorp"
import _ "github.com/go-sql-driver/mysql"
import _ "github.com/lib/pq"
import _ "github.com/mattn/go-sqlite3"
import "github.com/tobz/phosphorus/log"

var schemaMap = make(map[string]interface{})

func RegisterTableSchema(i interface{}, tableName string) {
    if _, ok := schemaMap[tableName]; ok {
        panic(fmt.Sprintf("table schema %s is already defined!", tableName))
    }

    schemaMap[tableName] = i

    log.Server.Info("database", "Registered '%s' schema", tableName)
}

func NewDatabaseConnection(databaseType, databaseDsn string) (*Database, error) {
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

    dbMap := &gorp.DbMap{ Db: databaseConnection, Dialect: dialect }

    // Now run our registered schema objects against the map to bring them in.
    for tableName, tableObject := range schemaMap {
        dbMap.AddTableWithName(tableObject, tableName)
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
