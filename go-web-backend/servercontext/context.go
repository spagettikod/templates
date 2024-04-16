package servercontext

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
)

type Context struct {
	Logger     *slog.Logger
	baseLogger *slog.Logger
	SessionID  string
}

type ContextKey string

const contextKey ContextKey = "serverContextKey"

func New() Context {
	return Context{baseLogger: slog.Default()}
}

func FromRequest(r *http.Request) Context {
	return r.Context().Value(contextKey).(Context)
}

func sessionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", b[0:4]), nil
}

func (c Context) Wrap(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// generate a session identifier to enable tracing of logs
		c.SessionID, err = sessionID()
		if err != nil {
			c.Logger.Error("error generating session id", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// add attributes to logger
		c.Logger = c.baseLogger.With("session", c.SessionID).WithGroup("request").With("method", r.Method).With("url", r.URL)

		// merge context the exist request context
		newCtx := context.WithValue(r.Context(), contextKey, c)

		// call wrapped handler
		fn(w, r.WithContext(newCtx))
	}
}
