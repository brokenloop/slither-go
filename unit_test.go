package main

import (
	"fmt"
	"testing"
	"time"
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
			Snake{
				Id:     "yousnakeid",
				Name:   "yousnakename",
				Health: 100,
				Body:   []Coord{Coord{X: 0, Y: 1}},
			},
		},
	},
	You: Snake{
		Id:     "yousnakeid",
		Name:   "yousnakename",
		Health: 100,
		Body:   []Coord{Coord{X: 0, Y: 1}},
	},
}

func TestParseWorldFromRequest(t *testing.T) {
	world := ParseWorldFromRequest(testRequest)
	world.SetHead(testRequest.You.Body[0])
	fromTile := world.From()

	expectedFrom := testRequest.You.Body[0]

	if fromTile.X != expectedFrom.Y || fromTile.Y != expectedFrom.X {
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
	snakeIndex := 0
	result := FindClosestFood(testRequest, snakeIndex)
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
	newWorld := world.DeepCopyWorld()

	if world[0][1].Kind != newWorld[0][1].Kind {
		t.Errorf("World hasn't been copied correctly. Expected %v got %v", world[0][1].Kind, newWorld[0][1].Kind)
	}

	world.SetTile(&Tile{
		Kind: KindMountain,
	}, 0, 1)

	t.Log(StringifyWorld(world))
	t.Log(StringifyWorld(newWorld))
	if world[0][1].Kind == newWorld[0][1].Kind {
		t.Errorf("World hasn't been copied correctly. Expected %v got %v", world[0][1].Kind, newWorld[0][1].Kind)
	}
}

func TestMoveSnake(t *testing.T) {
	// world := ParseWorldFromRequest(testRequest)
	snake := testRequest.You
	snake.Move("down", false)
	expectedResult := Coord{X: 0, Y: 2}
	result := snake.Body[0]
	if result != expectedResult {
		t.Errorf("Move is broken, expected %v got %v", expectedResult, result)
	}

	if len(snake.Body) != 1 {
		t.Error("snake grew and didn't eat")
	}

	snake.Move("down", true)
	t.Logf("%v", snake.Body)

	if len(snake.Body) != 2 {
		t.Error("snake ate and didn't grow")
	}

}

func TestRandomMove(t *testing.T) {
	world := ParseWorldFromRequest(testRequest)
	t.Log("\n" + StringifyWorld(world))
	snake := testRequest.You
	t.Log(snake.Body[0])
	for i := 0; i < 10; i++ {
		t.Logf("%v", snake.RandomMove(world))
	}
}

func TestOutOfBounds(t *testing.T) {
	maxSize := 10
	testCoords := []Coord{
		Coord{X: -1, Y: 0},
		Coord{X: 0, Y: -1},
		Coord{X: 10, Y: 0},
		Coord{X: 0, Y: 10},
	}
	for i := 0; i < len(testCoords); i++ {
		if !OutOfBounds(testCoords[i], maxSize) {
			t.Errorf("Broken! %v", testCoords[i])
		}
	}

	safeCoord := Coord{X: 0, Y: 0}
	if OutOfBounds(safeCoord, maxSize) {
		t.Errorf("Broken! %v", safeCoord)
	}
}

var deadRequest = GameRequest{
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
				Id:     "self-collision",
				Name:   "",
				Health: 100,
				Body:   []Coord{Coord{X: 3, Y: 3}, Coord{X: 3, Y: 3}, Coord{X: 3, Y: 2}},
			},
			Snake{
				Id:     "head on 1",
				Name:   "",
				Health: 100,
				Body:   []Coord{Coord{X: 5, Y: 5}, Coord{X: 5, Y: 4}},
			},
			Snake{
				Id:     "head on 2",
				Name:   "",
				Health: 100,
				Body:   []Coord{Coord{X: 5, Y: 5}, Coord{X: 5, Y: 6}},
			},
			Snake{
				Id:     "head on 3",
				Name:   "",
				Health: 100,
				Body:   []Coord{Coord{X: 5, Y: 5}},
			},
			Snake{
				Id:     "safe",
				Name:   "",
				Health: 100,
				Body:   []Coord{Coord{X: 2, Y: 2}},
			},
		},
	},
	// You: Snake{
	// 	Id:     "yousnakeid",
	// 	Name:   "yousnakename",
	// 	Health: 100,
	// 	Body:   []Coord{Coord{X: 0, Y: 1}},
	// },
}

func TestKillSnakes(t *testing.T) {
	world := ParseWorldFromRequest(deadRequest)
	deadRequest.KillSnakes(world)
	expectedAlive := 1
	alive := len(deadRequest.Board.Snakes)
	if alive != expectedAlive {
		t.Errorf("Kill broken: expected %v got %v", expectedAlive, alive)
	}
	// Simulate(world, testRequest)
}

var moveSimRequest = GameRequest{
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
				Body:   []Coord{Coord{X: 3, Y: 3}, Coord{X: 3, Y: 4}, Coord{X: 3, Y: 5}},
			},
			Snake{
				Id:     "yousnakeid",
				Name:   "yousnakename",
				Health: 100,
				Body:   []Coord{Coord{X: 0, Y: 1}, Coord{X: 1, Y: 1}, Coord{X: 1, Y: 0}, Coord{X: 2, Y: 0}, Coord{X: 3, Y: 0}},
			},
		},
	},
	You: Snake{
		Id:     "yousnakeid",
		Name:   "yousnakename",
		Health: 100,
		Body:   []Coord{Coord{X: 0, Y: 1}, Coord{X: 1, Y: 1}, Coord{X: 1, Y: 0}, Coord{X: 2, Y: 0}, Coord{X: 3, Y: 0}},
	},
}

