// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
	"fmt"
	"os"
	"flag"
)

func handleConn(c net.Conn, timeTZ string) {
	defer c.Close()
	_, err := io.WriteString(c, timeTZ)
	if err != nil {
		return // e.g., client disconnected
	}
}

// Reads the enviroment variables
func parametersParser() (string, int) {
	tz := os.Getenv("TZ")				// TZ
	var port = flag.Int("port", 1234, "Port")	// Port
	flag.Parse()
	return tz, *port;
}

// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

// Returns a formated string with the timezone
func getTime(tz string) string {
	t, err := TimeIn(time.Now(), tz)
	if err == nil {
		return fmt.Sprintf("%v\t: %v\n", t.Location(), t.Format("15:04:05"))
	} else {
		return fmt.Sprintf("%v\t: <timezone unknown>\n", t.Location())
	}
}

func main() {
	// get the enviroment variables
	tz, port := parametersParser()
	// get the time for the timezone
	timeTZ := getTime(tz)

	// listener string
	server := fmt.Sprintf("localhost:%v", port)

	// Connection
	listener, err := net.Listen("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn, timeTZ) // handle connections concurrently
	}
}
