// Package httphdr contains the names of HTTP headers.
//
// Please keep the values in their canonical form.
package httphdr

// Common standard headers for caching.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#caching.
const (
	CacheControl = "Cache-Control"
)

// Common standard headers for cookie management.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#cookies.
const (
	Cookie    = "Cookie"
	SetCookie = "Set-Cookie"
)

// Common standard headers for conditionals.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#conditionals.
const (
	IfMatch           = "If-Match"
	IfModifiedSince   = "If-Modified-Since"
	IfNoneMatch       = "If-None-Match"
	IfRange           = "If-Range"
	IfUnmodifiedSince = "If-Unmodified-Since"
	LastModified      = "Last-Modified"
	Vary              = "Vary"
)

// Common standard headers for content negotiation.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#content_negotiation.
const (
	Accept         = "Accept"
	AcceptEncoding = "Accept-Encoding"
	AcceptLanguage = "Accept-Language"
)

// Common standard headers for CORS.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#cors.
const (
	AccessControlAllowOrigin = "Access-Control-Allow-Origin"
	Origin                   = "Origin"
)

// Common standard headers for controlling downloads.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#downloads.
const (
	ContentDisposition = "Content-Disposition"
)

// Common standard headers for fetch metadata.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#fetch_metadata_request_headers.
const (
	SecFetchDest = "Sec-Fetch-Dest"
)

// Common standard headers for message body information.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#message_body_information.
const (
	ContentEncoding = "Content-Encoding"
	ContentLength   = "Content-Length"
	ContentType     = "Content-Type"
)

// Common standard headers for request contexts.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#request_context.
const (
	Host      = "Host"
	Referer   = "Referer"
	UserAgent = "User-Agent"
)

// Common standard headers for response contexts.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#response_context.
const (
	Server = "Server"
)

// Common standard headers for security.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#security.
const (
	ContentSecurityPolicy           = "Content-Security-Policy"
	ContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	StrictTransportSecurity         = "Strict-Transport-Security"
)

// Common standard headers for transfer coding.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#transfer_coding.
const (
	TransferEncoding = "Transfer-Encoding"
	Trailer          = "Trailer"
)

// Miscellaneous common standard headers.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#other.
const (
	RetryAfter = "Retry-After"
	AltSvc     = "Alt-Svc"
)

// Common extension headers.
const (
	AdminToken   = "Admin-Token"
	TrueClientIP = "True-Client-IP"

	XError        = "X-Error"
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-Ip"
	XRequestID    = "X-Request-Id"
)

// Common Cloudflare extension headers.
const (
	CFConnectingIP = "Cf-Connecting-Ip"
	CFIPCountry    = "Cf-Ipcountry"
)
