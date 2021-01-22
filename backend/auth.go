package PPA

type AuthToken struct {
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	Token string `json:"token"`
}

type RBAC interface {}
