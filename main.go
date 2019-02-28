// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/heroku/x/hmetrics/onload"
// )

// func main() {
// 	port := os.Getenv("PORT")

// 	if port == "" {
// 		log.Fatal("$PORT must be set")
// 	}

// 	router := gin.New()
// 	router.Use(gin.Logger())
// 	router.LoadHTMLGlob("templates/*.tmpl.html")
// 	router.Static("/static", "static")

// 	router.GET("/", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "index.tmpl.html", nil)
// 	})

// 	router.Run(":" + port)
// }

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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
	Health   int
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

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PING")
	returnString := "ssssSSsssSsSSSsssSSSSSssssSSss"
	json.NewEncoder(w).Encode(returnString)
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("START")
	decoder := json.NewDecoder(r.Body)
	var input GameRequest
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}

	// b, err := json.MarshalIndent(input, "", "  ")
	// fmt.Println(b)
	// spew.Dump(input)

	response := make(map[string]string)

	response["responscolor"] = "#ff00ff"
	response["headType"] = "bendr"
	response["tailType"] = "pixel"

	json.NewEncoder(w).Encode(response)
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("MOVE")
	decoder := json.NewDecoder(r.Body)
	var input GameRequest
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}
	// spew.Dump(input)
	move := FindMove(input)
	// move := FindMoveBySimulation(input)
	spew.Dump(move)
	response := make(map[string]string)
	response["move"] = move
	json.NewEncoder(w).Encode(response)
}

func EndHandler(w http.ResponseWriter, r *http.Request) {
}

// main function to boot up everything
func main() {
	// set the random seed for later use
	rand.Seed(time.Now().UTC().UnixNano())
	// TestStraightLine()

	router := mux.NewRouter()
	router.HandleFunc("/", PingHandler).Methods("GET")
	router.HandleFunc("/start", StartHandler).Methods("POST")
	router.HandleFunc("/move", MoveHandler).Methods("POST")
	router.HandleFunc("/end", EndHandler).Methods("POST")
	router.HandleFunc("/ping", PingHandler).Methods("GET", "POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// port := "8000"
	fmt.Println("Dispensing snakes on port: " + port)
	// log.Fatal(http.ListenAndServe(":"+port, router))
	log.Fatal(http.ListenAndServe(":"+port, router))
}
