package header

const (
	TextContentType = "text/plain"
	JSONContentType = "application/json"

	XContentTypeOptionsNosniff = "nosniff"
)

const (
	AcceptHeader        = "Accept"
	AcceptEncoding      = "Accept-Encoding"
	AuthorizationHeader = "Authorization"

	OriginHeader                    = "Origin"
	AccessControlAllowOriginHeader  = "Access-Control-Allow-Origin"
	AccessControlAllowMethods       = "Access-Control-Allow-Methods"
	AccessControlAllowHeadersHeader = "Access-Control-Allow-Headers"

	ContentLength             = "Content-Length"
	ContentType               = "Content-Type"
	RequestHeader             = "X-Request-ID"
	XContentTypeOptionsHeader = "X-Content-Type-Options"
	XCSRFToken                = "X-CSRF-Token"
)
