package secure

import (
	"github.com/matthewhartstonge/argon2"

)

// New initializes security service
func New() *Service {
	cfg := argon2.DefaultConfig()
	return &Service{cfg: cfg}
}

// Service holds security related methods
type Service struct {
	cfg		argon2.Config
}


// Hash hashes the password using bcrypt
func (self *Service) Hash(password string) string {
	raw, _ := self.cfg.Hash([]byte(password), nil)
	return string(raw.Encode())
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func (*Service) HashMatchesPassword(hash, password string) bool {
	res, _ := argon2.VerifyEncoded([]byte(password), []byte(hash))
	return res
}
