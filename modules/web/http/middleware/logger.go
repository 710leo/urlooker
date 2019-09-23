package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "", 0)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Printf("%v %s %s %s", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.URL.Path, r.RemoteAddr)

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Printf("%v %v %s in %v", time.Now().Format("2006/01/02 15:04:05"), res.Status(), http.StatusText(res.Status()), time.Since(start))
}