var moveSimRequest2 = GameRequest{
	Game: Game{
		Id: "gameid",
	},
	Turn: 0,
	Board: Board{
		Height: 5,
		Width:  5,
		Food:   []Coord{Coord{X: 0, Y: 0}, Coord{X: 0, Y: 1}, Coord{X: 0, Y: 2}},
		Snakes: []Snake{
			Snake{
				Id:     "yousnakeid",
				Name:   "yousnakename",
				Health: 100,
				Body: []Coord{Coord{X: 0, Y: 3},
					Coord{X: 1, Y: 3},
					Coord{X: 1, Y: 2},
					Coord{X: 1, Y: 1},
					Coord{X: 1, Y: 0},
					Coord{X: 2, Y: 0},
					Coord{X: 2, Y: 1},
					Coord{X: 2, Y: 2},
					Coord{X: 2, Y: 3},
					Coord{X: 2, Y: 4},
					Coord{X: 2, Y: 4},
					Coord{X: 2, Y: 4},
					Coord{X: 2, Y: 4},
					Coord{X: 2, Y: 4},
					Coord{X: 2, Y: 4}},
			},
		},
	},
	You: Snake{
		Id:     "yousnakeid",
		Name:   "yousnakename",
		Health: 100,
		Body: []Coord{Coord{X: 0, Y: 3},
			Coord{X: 1, Y: 3},
			Coord{X: 1, Y: 2},
			Coord{X: 1, Y: 1},
			Coord{X: 1, Y: 0},
			Coord{X: 2, Y: 0},
			Coord{X: 2, Y: 1},
			Coord{X: 2, Y: 2},
			Coord{X: 2, Y: 3},
			Coord{X: 2, Y: 4},
			Coord{X: 2, Y: 4},
			Coord{X: 2, Y: 4},
			Coord{X: 2, Y: 4},
			Coord{X: 2, Y: 4},
			Coord{X: 2, Y: 4}},
	},
}

// func TestDeepCopyRequest(t *testing.T) {
// 	original := moveSimRequest
// 	copy := original.DeepCopyRequest()
// 	snakes := copy.Board.Snakes
// 	for i := 0; i < len(snakes); i++ {
// 		origSnake := original.Board.Snakes[i]
// 		newSnake := copy.Board.Snakes[i]
// 		for j := 0; j < len(snakes); j++ {
// 			c1 := origSnake.Body[i]
// 			c2 := newSnake.Body[i]
// 			fmt.Println(c1)
// 			fmt.Println(c2)

// 		}

// 		fmt.Println(origSnake)
// 		fmt.Println(newSnake)
// 	}
// 	fmt.Println(copy)
// 	fmt.Println(snakes)
// 	fmt.Println(original)
// 	copy.You.Body[0] = Coord{X: 99, Y: 99}
// 	if copy.You.Body[0] == original.You.Body[0] {
// 		t.Log("Universal rules broken")
// 	}
// }

func TestFindMoveSimulation(t *testing.T) {
	world := ParseWorldFromRequest(moveSimRequest2)
	bestMove := FindMoveSimulation(world, moveSimRequest2)
	t.Log(bestMove)
	// deadRequest.KillSnakes(world)
	// expectedAlive := 1
	// alive := len(deadRequest.Board.Snakes)
	// if alive != expectedAlive {
	// 	t.Errorf("Kill broken: expected %v got %v", expectedAlive, alive)
	// }
	// Simulate(world, testRequest)
}

func TestWorldGenTime(t *testing.T) {
	start := time.Now()
	newS := GameRequest{}
	newWorld := World{}
	for i := 0; i < 100000; i++ {
		// newS = DeepCopyRequest(moveSimRequest)
		newWorld = ParseWorldFromRequest(moveSimRequest)

		newS = moveSimRequest
	}

	fmt.Println(newS)
	fmt.Println(newWorld)
	fmt.Println(time.Since(start))
}

func TestBlah(t *testing.T) {
	w := ParseWorldFromRequest(moveSimRequest)
	bestMove := FindMoveSimulation(w, moveSimRequest)
	fmt.Println(bestMove)
}

var multiAgentRequest = GameRequest{
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
				Body:   []Coord{Coord{X: 3, Y: 3}, Coord{X: 3, Y: 3}, Coord{X: 3, Y: 3}},
			},
			Snake{
				Id:     "themsnakeid2",
				Name:   "themsnakename2",
				Health: 100,
				Body:   []Coord{Coord{X: 0, Y: 3}, Coord{X: 0, Y: 3}, Coord{X: 0, Y: 3}},
			},
			Snake{
				Id:     "yousnakeid",
				Name:   "yousnakename",
				Health: 100,
				Body:   []Coord{Coord{X: 0, Y: 1}, Coord{X: 0, Y: 1}, Coord{X: 0, Y: 1}},
			},
		},
	},
	You: Snake{
		Id:     "yousnakeid",
		Name:   "yousnakename",
		Health: 100,
		Body:   []Coord{Coord{X: 0, Y: 1}, Coord{X: 0, Y: 1}, Coord{X: 0, Y: 1}},
	},
}

func TestMultiagentSimulation(t *testing.T) {
	world := ParseWorldFromRequest(multiAgentRequest)
	bestMove := FindMoveSimulation(world, multiAgentRequest)
	t.Log(bestMove)
}
