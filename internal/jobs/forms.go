package jobs

type submitJobCmd struct {
	ClientName      string `json:"client_name" binding:"required"`
	ClientPhone     string `json:"client_phone" binding:"required"`
	CarInfo         string `json:"car_info" binding:"required"`
	AppointmentInfo string `json:"appointment_info" binding:"required"`
	Notes           string `json:"notes"binding:"required"`
}

type updateJobCmd struct {
	ID              string `json:"_id" binding:"required"`
	ClientName      string `json:"client_name" omitempty`
	ClientPhone     string `json:"client_phone" omitempty`
	CarInfo         string `json:"car_info" omitempty`
	AppointmentInfo string `json:"appointment_info" omitempty`
	Notes           string `json:"notes" omitempty`
}
