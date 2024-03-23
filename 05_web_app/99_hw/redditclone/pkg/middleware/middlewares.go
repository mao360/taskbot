package middleware

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"reddit_clone/pkg/session"
	"strings"
	"time"
)

// нигде никакого редиректа быть не должно
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		token := strings.Split(authHeader, " ")
		userMap, err := session.CheckSess(token[1])
		if err != nil {
			panic(err) // fix it
		}

		ctx := context.WithValue(r.Context(), session.Key, userMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccessLog(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("access log middleware")
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
