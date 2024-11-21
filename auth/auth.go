package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "57824297361-uijfar3d9rkkrcn57ndro61gob1l8j4h.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	GithubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		ClientID:     "Ov23litcIJ8mWbAQ3SRR",
		ClientSecret: "66130a2e7e33f3b6b0a6b00b9e6d0877f9cdba7c",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	FacebookOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/facebook/callback",
		ClientID:     "1500867300637924",
		ClientSecret: "56e5674c68058163ae39345d4de20c07",
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}
)

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = 365 * 24 * 60 * 60
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: time.Now().Add(time.Second * time.Duration(expiration))}
	http.SetCookie(w, &cookie)
	return state
}
