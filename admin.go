package admin

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Admin struct {
	token      string
	pwd        string
	webPath    string
	handlerMap map[string]ReqHandler
	addr       string
}

func New(pwd, webPath, addr string) *Admin {
	return &Admin{
		pwd:        pwd,
		webPath:    webPath,
		addr:       addr,
		handlerMap: make(map[string]ReqHandler),
	}
}

func (a *Admin) Run() {
	http.HandleFunc("/req", a.reqHandler)
	http.HandleFunc("/login", a.loginHandler)
	http.HandleFunc("/logout", a.authHandler(a.logoutHandler))
	http.HandleFunc("/admin/", a.authHandler(http.StripPrefix("/admin/", http.FileServer(http.Dir(a.webPath+"/admin/"))).ServeHTTP))
	http.HandleFunc("/", handleFile(a.webPath+`index.html`))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(a.webPath+"/static/"))))
	http.ListenAndServe(a.addr, nil)
}

func handleFile(filePath string) handleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filePath)
	}
}

type handleFunc func(http.ResponseWriter, *http.Request)

func (a *Admin) authHandler(handler handleFunc) handleFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/" || a.checkSession(request) {
			handler(response, request)
		} else {
			http.Redirect(response, request, "/", 302)
		}
	}
}

func (a *Admin) checkSession(request *http.Request) bool {
	if cookie, err := request.Cookie("session"); err == nil {
		return a.token != "" && cookie.Value == a.token
	}
	return false
}

func (a *Admin) logoutHandler(response http.ResponseWriter, request *http.Request) {
	a.token = ""
	http.Redirect(response, request, "/", 302)
}

func (a *Admin) loginHandler(response http.ResponseWriter, request *http.Request) {
	pwd := request.FormValue("pwd")
	redirectTarget := "/"
	if a.auth(pwd) {
		a.setSession(a.newToken(), response)
		redirectTarget = "/admin/home.html"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func (a *Admin) setSession(token string, response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(response, cookie)
}

func (a *Admin) newToken() string {
	lenMax := 6
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenMax; i++ {
		a.token += string(rune(r.Int31()%26) + rune('a'))
	}
	println(a.token)
	return a.token
}

func (a *Admin) auth(pwd string) bool {
	return pwd == a.pwd
}

func (a *Admin) reqHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	q := request.FormValue("q")
	f := a.handlerMap[q]
	if f != nil {
		fmt.Fprintf(response, f(request))
	} else {
		logErr("no handler %v", q)
		fmt.Fprintf(response, "")
	}
}

type ReqHandler func(*http.Request) string

func (a *Admin) RegHandler(q string, f ReqHandler) {
	if a.handlerMap[q] == nil {
		a.handlerMap[q] = f
	} else {
		logErr("same ReqHandler %v", q)
	}
}

func logErr(format string, a ...interface{}) {
	fmt.Printf("[ERR]"+format+"\n", a)
}

func logInfo(format string, a ...interface{}) {
	fmt.Printf("[INFO]"+format+"\n", a)
}
