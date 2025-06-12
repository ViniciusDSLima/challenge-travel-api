package enums

type TravelRequestStatus string

const (
	TravelRequestStatusSolicited TravelRequestStatus = "SOLICITED"
	TravelRequestStatusApproved  TravelRequestStatus = "APPROVED"
	TravelRequestStatusRejected  TravelRequestStatus = "REJECTED"
)
