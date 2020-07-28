module github.com/PPA-Enterprises/crispy-fiesta

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/matthewhartstonge/argon2 v0.1.4 // indirect
	go.mongodb.org/mongo-driver v1.3.5
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	internal/clients v1.0.0
	internal/common v1.0.0
	internal/db v1.0.0
	internal/jobs v1.0.0
	internal/uid v1.0.0
	internal/users v1.0.0
)

replace internal/common => ./internal/common

replace internal/db => ./internal/db

replace internal/users => ./internal/users

replace internal/jobs => ./internal/jobs

replace internal/uid => ./internal/uid

replace internal/clients => ./internal/clients
