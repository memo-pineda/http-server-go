package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error reading request ", err.Error())
		return
	}

	fmt.Printf("Request: %s %s\n", request.Method, request.URL.Path)

	if request.URL.Path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	urlParts := strings.Split(request.URL.Path, "/")
	endpoint := urlParts[1]

	if len(urlParts) > 1 && endpoint == "echo" {
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(urlParts[2]), urlParts[2])
		conn.Write([]byte(resp))
		return
	}

	if endpoint == "user-agent" {
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(request.UserAgent()), request.UserAgent())
		conn.Write([]byte(resp))
		return
	}

	if endpoint == "files" {
		fmt.Println("Reaching files endpoint")
		filePath := "/" + urlParts[2]
		dir := os.Args[2]
		fmt.Println(dir + filePath)
		file, err := os.ReadFile(dir + filePath)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(file), file)
			conn.Write([]byte(resp))
		}
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
