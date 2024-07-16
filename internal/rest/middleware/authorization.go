package middleware

import (
	context "context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/richardbizik/gommentary/internal/appctx"
	"github.com/richardbizik/gommentary/internal/rest/handlers"
)

const (
	AUTH_HEADER = "authorization"
)

type JWTConfig struct {
	JWTRequired bool   `yaml:"required" json:"required" default:"true"`
	JWTHeader   string `yaml:"header" json:"header" default:"Authorization"`
}

type AuthToken struct {
	Jti string "json:\"jti\""
	Exp int64  "json:\"exp\""
	Iss string "json:\"iss\""
	Sub string "json:\"sub\""
}

func Authorization(config JWTConfig) func(next http.Handler) http.Handler {
	authorizationHeader := config.JWTHeader
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !config.JWTRequired {
				next.ServeHTTP(w, r)
				return
			}
			authorizationValue := r.Header.Get(authorizationHeader)

			if authorizationValue != "" {
				authorizationValue = authorizationValue[7:] //remove "bearer " string
			}

			if authorizationValue != "" {
				split := strings.Split(authorizationValue, ".")
				if len(split) > 2 {
					body := split[1]
					token := AuthToken{}
					var encoder = base64.URLEncoding.WithPadding(base64.NoPadding)
					decoded, err := encoder.DecodeString(body)
					if err != nil {
						setBadRequestAuthError(w, "Authorization header has wrong format")
						return
					}
					err = json.Unmarshal(decoded, &token)
					if err != nil {
						slog.Debug(fmt.Sprintf("%v", err))
						setBadRequestAuthError(w, "Authorization header has wrong format")
						return
					}

					if token.Sub != "" {
						ctx := r.Context()
						ctx = context.WithValue(ctx, appctx.USERNAME, &token.Sub)
						r = r.WithContext(ctx)
					}
				} else {
					setBadRequestAuthError(w, "Authorization header has wrong format")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func setBadRequestAuthError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write(GetJson(handlers.Error{
		Code:     "AUTH_ERROR",
		Message:  message,
		Severity: "ERROR",
	}))
	if err != nil {
		panic(err)
	}
}

func GetJson(i interface{}) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return b
}
