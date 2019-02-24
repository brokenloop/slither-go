package main

import (
	"fmt"
	"strings"
)

// Have to encode the world using {y, x} because requests come in the opposite orientation
func ParseWorldFromRequest(request GameRequest) World {
	var grid_size = request.Board.Width

	w := World{}
	for x := 0; x < grid_size; x++ {
		for y := 0; y < grid_size; y++ {
			w.SetTile(&Tile{
				Kind: KindPlain,
			}, y, x)
		}
	}

	for _, snake := range request.Board.Snakes {
		fmt.Println(snake)
		for _, coord := range snake.Body {
			w.SetTile(&Tile{
				Kind: KindBlocker,
			}, coord.Y, coord.X)
		}
	}

	// marking body
	for i, coord := range request.You.Body {
		if i != 0 {
			w.SetTile(&Tile{
				Kind: KindBlocker,
			}, coord.Y, coord.X)
		}
	}

	//marking head
	coord := request.You.Body[0]
	w.SetTile(&Tile{
		Kind: KindFrom,
	}, coord.Y, coord.X)

	// setting goal as first food for testing
	// goalCoord := request.Board.Food[0]
	// w.SetTile(&Tile{
	// 	Kind: KindTo,
	// }, goalCoord.Y, goalCoord.X)

	return w
}

// func SetGoal(w *World, g Coord) {
// 	w.SetTile(&Tile{
// 		Kind: KindTo,
// 	}, g.Y, g.X)
// }

// Sets the tile at g to a goal
func (w World) SetGoal(g Coord) {
	w.SetTile(&Tile{
		Kind: KindTo,
	}, g.Y, g.X)
}

// Sets the tile at g to plain
func (w World) StripGoal(g Coord) {
	w.SetTile(&Tile{
		Kind: KindPlain,
	}, g.Y, g.X)
}

// func ParseMove(head Coord, path []Pather) string {

// 	p := path[len(path)-2]
// 	pT := p.(*Tile)
// 	fmt.Println("MOVES")
// 	fmt.Println(head)
// 	fmt.Println(pT.X, pT.Y)
// 	var direction = [2]int{pT.X - head.X, pT.Y - head.Y}
// 	fmt.Printf("DIRECTION: %v", direction)
// 	return MoveMap(direction)
// }

func ParseMove(head Coord, moveCoord Coord) string {

	fmt.Println("\n\nMOVES")
	fmt.Println(head)
	fmt.Println(moveCoord.X, moveCoord.Y)
	var direction = [2]int{moveCoord.X - head.X, moveCoord.Y - head.Y}
	fmt.Printf("DIRECTION: %v", direction)
	return MoveMap(direction)
}

// needs to be tested - what happens if direction is malformed?
func MoveMap(direction [2]int) string {
	if direction[0] == 0 {
		if direction[1] == 1 {
			return "right"
		} else if direction[1] == -1 {
			return "left"
		}
	} else if direction[0] == 1 {
		return "down"
	} else {
		return "up"
	}
	return "right"
}

func PrintWorld(f func(...interface{}), w World) {
	// func PrintWorld(f func(string, ...interface{}), w World) {
	testres := []string{}
	for i := 0; i < len(w); i++ {
		testres = append(testres, "")
		for j := 0; j < len(w); j++ {
			testres[i] = testres[i] + "0"
		}
	}
	for _, v := range w {
		for _, value := range v {
			testres[value.X] = replaceAtIndex(
				testres[value.X],
				KindRunes[value.Kind],
				value.Y,
			)
		}
	}
	stringResult := strings.Join(testres, "\n")
	f("\n" + stringResult)
}

func StringifyWorld(w World) string {
	// func PrintWorld(f func(string, ...interface{}), w World) {
	testres := []string{}
	for i := 0; i < len(w); i++ {
		testres = append(testres, "")
		for j := 0; j < len(w); j++ {
			testres[i] = testres[i] + "0"
		}
	}
	for _, v := range w {
		for _, value := range v {
			testres[value.X] = replaceAtIndex(
				testres[value.X],
				KindRunes[value.Kind],
				value.Y,
			)
		}
	}
	stringResult := strings.Join(testres, "\n")
	return (stringResult)
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

// directions = {(0, -1) : 'left', (0, 1) : 'right', (-1, 0) : 'up', (1, 0) : 'down'}

func ManhattanDistance(c1 Coord, c2 Coord) int {
	absX := c1.X - c2.X
	absY := c1.Y - c2.Y
	if absX < 0 {
		absX = -absX
	}
	if absY < 0 {
		absY = -absY
	}
	return absX + absY
}
