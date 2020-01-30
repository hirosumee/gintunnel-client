package main

import (
	"flag"
	http_server "gintunnel-client/http-server"
	tcp_client "gintunnel-client/tcp-client"
	tcp_server "gintunnel-client/tcp-server"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func main() {
	hostname, port := getConfig()
	if hostname == "" {
		logrus.Fatal("Hostname is required")
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		http_server.StartHTTP()
		wg.Done()
	}()
	go func() {
		tcp_server.StartTCP(hostname, port)
		wg.Done()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		tcp_client.StartTcpClient(port)
		wg.Done()
	}()
	wg.Wait()
}
func getConfig() (hostname string, port string) {
	flag.StringVar(&hostname, "hostname", "", "hostname of tunnel")
	flag.StringVar(&port, "port", "80", "local port")
	flag.Parse()
	return
}
