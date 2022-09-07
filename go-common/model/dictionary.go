package model

type ResponseType string
type ResponseMsg string

func (r ResponseType) String() string {
	return string(r)
}

// MESSAGE
const (
	// error
	ERR_NO_WALLET_MSG           ResponseMsg = "wallet not found"
	ERR_USER_NO_AUTH_HEADER_MSG ResponseMsg = "no authentication header provided"
	ERR_USER_UNAUTHORIZED_MSG   ResponseMsg = "unauthorized"
	ERR_INTERNAL_MSG            ResponseMsg = "internal server error"
	ERR_BAD_REQUEST_MSG         ResponseMsg = "invalid request parameters"
	ERR_BLOCKED_MSG             ResponseMsg = "request blocked"
	// success
	SUCCESS_OK_MSG       ResponseMsg = "request success"
	SUCCESS_ACCEPTED_MSG ResponseMsg = "request accepted"
)

// TYPE
const (
	// error type
	ERR_NO_WALLET_TYPE      ResponseType = "WALLET_NOT_FOUND"
	ERR_NO_AUTH_HEADER_TYPE ResponseType = "NO_AUTHENTICATION_HEADER"
	ERR_UNAUTHORIZED_TYPE   ResponseType = "UNAUTHORIZED"
	ERR_BLOCKED_TYPE        ResponseType = "BLOCKED"
	ERR_BAD_REQUEST_TYPE    ResponseType = "BAD_REQUEST"
	ERR_INTERNAL_TYPE       ResponseType = "INTERNAL_SERVER_ERROR"
	// success
	SUCCESS_ACCEPTED_TYPE ResponseType = "ACCEPTED"
	SUCCESS_OK_TYPE       ResponseType = "SUCCESS"
)
