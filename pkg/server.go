package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func AcceptConnection(Listener net.Listener, ConnectionChannel chan net.Conn) {
	conn, err := Listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	ConnectionChannel <- conn
}

func server() error {
	fmt.Println("Logs from your program will appear here!")

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

			}
			if header == "\r\n" {
				break
			}

			key := strings.SplitAfterN(header, ":", 2)[0]
			key = strings.Trim(key, ":")
			value := strings.SplitAfterN(header, ":", 2)[1]
			value = strings.TrimSpace(value)
			Request.Header[key] = value
		}

		// http path handling
		RequestPath := Request.Path
		if RequestPath == "/" {
			fmt.Fprintf(TCPConnection, "HTTP/1.1 200 OK\r\n\r\n")

			// /echo/{str}
		} else if strings.Split(RequestPath, "/")[1] == "echo" {
			echo_string := strings.Split(RequestPath, "/")[2]
			fmt.Fprintf(TCPConnection, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v \r\n\r\n%v", len(echo_string), echo_string)
		} else if RequestPath == "/user-agent" {
			// `GET
			// /user-agent
			// HTTP/1.1
			// \r\n
			// Host: localhost:4221\r\n
			// User-Agent: foobar/1.2.3\r\n
			// Accept: */*\r\n
			// \r\n`
			UserAgent := Request.Header["User-Agent"]
			fmt.Fprintf(TCPConnection, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(UserAgent), UserAgent)
		} else {
			fmt.Fprintf(TCPConnection, "HTTP/1.1 404 Not Found\r\n\r\n")
		}
		TCPConnection.Close()
	}
	return nil
}
