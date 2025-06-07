package main

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Header  map[string]string
	Body    string
}

type HTTPResponse struct {
	Version       string
	StatusCode    int
	StatusMessage string
	Header        map[string]string
	Body          string
}
