package enums

type TravelRequestStatus string

const (
	TravelRequestStatusSolicited TravelRequestStatus = "SOLICITED"
	TravelRequestStatusApproved  TravelRequestStatus = "APPROVED"
	TravelRequestStatusCanceled  TravelRequestStatus = "CANCELED"
)
