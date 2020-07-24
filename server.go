package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"image"
	"time"
	"encoding/json"
)

func main() {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/fluid", fluidHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	fmt.Fprintf(w, "POST request successful\n")
	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

func fluidHandler(w http.ResponseWriter, r *http.Request) {
	p := image.Pt(20, 20)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			p.X++
			p.Y++
			reqBody, err := json.Marshal(map[string]string {
				"point": p.String(),
			})
			if err != nil {
				print(err)
			}
			resp, err := http.Post("http://localhost:8080/fluid", "application/json", bytes.NewBuffer(reqBody))
			defer resp.Body.Close()
		}
	}
}

// MovingPoint is a Point that has a Point and accelleration assigned to it.
type MovingPoint struct {
	p Point
	s Point
	a Point
}

func (m *MovingPoint) next() {
	m.p.Add(&m.s)
	m.s.Add(&m.a)
}

type Field struct {
	image.Rectangle
	c image.Point
}

func Fld(r image.Rectangle) Field {
	c := image.Pt((r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2)
	return Field{r, c}
}

type Point struct {
	x, y float64
}

func (s1 *Point) Add(s2 *Point) {
	s1.x += s2.x
	s1.y += s2.y
}