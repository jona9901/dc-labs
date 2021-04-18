package main

import (
	"fmt"
	"os"
	"io"
	"time"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
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

type User struct {
	Username string
	Password string
	Token string
}

var loggedUsers []*User
var files []*os.File
/*
A sample user
*/

var defaultUser = User {
	Username: "username",
	Password: "password",
}

func newUser(username string, password string) *User{
	token, _ := createToken(username)				// handle error

	u := User{
		Username: username,
		Password: password,
		Token: token,
	}
	return &u
}

func createToken(username string) (string, error) {
	var err error

	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")			// ad to env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func login(c *gin.Context) {						// asign token
	username, password, _ := c.Request.BasicAuth()				// If needed Basic auth returns: (username, password string, ok bool)

	if username != defaultUser.Username || password != defaultUser.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid credentials")
		return
	}

	user := newUser(username, password)
	loggedUsers = append(loggedUsers, user)

	message := fmt.Sprintf("Hi %s welcome to the DPIP system", user.Username)

	c.JSON(200, gin.H {
		"message": message,
		"token": user.Token,
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

	for i, u := range loggedUsers {
		if h.Token == u.Token {
			message := fmt.Sprintf("Bye %s, your token has been revoked", u.Username)

			loggedUsers = append(loggedUsers[:i], loggedUsers[i + 1:]...)

			c.JSON(200, gin.H {
				"message": message,
			})
			return
		}

		c.JSON(500, gin.H {					// err 500 -> bad token
			"message": "Error, not logged in",
		})
	}
}

type uploadBody struct {
	Body string `form:"data"`
}

func upload(c *gin.Context) {
	h := tokenHeader{}
	data := uploadBody{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(700, err)						// err 700 -> header error
	}

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(800, err)						// err 800 -> body error
	}

	for _, u := range loggedUsers {
		if h.Token == u.Token {
			c.Request.ParseMultipartForm(32 << 20)

			f, err := os.Open(data.Body)
			if err != nil {
				c.JSON(100, gin.H {				// 100 -> file upload error
					"message": err,
				})
				return
			}

			defer f.Close()

			body := strings.Split(data.Body, "/")
			filename := body[len(body) - 1]

			newFile, err := os.Create(filename)

			if err != nil {
				c.JSON(101, gin.H {				// 101 -> file create error
					"message": err,
				})
				return
			}

			if _, err := io.Copy(newFile, f); err != nil {
				c.JSON(102, gin.H {				// 103 -> file copy error
					"message": err,
				})
				return

			}

			files = append(files, newFile)

			fmt.Println(files)
			fi, err := f.Stat()
			if err != nil {
				c.JSON(103, gin.H {				// 103 -> file size error
					"message": err,
				})
				return
			}

			size := fmt.Sprintf("%d.kb", fi.Size())

			c.JSON(200, gin.H {
				"message": "An image has been successfully uploaded",
				"filename": filename,				// add the filename variable
				"size": size,
			})
			return
		}
	}

	c.JSON(500, gin.H {					// err 500 -> bad token
		"message": "Error, not logged in",
	})
}

func status(c *gin.Context) {
	h := tokenHeader{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(700, err)						// err 700 -> header error
	}

	for _, u := range loggedUsers {
		if h.Token == u.Token {
			message := fmt.Sprintf("Hi %s, the DPIP System is Up and Running", u.Username)

			c.JSON(200, gin.H {
				"message": message,
				"time": time.Now().Format("2006-Jan-02 15:04:05"),//("2015-03-07 11:06:39"),
			})
			return
		}
	}

	c.JSON(500, gin.H {					// err 500 -> bad token
		"message": "Error, not logged in",
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
