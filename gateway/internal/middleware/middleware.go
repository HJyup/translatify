package middleware

import (
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"net/http"
)

func WithMiddleware(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(handler))
}
