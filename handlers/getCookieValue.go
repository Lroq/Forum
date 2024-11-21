package handlers

import (
	"net/http"
)

func GetCookieValue(r *http.Request, cookieName string) (string) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return ""
		}
		return ""
	}
	return cookie.Value
}
