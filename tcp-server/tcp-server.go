package tcp_server

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
)

var hostname, port string

func StartTCP(_hostname string, _port string) {
	hostname = _hostname
	port = _port
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		logrus.Fatal(err)
		return
	}
	logrus.Info("start tcp sv")
	listen(listener)

}
func listen(listener net.Listener) {
	for {
		conn, err := (listener).Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	svConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		sendNotFound(conn)
		checkErr(err)
		return
	}
	editHeader(conn, svConn)
	go transfer(conn, svConn, 1)
	go transfer(svConn, conn, 2)
}
func editHeader(client net.Conn, server net.Conn) {
	buf := make([]byte, 2*1024)
	client.Read(buf)
	str := string(buf)
	requestName := getRequestName(str)
	logrus.Info(requestName)
	str = editHostname(str)
	server.Write([]byte(str))
}

func editHostname(str string) string {
	return strings.ReplaceAll(str, hostname, fmt.Sprintf("localhost:%s", port))
}

func transfer(src io.ReadCloser, dest io.WriteCloser, mode int8) {
	defer src.Close()
	defer dest.Close()
	if mode == 1 {
		io.Copy(dest, src)
	} else if mode == 2 {
		buf := make([]byte, 100)
		buffer := bytes.NewBuffer(buf)
		w := io.MultiWriter(dest, os.Stdout, buffer)
		io.Copy(w, src)
		logrus.Info(buffer.ReadString('\n'))
	}
}

func sendNotFound(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
	conn.Write([]byte("\r\n<h1>HI</h1>\r\n"))
	conn.Close()
}

func checkErr(err error) {
	if err == nil {
		println("Ok")
		return

	} else if netError, ok := err.(net.Error); ok && netError.Timeout() {
		println("Timeout")
		return
	}

	switch t := err.(type) {
	case *net.OpError:
		if t.Op == "dial" {
			println("Unknown host")
		} else if t.Op == "read" {
			println("Connection refused")
		}
	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			println("Connection refused")
		}
	}
}
func getRequestName(str string) string {
	var index = -1
	for i, v := range str[:] {
		if v == '\r' {
			index = i
			break
		}
	}
	if index == -1 {
		return ""
	}
	return str[:index]
}
