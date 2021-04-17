package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"context"
	"strings"
	"encoding/base64"
//	"github.com/gorilla/mux"
)

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
    fields := r.Context().Value(ctxKey{}).([]string)
    return fields[index]
}

var routes = []route {
	newRoute("GET", "/", home),
	newRoute("GET", "/login", login),
	newRoute("GET", "/logout", logout),
	newRoute("GET", "/upload", upload),
	newRoute("GET", "/status", status),
}

func home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		http.Error(w, "Not authorized", 401)
		return
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		http.Error(w, "Not authorized", 401)
		return
	}

	if pair[0] != "username" || pair[1] != "password" {
		http.Error(w, "Not authorized", 401)
		return
	}
	fmt.Fprintf(w, "Login %s", r)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Logout")
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Upload")
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status")
}

func Serve(w http.ResponseWriter, r *http.Request) {
    var allow []string
    for _, route := range routes {
        matches := route.regex.FindStringSubmatch(r.URL.Path)
        if len(matches) > 0 {
            if r.Method != route.method {
                allow = append(allow, route.method)
                continue
            }
            ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
            route.handler(w, r.WithContext(ctx))
            return
        }
    }
    if len(allow) > 0 {
        w.Header().Set("Allow", strings.Join(allow, ", "))
        http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
        return
    }
    http.NotFound(w, r)
}

func main() {
	fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/", Serve)


	if err := http.ListenAndServe(":8080", nil); err != nil {
		 log.Fatal(err)
	}
}
