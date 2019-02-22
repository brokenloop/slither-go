package main

import (
	"encoding/json"
	"testing"
)

func Sum(x int, y int) int {
	return x + y
}

var testRequest = GameRequest{
	Game: Game{
		Id: "gameid",
	},
	Turn: 0,
	Board: Board{
		Height: 10,
		Width:  10,
		Food:   []Coord{Coord{X: 5, Y: 5}},
		Snakes: []Snake{},
	},
	You: Snake{
		Id:    "yousnakeid",
		Name:  "yousnakename",
		Heath: 100,
		Body:  []Coord{Coord{X: 0, Y: 0}},
	},
}

func TestSum(t *testing.T) {
	total := Sum(5, 5)
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}

func TestParseWorldFromRequest(t *testing.T) {
	world := ParseWorldFromRequest(testRequest)
	fromTile := world.From()

	expectedFrom := testRequest.You.Body[0]

	if fromTile.X != expectedFrom.X || fromTile.Y != expectedFrom.Y {
		t.Errorf("Fromtile incorrectly set. got [%d, %d], wanted [%d, %d]", fromTile.X, fromTile.Y, expectedFrom.X, expectedFrom.Y)
	}
	PrintWorld(t.Log, world)
	// out, _ := json.Marshal(world)
	// if err != nil {
	// 	panic(err)
	// }
	// stringWorld := string(out)
	// t.Logf(stringWorld)
	//t.Logf("world looks like following \n %v", world)
}

func TestJsonMarshal(t *testing.T) {
	res, err := json.Marshal(testRequest)
	if err != nil {
		panic(err)
	}
	t.Logf(string(res))
}
