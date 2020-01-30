package http_server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartHTTP() {
	server := http.Server{Addr: ":8084", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("a new request")
		fmt.Fprintln(w, "haha "+r.URL.Path)
	})}
	logrus.Info("start http sv")
	logrus.Error(server.ListenAndServe())
}
