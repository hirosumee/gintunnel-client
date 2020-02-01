package main

import (
	"flag"
	tcp_client "gintunnel-client/tcp-client"
	tcp_server "gintunnel-client/tcp-server"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func main() {
	from, to := getConfig()
	if from == "" {
		logrus.Fatal("Hostname is required")
	}
	logrus.Infof("start with server hostname : %s and redirect to : %s", from, to)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		//http_server.StartHTTP()
		wg.Done()
	}()
	go func() {
		tcp_server.StartTCP(from, to)
		wg.Done()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		tcp_client.StartTcpClient(from)
		wg.Done()
	}()
	wg.Wait()
}
func getConfig() (hostname string, port string) {
	flag.StringVar(&hostname, "from", "localhost:8080", "hostname of tunnel")
	flag.StringVar(&port, "to", "localhost:8084", "address of redirecting hostname")
	flag.Parse()
	return
}
