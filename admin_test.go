package admin

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	admin := New("123", `web\`, ":9999")
	admin.RegHandler("test", func(req *http.Request) string {
		return "-->" + req.FormValue("p")
	})
	admin.Run()
}
