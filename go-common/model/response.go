package model

type Response struct {
	ResponseMessage ResponseMsg  `json:"message"`
	ResponseType    ResponseType `json:"type"`
	ResponseCode    string       `json:"code"`
}

type CommonResponse struct {
	Response
	ResponseData any `json:"data,omitempty"`
}

type CommonErrorResponse struct {
	Response
	InvalidArgs any `json:"invalid_args,omitempty"`
}

type CommonErrorResponseType struct {
	HttpCode            int
	CommonErrorResponse CommonErrorResponse
}

type CommonResponseType struct {
	HttpCode       int
	CommonResponse CommonResponse
}

type ErrorModel struct {
	Error     error        `json:"error"`
	ErrorType ResponseType `json:"error_type"`
}
