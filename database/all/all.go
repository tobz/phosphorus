package database

import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/database"
import "github.com/tobz/phosphorus/database/models"

func init() {
	database.RegisterTableSchema(models.Account{}, "accounts", func(dtm interfaces.DatabaseTableMap) {
		dtm.SetKeys(true, "account_id")
	})

	database.RegisterTableSchema(models.Character{}, "characters", func(dtm interfaces.DatabaseTableMap) {
		dtm.SetKeys(true, "character_id")
	})
}
