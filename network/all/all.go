package all

// This package imports and registers all packet handlers.
import (
	_ "github.com/tobz/phosphorus/network/handlers"
	_ "github.com/tobz/phosphorus/network/handlers/badnamecheck"
)
