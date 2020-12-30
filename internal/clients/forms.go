package clients

import(
	"internal/jobs/types"
)

type updateClientCmd struct {
	ID string `json:"_id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Jobs []types.Job `json:"jobs" binding:"required"`
}

type createClientCmd struct {
	Name string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}
