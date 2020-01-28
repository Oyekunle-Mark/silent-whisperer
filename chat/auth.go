package main

import "net/http"

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")

	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	a.next.ServeHTTP(w, r)
}

// MustAuth protects a handler using authHandler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
