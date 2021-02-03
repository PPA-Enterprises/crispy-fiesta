package secure

import (
	"fmt"
	"hash"
	"time"
	"strconv"
	"github.com/matthewhartstonge/argon2"

)

func New(h hash.Hash) *Service {
	cfg := argon2.DefaultConfig()
	return &Service{cfg: cfg, hasher: h}
}

type Service struct {
	cfg		argon2.Config
	hasher hash.Hash
}

func (self *Service) Hash(password string) string {
	raw, _ := self.cfg.Hash([]byte(password), nil)
	return string(raw.Encode())
}

func (*Service) HashMatchesPassword(hash, password string) bool {
	res, _ := argon2.VerifyEncoded([]byte(password), []byte(hash))
	return res
}

func (s *Service) Token(str string) string {
	s.hasher.Reset()
	fmt.Fprintf(s.hasher, "%s%s", str, strconv.Itoa(time.Now().Nanosecond()))
	return fmt.Sprintf("%x", s.hasher.Sum(nil))
}
