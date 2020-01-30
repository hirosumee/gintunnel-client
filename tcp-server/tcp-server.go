package tcp_server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"net"
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
	editHostname(conn, svConn)
	go transfer(conn, svConn)
	go transfer(svConn, conn)
}
func editHostname(client net.Conn, server net.Conn) {
	buf := make([]byte, 100)
	client.Read(buf)
	str := string(buf)
	str = strings.ReplaceAll(str, hostname, fmt.Sprintf("localhost:%s", port))
	logrus.Info(str)
	server.Write([]byte(str))
}

func transfer(src io.ReadCloser, dest io.WriteCloser) {
	defer src.Close()
	defer dest.Close()
	//r := io.TeeReader(src, dest)
	//go printAll(r)
	io.Copy(dest, src)
}
func printAll(r io.Reader) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
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
