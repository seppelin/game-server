package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type SessionNameKey struct{}

func SessionName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_name")
		if err != nil || sessionCookie.Value == "" {

			sessionCookie = GetCookie(uuid.New().String())
			http.SetCookie(w, sessionCookie)
			log.Println("new session:", sessionCookie.Value)
		} else {
			log.Println("cookie session:", sessionCookie.Value)
		}
		ctx := context.WithValue(r.Context(), SessionNameKey{}, sessionCookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSignOutCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_name",
		Value:    name,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
}

func GetCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_name",
		Value:    name,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
}
