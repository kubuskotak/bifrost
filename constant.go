package bifrost

const (
	HeaderAuthorization            = "Authorization"
	HeaderContentDisposition       = "Content-Disposition"
	HeaderContentEncoding          = "Content-Encoding"
	HeaderContentLength            = "Content-Length"
	HeaderContentType              = "Content-Type"
	HeaderCookie                   = "Cookie"
	HeaderXCSRFToken               = "X-CSRF-Token"
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	HeaderXTraceId                 = "X-Trace-Id"
	HeaderUberTraceId              = "Uber-Trace-Id"
)

// MIME types
const (
	charsetUTF8                          = "charset=UTF-8"
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextCSV                          = "text/csv"
	MIMETextCSVCharsetUTF8               = MIMETextCSV + "; " + charsetUTF8
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"

	defaultMemory = 32 << 20 // 32 MB
)
