package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"context"
	"strings"
)

/*
var routes = [] route {
	newRoute("GET", "/", home),
	newRoute("GET", "/login", login),
	newRoute("GET", "/logout", logout),
	newRoute("GET", "/upload", upload),
	newRoute("GET", "/status", status),
}
*/

/*
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login")
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
}*/

func login(c *gin.Context) {
	c.JSON(200, gin.H {
		"message": "Hi @username, welcome to the DPIP system",	// add the user variable
		"token": "ojIE89GzFw",					// add the token
	})
}

func main() {
	fmt.Printf("Starting server at port 8080\n")

	r := gin.Default()

	r.GET("/login", login)
	r.GET("/logout", logout)
	r.GET("/upload", upload)
	r.GET("/status", status)

	r.Run()
	/*fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/", Serve)


	if err := http.ListenAndServe(":8080", nil); err != nil {
		 log.Fatal(err)
	}*/
}
