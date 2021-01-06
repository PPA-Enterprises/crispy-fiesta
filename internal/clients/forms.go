package clients

import(
	"internal/jobs/types"
)

type updateClientCmd struct {
	Name string `json:"name", omitempty`
	Phone string `json:"phone", omitempty`
}

type createClientCmd struct {
	Name string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}
