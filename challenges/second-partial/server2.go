package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
//	"log"
//	"net/http"
//	"regexp"
//	"context"
//	"strings"
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

func login(c *gin.Context) {						// asign token
	user, _, _ := c.Request.BasicAuth()				// If needed Basic auth returns: (username, password string, ok bool)

	message := fmt.Sprintf("Hi %s, welcome to the DPIP system", user)

	c.JSON(200, gin.H {
		"message": message,
		"token": "ojIE89GzFw",					// add the token
	})
}

type tokenHeader struct {
	Token string `header:"Authorization"`
}

func logout(c *gin.Context) {
	h := tokenHeader{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(700, err)						// err 700 -> header error
	}

	c.JSON(200, gin.H {
		//"token": h.Token,  						// this is how you acces a token
		"message": "Bye @username, your token has been revoked",	// add the user variable
	})
}

type uploadBody struct {
	Body string `form:"data"`
}

func upload(c *gin.Context) {
	h := tokenHeader{}
	//body := c.Clone() //c.Request.Body
//	body := c.FormValue("Body")
	data := uploadStruct{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(700, err)						// err 700 -> header error
	}

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(800, err)						// err 800 -> body error
	}

	c.JSON(200, gin.H {
		"message": data, // "Hi @username, welcome to the DPIP system",	// add the user variable
		"filename": "image.png",				// add the filename variable
		"size":  "500kb",					// add the file size
	})
}

func status(c *gin.Context) {
	c.JSON(200, gin.H {
		"message": "Hi #username, the DPIP System is Up and Running",	// add the user variable
		"time": "2015-03-07 11:06:39",					// add the token
	})
}

func main() {
	fmt.Printf("Starting server at port 8080\n")

	server := gin.Default()

	server.GET("/login", login)
	server.GET("/logout", logout)
	server.GET("/upload", upload)
	server.GET("/status", status)

	server.Run()
}
