package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func AcceptConnection(Listener net.Listener, ConnectionChannel chan net.Conn) {
	conn, err := Listener.Accept()
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	ConnectionChannel <- conn
}

func handleRequest(TCPConnection net.Conn, Server HTTPServer, ServerDir *string) {
	var RequestLine string
	var Request HTTPRequest
	ReqBuffer := bufio.NewReader(TCPConnection)
	RequestLine, _ = ReqBuffer.ReadString('\n')

	Request.Method = strings.Fields(RequestLine)[0]
	Request.Path = strings.Fields(RequestLine)[1]
	Request.Version = strings.Fields(RequestLine)[2]

	Request.Header = make(map[string]string)
	// parsing headers
	for {
		header, err := ReqBuffer.ReadString('\n')
		if err != nil {
			fmt.Fprintf(TCPConnection, "HTTP/1.1 400 Bad Request\r\n\r\n")
			TCPConnection.Close()
			return
		}
		if header == "\r\n" {
			break
		}

		key := strings.SplitAfterN(header, ":", 2)[0]
		key = strings.Trim(key, ":")
		value := strings.SplitAfterN(header, ":", 2)[1]
		value = strings.TrimSpace(value)
		Request.Header[strings.ToLower(key)] = strings.ToLower(value)
	}

	// read the body
	if len, ok := Request.Header["content-length"]; ok {
		var body []byte
		l, err := strconv.Atoi(len)
		fmt.Println(l)
		if err == nil {
			for i := 0; i < l; i++ {
				b, err := ReqBuffer.ReadByte()
				if err == nil {
					body = append(body, b)
				}
			}
			Request.Body = body
		}
	}
	// http path handling for GET Method
	RequestPath := Request.Path
	RequestMethod := Request.Method
	if RequestMethod == "GET" {
		if RequestPath == "/" {
			r := NewHTTPResponse(Request, Server, HTTPResponse{
				Version: Request.Version,
			})
			// response generator
			HTTPRespond(r, TCPConnection)
			// /echo/{str}
		} else if strings.Split(RequestPath, "/")[1] == "echo" {
			echo_string := strings.Split(RequestPath, "/")[2]
			m := mimetype.Detect([]byte(echo_string))
			mime := strings.Split(m.String(), ";")[0]
			r := NewHTTPResponse(Request, Server, HTTPResponse{
				Version: Request.Version,
				Header:  map[string]string{"Content-Type": mime, "Content-Length": strconv.Itoa(len(echo_string))},
				Body:    []byte(echo_string),
			})

			HTTPRespond(r, TCPConnection)
		} else if RequestPath == "/user-agent" {
			UserAgent := Request.Header["user-agent"]
			m := mimetype.Detect([]byte(UserAgent))
			mime := strings.Split(m.String(), ";")[0]
			r := NewHTTPResponse(Request, Server, HTTPResponse{
				Version: Request.Version,
				Header:  map[string]string{"Content-Type": mime, "Content-Length": strconv.Itoa(len(UserAgent))},
				Body:    []byte(UserAgent),
			})

			// fmt.Fprintf(TCPConnection, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(UserAgent), UserAgent)
			// using the function instead
			HTTPRespond(r, TCPConnection)
			//files endpoint
		} else if strings.Split(RequestPath, "/")[1] == "files" {
			FileName := strings.Split(RequestPath, "/")[2]
			FileContent, err := os.ReadFile(*ServerDir + "/" + FileName)
			if err == nil {
				m := mimetype.Detect([]byte(FileContent))
				mime := strings.Split(m.String(), ";")[0]
				r := NewHTTPResponse(Request, Server, HTTPResponse{
					Version: Request.Version,
					Header:  map[string]string{"Content-Type": mime, "Content-Length": strconv.Itoa(len(FileContent))},
					Body:    []byte(FileContent),
				})
				// fmt.Fprintf(TCPConnection, "HTTP/1.1 200 OK\r\nContent-Type: application/octet\r\nContent-Length: %v\r\n\r\n%s", len, FileContent)
				HTTPRespond(r, TCPConnection)
			} else {

				// fmt.Fprintf(TCPConnection, "HTTP/1.1 404 Not Found\r\n\r\n")
				r := NewHTTPResponse(Request, Server, HTTPResponse{
					Version:       Request.Version,
					StatusCode:    404,
					StatusMessage: "Not Found",
				})
				HTTPRespond(r, TCPConnection)
			}

		} else {
			// fmt.Fprintf(TCPConnection, "HTTP/1.1 404 Not Found\r\n\r\n")
			r := NewHTTPResponse(Request, Server, HTTPResponse{
				Version:       Request.Version,
				StatusCode:    404,
				StatusMessage: "Not Found",
			})
			HTTPRespond(r, TCPConnection)
		}
		// close connection
		TCPConnection.Close()
	} else if RequestMethod == "POST" {
		if strings.Split(RequestPath, "/")[1] == "files" {
			FileName := strings.Split(RequestPath, "/")[2]
			err := os.WriteFile(*ServerDir+"/"+FileName, Request.Body, 0666)
			if err == nil {
				// fmt.Fprintf(TCPConnection, "HTTP/1.1 201 Created\r\n\r\n")
				r := NewHTTPResponse(Request, Server, HTTPResponse{
					Version:       Request.Version,
					StatusCode:    201,
					StatusMessage: "Created",
				})
				HTTPRespond(r, TCPConnection)
			} else {
				// fmt.Fprintf(TCPConnection, "HTTP/1.1 400 Bad Request\r\n\r\n")
				r := NewHTTPResponse(Request, Server, HTTPResponse{
					Version:       Request.Version,
					StatusCode:    400,
					StatusMessage: "Bad Request",
				})
				HTTPRespond(r, TCPConnection)

			}

		}
		// close connection
		TCPConnection.Close()
	}
}

func server() error {
	fmt.Println("Logs from your program will appear here!")
	Server := NewHTTPServer(ServerEncodings)
	ServerDir := flag.String("directory", "/tmp", "Specify your download directory")
	flag.Parse()
	ConnectionChannel := make(chan net.Conn)
	Listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		go AcceptConnection(Listener, ConnectionChannel)
		TCPConnection, OK := <-ConnectionChannel
		if !OK {
			break
		}
		go handleRequest(TCPConnection, Server, ServerDir)
	}
	return nil
}
