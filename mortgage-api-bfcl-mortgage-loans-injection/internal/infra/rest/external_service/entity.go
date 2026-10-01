package external_service

type ExternalServiceRequest struct {
	FieldExternalService string `json:"fieldExternalService"`
}

type ExternalServiceResponse struct {
	Field string `json:"field"`
}
