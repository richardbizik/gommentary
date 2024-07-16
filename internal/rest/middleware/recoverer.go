package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Error struct {
	Type    ResultType "json:\"severity\""
	Message string     "json:\"message\""
	Code    string     "json:\"code\""
}

type ResultType string

const (
	ERROR ResultType = "ERROR"
	WARN  ResultType = "WARN"
	INFO  ResultType = "INFO"
)

func getErrorResult(severity ResultType, message string, code string) Error {
	return Error{
		Type:    severity,
		Message: message,
		Code:    code,
	}
}

func handleError(w http.ResponseWriter, r *http.Request, recovered interface{}) {
	switch recovered := recovered.(type) {
	//defaults
	case error:
		writeError(w, r, http.StatusInternalServerError, getErrorResult(ERROR, recovered.Error(), "INTERNAL_SERVER_ERROR"))
	default:
		writeError(w, r, http.StatusInternalServerError, getErrorResult(ERROR, fmt.Sprintf("%v", recovered), "INTERNAL_SERVER_ERROR"))
	}
}

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}
				handleError(w, r, rvr)

				slog.Error(fmt.Sprint(rvr))
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, resp interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	body, err := json.Marshal(resp)
	if err != nil {
		slog.Error(err.Error())
	} else {
		_, err = w.Write(body)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}
