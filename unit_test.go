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
		Height: 10,
		Width:  10,
		Food:   []Coord{Coord{X: 9, Y: 9}, Coord{X: 0, Y: 5}, Coord{X: 7, Y: 8}},
		Snakes: []Snake{
			Snake{
				Id:     "themsnakeid",
				Name:   "themsnakename",
				Health: 100,
				Body:   []Coord{Coord{X: 3, Y: 3}, Coord{X: 3, Y: 2}, Coord{X: 3, Y: 1}},
			},
		},
	},
	You: Snake{
		Id:     "yousnakeid",
		Name:   "yousnakename",
		Health: 100,
		Body:   []Coord{Coord{X: 0, Y: 0}},
	},
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

// func TestFindMove(t *testing.T) {
// 	expectedMove := "down"
// 	move := FindMove(testRequest)
// 	if move != expectedMove {
// 		t.Errorf("Incorrect move. Got %v, expected %v", move, expectedMove)
// 	}
// }

// func TestParseMove(t *testing.T) {
// 	world := ParseWorldFromRequest(testRequest)
// 	p, _, _ := Path(world.From(), world.To())
// 	from := Coord{
// 		world.From().X,
// 		world.From().Y,
// 	}

// }

// directions = {(0, -1) : 'left', (0, 1) : 'right', (-1, 0) : 'up', (1, 0) : 'down'}
func TestMoveMap(t *testing.T) {

	up := [2]int{-1, 0}
	down := [2]int{1, 0}
	left := [2]int{0, -1}
	right := [2]int{0, 1}

	upresponse := MoveMap(up)
	if upresponse != "up" {
		t.Errorf("expected up got %v", upresponse)
	}
	downresponse := MoveMap(down)
	if downresponse != "down" {
		t.Errorf("expected down got %v", downresponse)
	}
	leftresponse := MoveMap(left)
	if leftresponse != "left" {
		t.Errorf("expected left got %v", leftresponse)
	}
	rightresponse := MoveMap(right)
	if rightresponse != "right" {
		t.Errorf("expected right got %v", rightresponse)
	}
}

func TestManhattanDistance(t *testing.T) {
	c1 := Coord{0, 0}
	c2 := Coord{5, 5}
	dist := ManhattanDistance(c1, c2)
	expectedDist := 10
	if dist != expectedDist {
		t.Errorf("Expected %v, got %v", expectedDist, dist)
	}
}

func TestFindClosestFood(t *testing.T) {
	expectedResult := testRequest.Board.Food[1]
	result := FindClosestFood(testRequest)
	if result != expectedResult {
		t.Errorf("Expected %v, got %v", expectedResult, result)
	}
}

func TestSortByDistance(t *testing.T) {
	testInput := []Destination{
		Destination{
			Dist: 8,
			Loc:  Coord{X: 0, Y: 0},
		},
		Destination{
			Dist: 1,
			Loc:  Coord{X: 0, Y: 0},
		},
		Destination{
			Dist: 27,
			Loc:  Coord{X: 0, Y: 0},
		},
		Destination{
			Dist: 4,
			Loc:  Coord{X: 0, Y: 0},
		},
		Destination{
			Dist: 16,
			Loc:  Coord{X: 0, Y: 0},
		},
	}
	result := SortByDistance(testInput)
	for i := 1; i < len(result); i++ {
		// t.Logf(strconv.Itoa(result[i].Dist))
		if result[i].Dist < result[i-1].Dist {
			t.Errorf("Sort failed")
		}
	}
}

type FloodFillTestMap struct {
	Map           string
	Coord         Coord
	ExpectedValue int
}

func TestFloodFill(t *testing.T) {
	mapList := []FloodFillTestMap{
		FloodFillTestMap{
			Map: `
.....
.....
.....
.....
.....`,
			Coord:         Coord{X: 0, Y: 0},
			ExpectedValue: 25,
		},
		FloodFillTestMap{
			Map: `
.....
XXX..
..X..
..X..
..X..`,
			Coord:         Coord{X: 2, Y: 0},
			ExpectedValue: 6,
		},
		FloodFillTestMap{
			Map: `
.X...
.X...
.X...
FX...
.....`,
			Coord:         Coord{X: 0, Y: 0},
			ExpectedValue: 3,
		},
		FloodFillTestMap{
			Map: `
.......
....XXX
....F.X
.....XX
.XX.XX.
.XXXX..
.......`,
			Coord:         Coord{X: 2, Y: 5},
			ExpectedValue: 1,
		},
	}

	for i := 0; i < len(mapList); i++ {
		testMap := mapList[i]
		world := ParseWorld(testMap.Map)
		value := FloodFill(world, testMap.Coord)
		if value != testMap.ExpectedValue {
			t.Errorf("Floodfill error: expected %v, got %v", testMap.ExpectedValue, value)
		}
	}
}

func TestDeepCopyWorld(t *testing.T) {
	world := ParseWorldFromRequest(testRequest)
	newWorld := DeepCopyWorld(world)

	t.Logf("%v", world[0][0].Kind)
	t.Logf("%v", newWorld[0][0].Kind)

	if world[0][0].Kind != newWorld[0][0].Kind {
		t.Error("World hasn't been copied correctly")
	}

	world.SetTile(&Tile{
		Kind: KindPlain,
	}, 0, 0)

	t.Logf("%v", world[0][0].Kind)
	t.Logf("%v", newWorld[0][0].Kind)

	if world[0][0].Kind == newWorld[0][0].Kind {
		t.Error("World hasn't been copied correctly")
	}

}
