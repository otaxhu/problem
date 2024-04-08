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
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue",
		Title:  "100 Continue",
		Status: 100,
	},
	101: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-101-switching-protocols",
		Title:  "101 Switching Protocols",
		Status: 101,
	},
	102: {
		Type:   "https://www.rfc-editor.org/rfc/rfc2518.html#section-10.1",
		Title:  "102 Processing",
		Status: 102,
	},
	103: {
		Type:   "https://www.rfc-editor.org/rfc/rfc8297.html#section-2",
		Title:  "103 Early Hints",
		Status: 103,
	},
	200: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok",
		Title:  "200 OK",
		Status: 200,
	},
	201: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-201-created",
		Title:  "201 Created",
		Status: 201,
	},
	202: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-202-accepted",
		Title:  "202 Accepted",
		Status: 202,
	},
	203: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-203-non-authoritative-infor",
		Title:  "203 Non-Authoritative Information",
		Status: 203,
	},
	204: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-204-no-content",
		Title:  "204 No Content",
		Status: 204,
	},
	205: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-205-reset-content",
		Title:  "205 Reset Content",
		Status: 205,
	},
	206: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-206-partial-content",
		Title:  "206 Partial Content",
		Status: 206,
	},
	207: {
		Type:   "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.1",
		Title:  "207 Multi-Status",
		Status: 207,
	},
	208: {
		Type:   "https://www.rfc-editor.org/rfc/rfc5842.html#section-7.1",
		Title:  "208 Already Reported",
		Status: 208,
	},
	226: {
		Type:   "https://www.rfc-editor.org/rfc/rfc3229.html#section-10.4.1",
		Title:  "226 IM Used",
		Status: 226,
	},
	300: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-300-multiple-choices",
		Title:  "300 Multiple Choices",
		Status: 300,
	},
	301: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-301-moved-permanently",
		Title:  "301 Moved Permanently",
		Status: 301,
	},
	302: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-302-found",
		Title:  "302 Found",
		Status: 302,
	},
	303: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-303-see-other",
		Title:  "303 See Other",
		Status: 303,
	},
	304: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-304-not-modified",
		Title:  "304 Not Modified",
		Status: 304,
	},
	305: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-305-use-proxy",
		Title:  "305 Use Proxy",
		Status: 305,
	},
	307: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-307-temporary-redirect",
		Title:  "307 Temporary Redirect",
		Status: 307,
	},
	308: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-308-permanent-redirect",
		Title:  "308 Permanent Redirect",
		Status: 308,
	},
	400: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request",
		Title:  "400 Bad Request",
		Status: 400,
	},
	401: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized",
		Title:  "401 Unauthorized",
		Status: 401,
	},
	402: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-402-payment-required",
		Title:  "402 Payment Required",
		Status: 402,
	},
	403: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden",
		Title:  "403 Forbidden",
		Status: 403,
	},
	404: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found",
		Title:  "404 Not Found",
		Status: 404,
	},
	405: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-405-method-not-allowed",
		Title:  "405 Method Not Allowed",
		Status: 405,
	},
	406: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-406-not-acceptable",
		Title:  "406 Not Acceptable",
		Status: 406,
	},
	407: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-407-proxy-authentication-re",
		Title:  "407 Proxy Authentication Required",
		Status: 407,
	},
	408: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-408-request-timeout",
		Title:  "408 Request Timeout",
		Status: 408,
	},
	409: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-409-conflict",
		Title:  "409 Conflict",
		Status: 409,
	},
	410: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-410-gone",
		Title:  "410 Gone",
		Status: 410,
	},
	411: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-411-length-required",
		Title:  "411 Length Required",
		Status: 411,
	},
	412: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-412-precondition-failed",
		Title:  "412 Precondition Failed",
		Status: 412,
	},
	413: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-413-content-too-large",
		Title:  "413 Content Too Large",
		Status: 413,
	},
	414: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-414-uri-too-long",
		Title:  "414 URI Too Long",
		Status: 414,
	},
	415: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-415-unsupported-media-type",
		Title:  "415 Unsupported Media Type",
		Status: 415,
	},
	416: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-416-range-not-satisfiable",
		Title:  "416 Range Not Satisfiable",
		Status: 416,
	},
	417: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-417-expectation-failed",
		Title:  "417 Expectation Failed",
		Status: 417,
	},
	418: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-418-unused",
		Title:  "418 I'm a teapot",
		Status: 418,
	},
	421: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request",
		Title:  "421 Misdirected Request",
		Status: 421,
	},
	422: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content",
		Title:  "421 Unprocessable Content",
		Status: 422,
	},
	423: {
		Type:   "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.3",
		Title:  "423 Locked",
		Status: 423,
	},
	424: {
		Type:   "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.4",
		Title:  "424 Failed Dependency",
		Status: 424,
	},
	425: {
		Type:   "https://www.rfc-editor.org/rfc/rfc8470.html#section-5.2",
		Title:  "425 Too Early",
		Status: 425,
	},
	426: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-426-upgrade-required",
		Title:  "426 Upgrade Required",
		Status: 426,
	},
	428: {
		Type:   "https://www.rfc-editor.org/rfc/rfc6585.html#section-3",
		Title:  "428 Precondition Required",
		Status: 428,
	},
	429: {
		Type:   "https://www.rfc-editor.org/rfc/rfc6585.html#section-4",
		Title:  "429 Too Many Requests",
		Status: 429,
	},
	431: {
		Type:   "https://www.rfc-editor.org/rfc/rfc6585.html#section-5",
		Title:  "431 Request Header Fields Too Large",
		Status: 431,
	},
	451: {
		Type:   "https://www.rfc-editor.org/rfc/rfc7725.html#section-3",
		Title:  "451 Unavailable For Legal Reasons",
		Status: 451,
	},
	500: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error",
		Title:  "500 Internal Server Error",
		Status: 500,
	},
	501: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-501-not-implemented",
		Title:  "501 Not Implemented",
		Status: 501,
	},
	502: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway",
		Title:  "502 Bad Gateway",
		Status: 502,
	},
	503: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-503-service-unavailable",
		Title:  "503 Service Unavailable",
		Status: 503,
	},
	504: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-504-gateway-timeout",
		Title:  "504 Gateway Timeout",
		Status: 504,
	},
	505: {
		Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-505-http-version-not-suppor",
		Title:  "505 HTTP Version Not Supported",
		Status: 505,
	},
	506: {
		Type:   "https://www.rfc-editor.org/rfc/rfc2295.html#section-8.1",
		Title:  "506 Variant Also Negotiates",
		Status: 506,
	},
	507: {
		Type:   "https://www.rfc-editor.org/rfc/rfc4918.html#section-11.5",
		Title:  "507 Insufficient Storage",
		Status: 507,
	},
	508: {
		Type:   "https://www.rfc-editor.org/rfc/rfc5842.html#section-7.2",
		Title:  "508 Loop Detected",
		Status: 508,
	},
	510: {
		Type:   "https://www.rfc-editor.org/rfc/rfc2774.html#section-7",
		Title:  "510 Not Extended",
		Status: 510,
	},
	511: {
		Type:   "https://www.rfc-editor.org/rfc/rfc6585.html#section-6",
		Title:  "511 Network Authentication Required",
		Status: 511,
	},
}
