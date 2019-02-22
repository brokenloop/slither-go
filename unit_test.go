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
		Food:   []Coord{Coord{X: 9, Y: 9}},
		Snakes: []Snake{
			Snake{
				Id:    "themsnakeid",
				Name:  "themsnakename",
				Heath: 100,
				Body:  []Coord{Coord{X: 3, Y: 3}, Coord{X: 3, Y: 2}, Coord{X: 3, Y: 1}},
			},
		},
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
}

func TestPrintFindMove(t *testing.T) {
	world := ParseWorldFromRequest(testRequest)
	p, _, found := Path(world.From(), world.To())
	if !found {
		t.Log("Could not find a path")
	} else {
		t.Log("Resulting path\n", world.RenderPath(p))
	}
}

func TestFindMove(t *testing.T) {
	expectedMove := "down"
	move := FindMove(testRequest)
	if move != expectedMove {
		t.Errorf("Incorrect move. Got %v, expected %v", move, expectedMove)
	}
}

// func TestParseMove(t *testing.T) {
// 	world := ParseWorldFromRequest(testRequest)
// 	p, _, _ := Path(world.From(), world.To())
// 	from := Coord{
// 		world.From().X,
// 		world.From().Y,
// 	}

// }

func TestJsonMarshal(t *testing.T) {
	res, err := json.Marshal(testRequest)
	if err != nil {
		panic(err)
	}
	t.Logf(string(res))
}
