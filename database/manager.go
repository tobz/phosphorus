package database

import "fmt"
import "database/sql"
import "github.com/jmoiron/modl"
import _ "github.com/go-sql-driver/mysql"
import _ "github.com/lib/pq"
import _ "github.com/mattn/go-sqlite3"
import "github.com/tobz/phosphorus/interfaces"

type tableMetadata struct {
    obj interface{}
    post func(interfaces.DatabaseTableMap)
}
var schemaMap = make(map[string]tableMetadata)

func RegisterTableSchema(i interface{}, tableName string, post func(interfaces.DatabaseTableMap)) {
    if _, ok := schemaMap[tableName]; ok {
        panic(fmt.Sprintf("table schema %s is already defined!", tableName))
    }

    schemaMap[tableName] = tableMetadata{i, post}
}

func NewDatabaseConnection(c interfaces.Config) (interfaces.Database, error) {
    // Set up our driver and dialect before handing over to modl.
    databaseType, err := c.GetAsString("database/type")
    if err != nil {
        return nil, err
    }

    dialect, err := getDatabaseDialect(databaseType)
    if err != nil {
        return nil, err
    }

    // Now set up our connection.
    databaseDsn, err := c.GetAsString("database/dsn")
    if err != nil {
        return nil, err
    }

    databaseConnection, err := sql.Open(databaseType, databaseDsn)
    if err != nil {
        return nil, err
    }

    dbMap := modl.NewDbMap(databaseConnection, dialect)

    // Now run our registered schema objects against the map to bring them in.
    for tableName, tableMetadata := range schemaMap {
        tableMap := dbMap.AddTableWithName(tableMetadata.obj, tableName)
        if tableMetadata.post != nil {
            tableMetadata.post(tableMap)
        }
    }

    return dbMap, nil
}

func getDatabaseDialect(databaseType string) (modl.Dialect, error) {
    switch databaseType {
    case "mysql":
        return modl.MySQLDialect{"InnoDB", "UTF-8"}, nil
    case "postgres":
        return modl.PostgresDialect{}, nil
    case "sqlite3":
        return modl.SqliteDialect{}, nil
    }

    return nil, fmt.Errorf("couldn't find matching dialect")
}
