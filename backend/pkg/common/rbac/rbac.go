package rbac
import (
	"PPA"
)

type Service struct{}

func boolToErr(b bool) error {
	if b {
		return nil
	}
	return PPA.Forbidden
}
