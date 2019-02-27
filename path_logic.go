package main

import (
	"strconv"
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

	//blocking snake locations
	for _, snake := range request.Board.Snakes {
		for i := 0; i < len(snake.Body)-1; i++ {
			coord := snake.Body[i]
			w.SetTile(&Tile{
				Kind: KindBlocker,
			}, coord.Y, coord.X)
		}
	}

	//marking head
	// coord := request.You.Body[0]
	// w.SetTile(&Tile{
	// 	Kind: KindFrom,
	// }, coord.Y, coord.X)

	return w
}

// Sets the tile at g to a goal
func (w World) SetGoal(g Coord) {
	w.SetTile(&Tile{
		Kind: KindTo,
	}, g.Y, g.X)
}

// Sets the tile at g to a goal
func (w World) SetHead(g Coord) {
	w.SetTile(&Tile{
		Kind: KindFrom,
	}, g.Y, g.X)
}

func (w World) IsEmpty(g Coord) bool {
	return w.Tile(g.Y, g.X).Kind == KindPlain
}

// Sets the tile at g to plain
func (w World) StripGoal(g Coord) {
	w.SetTile(&Tile{
		Kind: KindPlain,
	}, g.Y, g.X)
}

func ParseMove(head Coord, moveCoord Coord) string {
	var direction = [2]int{moveCoord.X - head.X, moveCoord.Y - head.Y}
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
			// fmt.Print(value)s
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

func FloodFill(world World, c Coord) int {
	// deepcopy of world so as not to affect the world the game is using
	w := world
	visited := make(map[string]bool)
	fromTile := w.FirstOfKind(KindFrom)

	if fromTile != nil {
		tileKey := strconv.Itoa(fromTile.X) + strconv.Itoa(fromTile.Y)
		visited[tileKey] = true
		// w.SetTile(&Tile{
		// 	Kind: KindBlocker,
		// }, fromTile.X, fromTile.Y)
	}
	return FloodFillUtil(visited, w, c)
}

func FloodFillUtil(visited map[string]bool, w World, c Coord) int {
	tile := w.Tile(c.Y, c.X)
	tileKey := strconv.Itoa(tile.X) + strconv.Itoa(tile.Y)
	neighbors := tile.PathNeighbors()
	// checking if tile has been visited
	if _, ok := visited[tileKey]; ok {
		return 0
	}
	// if len(neighbors) == 0 {
	// 	return 0
	// }
	// w.SetTile(&Tile{
	// 	Kind: KindBlocker,
	// }, c.Y, c.X)
	result := 1
	visited[tileKey] = true
	for i := 0; i < len(neighbors); i++ {
		neighbor := neighbors[i].(*Tile)
		neighborCoord := Coord{X: neighbor.Y, Y: neighbor.X}
		result = result + FloodFillUtil(visited, w, neighborCoord)
	}
	return result
}
