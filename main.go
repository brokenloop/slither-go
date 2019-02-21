package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

type test_struct struct {
	Test string
}

type Coord struct {
	X, Y int
}

type Snake struct {
	Id, Name string
	Heath    int
	Body     []Coord
}

type Board struct {
	Height, Width int
	Food          []Coord
	Snakes        []Snake
}

type Game struct {
	Id string
}

type GameRequest struct {
	Game  Game
	Turn  int
	Board Board
	You   Snake
}

func FindMove(g GameRequest) {
	world := ParseWorldFromRequest(g)
	p, _, found := Path(world.From(), world.To())
	if !found {
		fmt.Println("Could not find a path")
	} else {
		fmt.Println("Resulting path\n", world.RenderPath(p))
	}

}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	returnString := "ssssSSsssSsSSSsssSSSSSssssSSss"
	json.NewEncoder(w).Encode(returnString)
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var input GameRequest
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(input, "", "  ")
	fmt.Println(b)
	spew.Dump(input)

	// FindMove(input)
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
}

func EndHandler(w http.ResponseWriter, r *http.Request) {
}

// main function to boot up everything
func main() {
	TestStraightLine()
	router := mux.NewRouter()
	router.HandleFunc("/", PingHandler).Methods("GET")
	router.HandleFunc("/start", StartHandler).Methods("POST")
	router.HandleFunc("/move", MoveHandler).Methods("POST")
	router.HandleFunc("/end", EndHandler).Methods("POST")
	router.HandleFunc("/ping", PingHandler).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
