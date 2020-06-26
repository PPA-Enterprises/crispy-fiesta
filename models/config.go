package models

import (
	"os"

	"github.com/PPA-Enterprises/crispy-fiesta/db"
)

var server = os.Getenv("DATABASE")

var databaseName = os.Getenv("DATABASE_NAME")

var dbConnect = db.NewConnection(server, databaseName)
