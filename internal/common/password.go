package common

import (
	"github.com/matthewhartstonge/argon2"
)
var argon2Config argon2.Config = argon2.DefaultConfig()

func HashFromPlaintext(password string) (string, error) {
	raw, err := argon2Config.Hash([]byte(password), nil)
	if err != nil {
		return string(""), err
	}
	return string(raw.Encode()), nil
}

func VerifyPassword(password string, hashedPassword string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(hashedPassword))
}
