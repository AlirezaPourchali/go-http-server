package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Header  map[string]string
	Body    []byte
}

type HTTPResponse struct {
	Version       string
	StatusCode    int
	StatusMessage string
	Header        map[string]string
	Body          []byte
}

type HTTPServer struct {
	Encoding []string
}

func NewHTTPServer(encode []string) HTTPServer {
	s := HTTPServer{
		encode,
	}
	return s
}

func NewHTTPResponse(req HTTPRequest, s HTTPServer, config ...HTTPResponse) HTTPResponse {
	// Set defaults
	response := HTTPResponse{
		Version:       "HTTP/1.1",
		StatusCode:    200,
		StatusMessage: "OK",
		Header:        map[string]string{"Server": "AlirezaPourchali"},
		Body:          nil,
	}

	// Override with provided config
	if len(config) > 0 {
		cfg := config[0]

		if cfg.Version != "" {
			response.Version = cfg.Version
		}
		if cfg.StatusCode != 0 {
			response.StatusCode = cfg.StatusCode
		}
		if cfg.StatusMessage != "" {
			response.StatusMessage = cfg.StatusMessage
		}
		if cfg.Header != nil {
			response.Header = cfg.Header
		}
		if cfg.Body != nil {
			response.Body = cfg.Body
		}
	}

	if v, ok := req.Header["accept-encoding"]; ok {
		// Accept-Encoding: encoding-1, encoding-2, encoding-3
		encodings := strings.Split(v, ",")
		var h string
		for i := range s.Encoding {
			for j := range encodings {
				if s.Encoding[i] == strings.TrimSpace(encodings[j]) {
					if len(h) == 0 {
						h += s.Encoding[i]
					} else {
						h += s.Encoding[i] + ", "
					}
				}
			}
		}
		if h != "" {
			response.Header["Content-Encoding"] = h
		}
	}

	return response
}

func HTTPRespond(r HTTPResponse, c net.Conn) {

	if encoding, ok := r.Header["Content-Encoding"]; ok {
		if encoding == "gzip" {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			_, err := gz.Write(r.Body)
			gz.Close()
			if err != nil {
				log.Fatal(err)
			} else {
				r.Body = b.Bytes()
			}
			r.Header["Content-Length"] = strconv.Itoa(len(r.Body))
		}
	}
	var h string
	for key, value := range r.Header {
		h += key + ": " + value + "\r\n"
	}

	fmt.Fprintf(c, "%v %v %v\r\n%v\r\n%s", r.Version, r.StatusCode, r.StatusMessage, h, r.Body)
}
