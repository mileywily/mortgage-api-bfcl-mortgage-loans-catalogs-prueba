package create

type RequestDto struct {
	Header HeaderDto
	Uri    UriDto
	Body   BodyRequestDto
}

type HeaderDto struct {
	XField string `header:"X-Field" binding:"required"`
}

type UriDto struct {
	Something string `uri:"something" binding:"required"`
}

type BodyRequestDto struct {
	Body BodyDto `json:"body" binding:"required"`
}

type BodyDto struct {
	Field2 string `json:"field2" binding:"required"`
}

type ResponseDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
