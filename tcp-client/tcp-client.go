package tcp_client

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
)

func StartTcpClient(hostname string) {
	svConn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		logrus.Error(err)
		return
	}
	r := bufio.NewReader(svConn)
	w := bufio.NewWriter(svConn)
	_, _ = w.WriteString(fmt.Sprintf("REG %s\n", hostname))
	_ = w.Flush()
	for {
		temp, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			logrus.Error(err)
			break
		}
		message := strings.TrimSpace(temp)
		splited := strings.Split(message, " ")
		cmd := strings.TrimSpace(splited[0])
		var data string
		if len(splited) > 1 {
			data = strings.TrimSpace(splited[1])
		}
		switch cmd {
		case "REG-RES":
			{
				logrus.Info(data)
			}
		case "PING":
			{
				w.WriteString("PONG\r\n")
				w.Flush()
			}
		}

	}
	svConn.Close()
}
