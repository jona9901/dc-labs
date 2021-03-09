/*
This is the client program
*/

package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// Argument parser
func parser(args []string, ch chan string) {
	for _, arg := range  args {
		server := strings.Split(arg, "=")[1]
		ch <- server // asign the string server to the channel
	}
	close(ch)
}

// get the connection with the server
func dial(ch chan string) {
	for i := range ch {
		// asing the channel to a string
		server := i
		
		// Get the connection to the server
		conn, err := net.Dial("tcp", server)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		client(os.Stdout, conn)
	}
}

// client side
func client(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Get the enviroment variables
	args := os.Args[1:]

	// Create the unbuffured channel
	ch := make(chan string)
	go parser(args, ch)
	dial(ch)
}
