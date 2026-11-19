package authorize

type ErrorType string

const (
	ErrorNone                ErrorType = ""
	ErrorInvalidRequest      ErrorType = "invalid_request"
	ErrorUnsupportedResponse ErrorType = "unsupported_response_type"
	ErrorInvalidScope        ErrorType = "invalid_scope"
	ErrorInvalidClient       ErrorType = "invalid_client"
	ErrorInvalidRedirectURI  ErrorType = "invalid_redirect_uri"
	ErrorAccessDenied        ErrorType = "access_denied"
	ErrorServerError         ErrorType = "server_error"
)

func (e ErrorType) String() string {
	return string(e)
}
