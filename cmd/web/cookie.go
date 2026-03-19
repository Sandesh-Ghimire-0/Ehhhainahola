package main

import (
	"net/http"
)
func GetCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("username")
	if err != nil {
    	return nil
	} else {
    	return cookie
	}
}
