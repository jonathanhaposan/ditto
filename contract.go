package main

type Contract struct {
	HttpRequest                  *HTTPRequest                  `json:"httpRequest"`
	HttpResponse                 *HTTPResponse                 `json:"httpResponse"`
	HttpOverrideForwardedRequest *HttpOverrideForwardedRequest `json:"httpOverrideForwardedRequest"`
}

type HTTPRequest struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Headers     map[string]string `json:"headers"`
	Body        Body              `json:"body"`
	Protocol    string            `json:"protocol"`
	Host        string            `json:"host"`
	QueryString map[string]string `json:"queryString"`
}

type HTTPResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Body struct {
	JSON string `json:"json"`
}

type HttpOverrideForwardedRequest struct {
	HttpRequest HTTPRequest `json:"httpRequest"`
}
