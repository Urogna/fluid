package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"time"
	"github.com/Urogna/point"
	"os"
    "os/exec"
)

func main() {
	/*
		fileServer := http.FileServer(http.Dir("./static"))
		http.Handle("/", fileServer)
		http.HandleFunc("/hello", helloHandler)
		http.HandleFunc("/form", formHandler)
		http.HandleFunc("/fluid", fluidHandler)

		fmt.Printf("Starting server at port 8080\n")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	*/
	//var print [point.H][point.W]int
	grid := point.FieldGrid(point.FX, point.FY, point.W, point.H)
	balls := point.RandomMovingPoints(point.N, point.W, point.H)
	ticker := time.NewTicker(time.Second/30)
	point.SetAverage(grid, point.N)
	for {
		select {
		case t := <-ticker.C:
			cmd := exec.Command("clear") //Linux example, its tested
			cmd.Stdout = os.Stdout
			cmd.Run()
			point.SetBallsNumber(balls, grid)
			for i := range balls {
				b := &balls[i]
				//b.Print(&print)
				b.SetAccelleration(grid)
				b.Next()
			}
			printGrid2(grid)
			//printGrid(&print)
			point.ResetField(grid)
			
			fmt.Println(t.Second())
		}
	}
}

func printGrid(g *[point.H][point.W]int) {
	for y := range g {
		for x := range g[y] {
			if g[y][x] != 0 {
				fmt.Print("0")
				g[y][x] = 0
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func printGrid2(g [][]point.Field) {
	for y := range g {
		for x := range g[y] {
			n := g[y][x].N
			g[y][x].N = 0
			switch {
			case n >= 12:
				fmt.Print("|:|")
			case n >= 12 && n < 10:
				fmt.Print(":|:")
			case n < 10 && n >= 7:
				fmt.Print(":::")
			case n < 7 && n >= 5:
				fmt.Print(":.:")
			case n < 5 && n >= 3:
				fmt.Print("...")
			case n < 3 && n >= 0:
				fmt.Print(" . ")
			case n < 0:
				fmt.Print("   ")
			default:
				fmt.Print("|||")
			}
		}
		fmt.Println()
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
			reqBody, err := json.Marshal(map[string]string{
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
