package problem

// Used in AddExtensionMember function, to know if the provided key is one of these.
var standardMembers = map[string]struct{}{
	"type":     {},
	"status":   {},
	"title":    {},
	"detail":   {},
	"instance": {},
}

// Default values for a problem details response
//
// https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
var problems = map[int]Problem{
	100: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue",
		Title: "100 Continue",
	},
	101: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-101-switching-protocols",
		Title: "101 Switching Protocols",
	},
	102: {
		Type:  "https://www.rfc-editor.org/rfc/rfc2518.html#section-10.1",
		Title: "102 Processing",
	},
	103: {
		Type:  "https://www.rfc-editor.org/rfc/rfc8297.html#section-2",
		Title: "103 Early Hints",
	},
	200: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok",
		Title: "200 OK",
	},
	201: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-201-created",
		Title: "201 Created",
	},
	202: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-202-accepted",
		Title: "202 Accepted",
	},
	203: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-203-non-authoritative-infor",
		Title: "203 Non-Authoritative Information",
	},
	204: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-204-no-content",
		Title: "204 No Content",
	},
	205: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-205-reset-content",
		Title: "205 Reset Content",
	},
	206: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-206-partial-content",
		Title: "206 Partial Content",
	},
	207: {
		Type:  "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.1",
		Title: "207 Multi-Status",
	},
	208: {
		Type:  "https://www.rfc-editor.org/rfc/rfc5842.html#section-7.1",
		Title: "208 Already Reported",
	},
	226: {
		Type:  "https://www.rfc-editor.org/rfc/rfc3229.html#section-10.4.1",
		Title: "226 IM Used",
	},
	300: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-300-multiple-choices",
		Title: "300 Multiple Choices",
	},
	301: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-301-moved-permanently",
		Title: "301 Moved Permanently",
	},
	302: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-302-found",
		Title: "302 Found",
	},
	303: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-303-see-other",
		Title: "303 See Other",
	},
	304: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-304-not-modified",
		Title: "304 Not Modified",
	},
	305: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-305-use-proxy",
		Title: "305 Use Proxy",
	},
	307: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-307-temporary-redirect",
		Title: "307 Temporary Redirect",
	},
	308: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-308-permanent-redirect",
		Title: "308 Permanent Redirect",
	},
	400: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request",
		Title: "400 Bad Request",
	},
	401: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized",
		Title: "401 Unauthorized",
	},
	402: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-402-payment-required",
		Title: "402 Payment Required",
	},
	403: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden",
		Title: "403 Forbidden",
	},
	404: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found",
		Title: "404 Not Found",
	},
	405: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-405-method-not-allowed",
		Title: "405 Method Not Allowed",
	},
	406: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-406-not-acceptable",
		Title: "406 Not Acceptable",
	},
	407: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-407-proxy-authentication-re",
		Title: "407 Proxy Authentication Required",
	},
	408: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-408-request-timeout",
		Title: "408 Request Timeout",
	},
	409: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-409-conflict",
		Title: "409 Conflict",
	},
	410: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-410-gone",
		Title: "410 Gone",
	},
	411: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-411-length-required",
		Title: "411 Length Required",
	},
	412: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-412-precondition-failed",
		Title: "412 Precondition Failed",
	},
	413: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-413-content-too-large",
		Title: "413 Content Too Large",
	},
	414: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-414-uri-too-long",
		Title: "414 URI Too Long",
	},
	415: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-415-unsupported-media-type",
		Title: "415 Unsupported Media Type",
	},
	416: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-416-range-not-satisfiable",
		Title: "416 Range Not Satisfiable",
	},
	417: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-417-expectation-failed",
		Title: "417 Expectation Failed",
	},
	418: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-418-unused",
		Title: "418 I'm a teapot",
	},
	421: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request",
		Title: "421 Misdirected Request",
	},
	422: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content",
		Title: "421 Unprocessable Content",
	},
	423: {
		Type:  "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.3",
		Title: "423 Locked",
	},
	424: {
		Type:  "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.4",
		Title: "424 Failed Dependency",
	},
	425: {
		Type:  "https://www.rfc-editor.org/rfc/rfc8470.html#section-5.2",
		Title: "425 Too Early",
	},
	426: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-426-upgrade-required",
		Title: "426 Upgrade Required",
	},
	428: {
		Type:  "https://www.rfc-editor.org/rfc/rfc6585.html#section-3",
		Title: "428 Precondition Required",
	},
	429: {
		Type:  "https://www.rfc-editor.org/rfc/rfc6585.html#section-4",
		Title: "429 Too Many Requests",
	},
	431: {
		Type:  "https://www.rfc-editor.org/rfc/rfc6585.html#section-5",
		Title: "431 Request Header Fields Too Large",
	},
	451: {
		Type:  "https://www.rfc-editor.org/rfc/rfc7725.html#section-3",
		Title: "451 Unavailable For Legal Reasons",
	},
	500: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error",
		Title: "500 Internal Server Error",
	},
	501: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-501-not-implemented",
		Title: "501 Not Implemented",
	},
	502: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway",
		Title: "502 Bad Gateway",
	},
	503: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-503-service-unavailable",
		Title: "503 Service Unavailable",
	},
	504: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-504-gateway-timeout",
		Title: "504 Gateway Timeout",
	},
	505: {
		Type:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-505-http-version-not-suppor",
		Title: "505 HTTP Version Not Supported",
	},
	506: {
		Type:  "https://www.rfc-editor.org/rfc/rfc2295.html#section-8.1",
		Title: "506 Variant Also Negotiates",
	},
	507: {
		Type:  "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.5",
		Title: "507 Insufficient Storage",
	},
	508: {
		Type:  "https://www.rfc-editor.org/rfc/rfc5842.html#section-7.2",
		Title: "508 Loop Detected",
	},
	510: {
		Type:  "https://www.rfc-editor.org/rfc/rfc2774.html#section-7",
		Title: "510 Not Extended",
	},
	511: {
		Type:  "https://www.rfc-editor.org/rfc/rfc6585.html#section-6",
		Title: "511 Network Authentication Required",
	},
}
