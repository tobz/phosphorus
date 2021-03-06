package all

// This package imports and registers all packet handlers.
import (
	_ "github.com/tobz/phosphorus/network/handlers"
	_ "github.com/tobz/phosphorus/network/handlers/character"
	_ "github.com/tobz/phosphorus/network/handlers/cryptkey"
	_ "github.com/tobz/phosphorus/network/handlers/loginrequest"
	_ "github.com/tobz/phosphorus/network/handlers/ping"
	_ "github.com/tobz/phosphorus/network/handlers/realmselection"
)
