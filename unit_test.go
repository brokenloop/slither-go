package main

import (
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
		Height: 5,
		Width:  5,
		Food:   []Coord{Coord{X: 4, Y: 4}},
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
	// out, err := json.Marshal(world)
	// if err != nil {
	// 	panic(err)
	// }
	// stringWorld := string(out)
	// t.Logf("")
	t.Logf("world looks like following \n %v", world)
}
