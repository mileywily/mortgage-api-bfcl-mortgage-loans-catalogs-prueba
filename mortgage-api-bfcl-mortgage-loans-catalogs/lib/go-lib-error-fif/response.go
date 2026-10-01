package error_fif

type Status string

const (
	NO_CONTENT            Status = "NO_CONTENT"
	BAD_REQUEST           Status = "INVALID_ARGUMENT"
	UNAUTHORIZED          Status = "UNAUTHORIZED"
	FORBIDDEN             Status = "FORBIDDEN"
	NOT_FOUND             Status = "NOT_FOUND"
	NOT_ACCEPTABLE        Status = "NOT_ACCEPTABLE"
	UNPROCESSABLE_ENTITY  Status = "UNPROCESSABLE_ENTITY"
	PRECONDITION_REQUIRED Status = "PRECONDITION_REQUIRED"
	INTERNAL_ERROR        Status = "INTERNAL_ERROR"
	BAD_GATEWAY           Status = "BAD_GATEWAY"
	GATEWAY_TIMEOUT       Status = "GATEWAY_TIMEOUT"
)

type ResponseBody struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Status  Status        `json:"status"`
	TraceId string        `json:"traceId,omitempty"`
	Details []interface{} `json:"details,omitempty"`
}

type FieldViolation struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}
