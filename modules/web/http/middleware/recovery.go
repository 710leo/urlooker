package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/710leo/urlooker/modules/web/http/errors"
	"github.com/710leo/urlooker/modules/web/http/render"
)

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[panic] ", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			if customError, ok := err.(errors.Error); ok {
				msg := fmt.Sprintf("[%s:%d] %s [Error:]%s", customError.File, customError.Line, customError.Time, customError.Msg)
				log.Printf(msg)

				if isAjax(r) {
					render.Message(w, customError.Msg)
					return
				}

				if customError.Code == http.StatusUnauthorized || customError.Code == http.StatusForbidden {
					http.Redirect(w, r, "/auth/login", 302)
					return
				}

				render.Put(r, "Error", customError.Msg)
				render.HTML(r, w, "inc/error")
				return
			}

			// Negroni part
			w.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			log.Printf(f, err, stack)

			if rec.PrintStack {
				fmt.Fprintf(w, f, err, stack)
			}
		}
	}()

	next(w, r)
}

func isAjax(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}
