package middleware

import (
	"net/http"
	"time"

	"github.com/gogf/gf/os/gsession"
	"github.com/spf13/viper"
)

var (
	sessionCookie  = "webvpn_session"
	sessionManager = gsession.New(86400*time.Second, gsession.NewStorageMemory())
)

type Authentication struct{}

func (m *Authentication) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("webvpn_session")
		if err != nil || cookie == nil {
			http.SetCookie(w, &http.Cookie{
				Name:     sessionCookie,
				Value:    gsession.NewSessionId(),
				Expires:  time.Time{},
				MaxAge:   60 * 60 * 24 * 7,
				HttpOnly: true,
			})
			m.redirectToLogin(w, r)
			return
		}

		session := sessionManager.New(cookie.Name)
		userId := session.GetString("userId", "")

		if userId == "" { // not logged in now
			m.redirectToLogin(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (*Authentication) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/.webvpn/callback" {
		return
	}

	authenticationServiceURL := viper.GetString("authentication.service_scheme") + "://" + viper.GetString("authentication.service_url")
	http.Redirect(w, r, authenticationServiceURL+"/authentication/login", http.StatusFound)
}
