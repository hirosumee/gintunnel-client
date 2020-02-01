package tcp_server

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var from, to string

func StartTCP(_from string, _to string) {
	from = _from
	to = _to
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
	svConn, err := net.Dial("tcp", to)
	if err != nil {
		sendError(conn)
		logrus.Error(err)
		return
	}
	name := editHeader(conn, svConn)
	c := make(chan string)
	go transferRequest(conn, svConn)
	go transferResponse(c, svConn, conn)
	meta := <- c
	logrus.Info(name + " " + meta)
}
func editHeader(client net.Conn, server net.Conn) string {
	buf := make([]byte, 2*1024)
	client.Read(buf)
	str := string(buf)
	requestName := getRequestName(str)
	str = editHostname(str)
	server.Write([]byte(str))
	return requestName
}

func editHostname(str string) string {
	return strings.ReplaceAll(str, from, to)
}
func transferRequest(src io.ReadCloser, dest io.WriteCloser) {
	defer src.Close()
	defer dest.Close()
	io.Copy(dest, src)
}
func transferResponse(c chan string ,src io.ReadCloser, dest io.WriteCloser) {
	defer src.Close()
	defer dest.Close()
	var b bytes.Buffer
	w := io.MultiWriter(dest, &b)
	io.Copy(w, src)
	meta, _  := b.ReadString('\n')
	c <- strings.TrimSpace((meta)[8:])
}

func sendError(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 503 Service Unavailable Error\r\n\r\n"))
	conn.Write(getOrRead503Page())
	conn.Write([]byte("\r\n"))
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

var page503 []byte

func getOrRead503Page() []byte {
	if len(page503) == 0 {
		var path, _ = filepath.Abs("./tcp-server/pages/503.html")
		file, err := os.Open(path)
		if err != nil {
			logrus.Error(file)
			return page503
		}
		page503 = make([]byte, 164*1024)
		_, err = file.Read(page503)
		if err != nil {
			logrus.Error(file)
		}
		return page503
	} else {
		return page503
	}
}
