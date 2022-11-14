package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func Log(logger *logrus.Entry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("New request: method - %s, url - %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
