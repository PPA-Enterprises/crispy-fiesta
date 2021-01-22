module PPA

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/matthewhartstonge/argon2 v0.1.4 // indirect
	go.mongodb.org/mongo-driver v1.4.4
	gopkg.in/yaml.v2 v2.4.0 // indirect
	pkg v1.0.0
)

replace pkg => ./pkg
