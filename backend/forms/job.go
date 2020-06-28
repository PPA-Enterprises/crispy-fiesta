package forms

type SubmitJobCmd struct {
	ClientInfo string `json:"client_info" binding:"required"`
	CarInfo string `json:"car_info" binding:"required"`
	AppointmentInfo string `json:"appointment_info" binding:"required"`
}
