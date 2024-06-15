package cookie

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,   // MaxAge in seconds (0 means no Max-Age set)
		Path:     path,     // Cookie path
		Domain:   domain,   // Cookie domain
		Secure:   secure,   // HTTPS only
		HttpOnly: httpOnly, // HTTP only (not accessible via JavaScript)
	}
	http.SetCookie(w, cookie)
}

// GetCookie returns the value of the named cookie.
func GetCookie(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// DeleteCookie deletes a cookie with the given name by setting its MaxAge to -1.
func DeleteCookie(w http.ResponseWriter, name, path, domain string) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   "",
		MaxAge:  -1, // MaxAge < 0 means delete cookie now
		Path:    path,
		Domain:  domain,
		Expires: time.Unix(1, 0), // Expire immediately
	}
	http.SetCookie(w, cookie)
}
