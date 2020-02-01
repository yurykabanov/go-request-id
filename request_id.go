package requestid

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
)

const (
	Header = "X-Request-Id"
	ContextKey = "__request_id__"
)

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, ContextKey, requestId)
}

func GetFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	value, ok := ctx.Value(ContextKey).(string)
	return value, ok
}

func Middleware(next http.Handler, nextRequestId func() string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get(Header)

		if requestId == "" {
			requestId = nextRequestId()
		}

		ctx := WithRequestId(r.Context(), requestId)

		w.Header().Set(Header, requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DefaultRequestIdProvider() string {
	var buf = make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%02x", buf)
}
