package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"math"
)

type Point struct {
	X, Y float64
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//generatePoints array
func generatePoints(s string) ([]Point, error) {

	points := []Point{}

	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	vals := strings.Split(s, ",")
	if len(vals) < 2 {
		return []Point{}, fmt.Errorf("Point [%v] was not well defined", s)
	}

	var x, y float64

	for idx, val := range vals {

		if idx%2 == 0 {
			x, _ = strconv.ParseFloat(val, 64)
		} else {
			y, _ = strconv.ParseFloat(val, 64)
			points = append(points, Point{x, y})
		}
	}
	return points, nil
}

// getArea gets the area inside from a given shape
func getArea(points []Point) float64 {
	buff := 0.0
	buff2 := 0.0

	size := len(points)

	for i:= 0; i < size; i++ {
		if i == size - 1{
			buff += points[i].X * points[0].Y
			buff2 += points[i].Y * points[0].X
		} else {
			buff += points[i].X * points[i+1].Y
			buff2 += points[i].Y * points[i+1].X
		}
	}

	return math.Abs(buff - buff2) / 2
}

// Calculate the distance betwen two points
func distance(a, b Point) float64{
	return math.Sqrt(math.Pow(b.X - a.X, 2) + math.Pow(b.Y - a.Y, 2))
}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {
	perimeter := 0.0
	size := len(points)

	for i := 0; i < size; i++ {
		if i == size - 1 {
			perimeter += distance(points[i], points[0])
		} else {
			perimeter += distance(points[i], points[i + 1])
		}
	}
	
	return perimeter
}

func collisions(vertices []Point) bool {
	return true
}

// Given three colinear points p, q, r, the fucntion checks if th e point q lies on line segment 'pr'
func onSegment(p Point, q Point, r Point) bool {
	if ((q.X <= math.Max(p.X, r.X)) && (q.X >= math.Min(p.X, r.X)) && (q.Y <= math.Max(p.Y, r.Y)) && (q.Y >= math.Min(p.Y, r.Y))) {
		return false	
	}
	return true
}

// to find the orientation of an ordered triplet (p, q, r)
// function returns the following values
// 0: Colinear points
// 1: Clockwise Points
// 2: Counterclockwise
func orientation(p Point, q Point, r Point) int {
	val := 0.0

	val = ((q.Y - p.Y) * (r.X - q.X)) - ((q.X - p.X) * (r.Y - q.Y))

	if (val > 0) {
		return 1 // Clockwise orientation
	} else if (val < 0) {
		return 2 // Counterclockwise orientation
	} else {
		return 0 // Colinear orientation
	}
}

func doIntersect(p1 Point, q1 Point, p2 Point, q2 Point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if ((o1 != o2) && (o3 != o4)) {
		return true
	}
	if ((o1 == 0) && onSegment(p1, p2, q1)) {
		return true
	}
	if ((o2 == 0) && onSegment(p1, q2, q1)) {
		return true
	}
	if ((o3 == 0) && onSegment(p2, p1, q2)) {
		return true
	}
	if ((o4 == 0) && onSegment(p2, q1, q2)) {
		return true
	}
	return false
}

// handler handles the web request and reponds it
func handler(w http.ResponseWriter, r *http.Request) {

	var vertices []Point
	for k, v := range r.URL.Query() {
		if k == "vertices" {
			points, err := generatePoints(v[0])
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("error: %v", err))
				return
			}
			vertices = points
			break
		}
	}

	// Results gathering
	area := getArea(vertices)
	perimeter := getPerimeter(vertices)

	//onSegment(vertices[0], vertices[1], vertices[2])

	// Logging in the server side
	log.Printf("Received vertices array: %v", vertices)

	// Response construction
	response := fmt.Sprintf("Welcome to the Remote Shapes Analyzer\n")
	response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", len(vertices))
	response += fmt.Sprintf(" - Vertices        : %v\n", vertices)

	// Check if collision
	size := len(vertices)
	error := false
	
	if size > 2 {
		if size > 3 {
			for i := 0; i < size - 1; i++ {
				if i == size - 2 {
					if doIntersect(vertices[i], vertices[i + 1], vertices[0], vertices[1]){
						response += fmt.Sprintf("ERROR - There is an intersection betwen the segments\n")
						error = true
					}
				} else {
					for j := i + 2; j < size - 1; j++ {
						if doIntersect(vertices[i], vertices[i + 1], vertices[j], vertices[j + 1]) {
							response += fmt.Sprintf("ERROR - There is an intersection betwen the segments\n")
							error = true
						}
					}
				}
			}
		}

		if error == false {
			response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
			response += fmt.Sprintf(" - Area            : %v\n", area)
		}

		error = false

	} else {
		response += fmt.Sprintf("ERROR - Your shape is not compliying with the minimum number of vertices.\n")
		error = true
	}

	// Send response to client
	fmt.Fprintf(w, response)
}
