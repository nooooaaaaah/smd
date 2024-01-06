package utils

import "net/http"

type AuthUtil struct{}

func (a *AuthUtil) getUserIDFromRequest(r *http.Request) (string, error) {
	// Implementation depends on )your authentication mechanism
	return "1", nil
}
