package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func NoSurf(next http.Handler) http.Handler {
	csrfHanlder := nosurf.New(next)

	csrfHanlder.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHanlder
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
