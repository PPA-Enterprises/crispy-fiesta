package models

import (
	"github.com/PPA-Enterprises/crispy-fiesta/db"
)

var server = "mongodb://localhost:27017"

var databaseName = "PPA"

var dbConnect = db.NewConnection(server)
