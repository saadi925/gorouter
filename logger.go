package gorouter

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

var golog = logrus.New()

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		golog.WithFields(logrus.Fields{
			"method": req.Method,
			"url":    req.URL.Path,
			"remote": req.RemoteAddr,
		}).Info("Received request")
		next.ServeHTTP(w, req)
	})
}
