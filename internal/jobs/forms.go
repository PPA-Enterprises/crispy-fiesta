package jobs

type submitJobCmd struct {
	ClientName      string `json:"client_info" binding:"required"`
	CarInfo         string `json:"car_info" binding:"required"`
	AppointmentInfo string `json:"appointment_info" binding:"required"`
	Notes           string `json:"notes"binding:"required"`
}

type updateJobCmd struct {
	ID              string `json:"_id" binding:"required"`
	ClientInfo      string `json:"client_info" binding:"required"`
	CarInfo         string `json:"car_info" binding:"required"`
	AppointmentInfo string `json:"appointment_info" binding:"required"`
	Notes           string `json:"notes"binding:"required"`
}
