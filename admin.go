package main

import (
	. "./config"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var token = ""

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", authHandler(logoutHandler))
	http.HandleFunc("/home", authHandler(homeHandler))
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":9999", nil)
}

type handleFunc func(http.ResponseWriter, *http.Request)

func authHandler(handler handleFunc) handleFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if checkSession(request) {
			handler(response, request)
		} else {
			http.Redirect(response, request, "/", 302)
		}
	}
}

func indexHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, IndexPage)
}

func homeHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, HomePage)
}

func checkSession(request *http.Request) bool {
	if cookie, err := request.Cookie("session"); err == nil {
		return token != "" && cookie.Value == token
	}
	return false
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	token = ""
	http.Redirect(response, request, "/", 302)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	pwd := request.FormValue("pwd")
	redirectTarget := "/"
	if auth(pwd) {
		setSession(newToken(), response)
		redirectTarget = "/home"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func setSession(token string, response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(response, cookie)
}

func newToken() string {
	lenMax := 6
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenMax; i++ {
		token += string(rune(r.Int31()%26) + rune('a'))
	}
	println(token)
	return token
}

func auth(pwd string) bool {
	return pwd == Pwd
}
