package main

import (
	"fmt"
	"math/rand"
)

// func (g GameRequest) MoveSnake(s Snake, m string) {
// 	g.Board.s
// }

type SimulationResult struct {
	// whether the snake is alive at the end of this simulation
	alive bool

	// how many moves the simulation ran
	moves int

	// the first move of this simulation
	move string

	// cause of death
	// cause string
	log string
}

func FindMoveSimulation(w World, g GameRequest) string {
	headCoord := g.You.Body[0]
	headTile := w[headCoord.Y][headCoord.X]
	neighbors := headTile.PathNeighbors()
	results := []SimulationResult{}
	for i := 0; i < len(neighbors); i++ {
		// foodEaten := 0
		// alive := true
		gCopy := g
		wCopy := ParseWorldFromRequest(g)
		myMove := ParseMoveFromNeighbor(headCoord, neighbors[i])
		if len(myMove) < 0 {
			fmt.Println("DEAD END")
			myMove = "left"
		}
		result := Simulate2(wCopy, gCopy, g.You.Id, myMove)
		results = append(results, result)
	}

	// choose best move
	bestResult := SimulationResult{
		alive: false,
		moves: 0,
		move:  "",
	}
	for i := 0; i < len(results); i++ {
		result := results[i]
		if result.alive && !bestResult.alive {
			bestResult = result
		} else if result.moves > bestResult.moves {
			bestResult = result
		}
	}
	fmt.Println("\n\n" + StringifyWorld(w))
	fmt.Println(g.You.Name)
	fmt.Println("RESULT")
	// fmt.Println(results)
	fmt.Println(bestResult.alive)
	fmt.Println(bestResult.moves)
	fmt.Println(bestResult.move)
	fmt.Println(bestResult.log)
	return bestResult.move
}

func (g GameRequest) CopySnake() {
	myId := g.You.Id
	for i := 0; i < len(g.Board.Snakes); i++ {
		if g.Board.Snakes[i].Id == myId {
			g.Board.Snakes[i] = g.You
		}
	}
}

func ParseMoveFromNeighbor(head Coord, neighbor Pather) string {
	nT := neighbor.(*Tile)
	moveCoord := Coord{X: nT.Y, Y: nT.X}
	move := InternalParseMove(head, moveCoord)
	return move
}

func Simulate2(w World, g GameRequest, myId string, firstMove string) SimulationResult {
	// should this be here?!

	simulations := ""

	simulations = simulations + "\n" + firstMove
	// fmt.Println(firstMove)
	// fmt.Println()
	// fmt.Println(StringifyWorld(w))
	movesToSimulate := 10
	for j := 1; j < movesToSimulate+1; j++ {

		// fmt.Println()
		// fmt.Println(i)
		// fmt.Println(g)
		for i := 0; i < len(g.Board.Snakes); i++ {
			if g.Board.Snakes[i].Id == myId {
				// fmt.Println("My Id")
				// fmt.Println(myId)
				if i == 0 {
					g.Board.Snakes[i].Move(firstMove, false)
				} else {
					move := g.Board.Snakes[i].RandomMove(w)
					if move == "" {
						move = "right"
					}
					g.Board.Snakes[i].Move(move, false)
				}
			} else {
				move := g.Board.Snakes[i].RandomMove(w)
				if move == "" {
					move = "right"
				}
				g.Board.Snakes[i].Move(move, false)
			}
		}
		g.KillSnakes(w)
		w = ParseWorldFromRequest(g)
		simulations = simulations + "\n\n" + StringifyWorld(w)
		// fmt.Println(StringifyWorld(w))
		if !g.SnakeAlive(myId) {
			fmt.Println("DEAD")
			return SimulationResult{
				alive: false,
				moves: j,
				move:  firstMove,
				log:   simulations,
			}
		}

	}
	return SimulationResult{
		alive: true,
		moves: movesToSimulate,
		move:  firstMove,
		log:   simulations,
	}
}

// func InternalMoveToExternal(m string) string {
// 	result := ""
// 	switch m {
// 	case "left":
// 		result = "up"
// 		// result = [2]int{-1, 0}
// 	case "right":
// 		result = "down"

// 		// result = [2]int{1, 0}
// 	case "up":
// 		result = "left"

// 		// result = [2]int{0, -1}
// 	case "down":
// 		result = "right"
// 		// result = [2]int{0, 1}
// 	}
// 	return result
// }

// func (g GameRequest) CheckEat(c Coord) bool {

// }

func (g *GameRequest) SnakeAlive(id string) bool {
	for i := 0; i < len(g.Board.Snakes); i++ {
		if g.Board.Snakes[i].Id == id {
			return true
		}
	}
	return false
}

func Simulate(w World, g GameRequest) {
	for i := 0; i < 10; i++ {
		if len(g.Board.Snakes) > 0 {
			// fmt.Println(StringifyWorld(w))
			// fmt.Println(i)
			// fmt.Println(g)
			for i := 0; i < len(g.Board.Snakes); i++ {
				move := g.Board.Snakes[i].RandomMove(w)
				g.Board.Snakes[i].Move(move, true)
				// g.Board.Snakes[i].Move("down", false)
			}
			g.KillSnakes(w)
			w = ParseWorldFromRequest(g)
		}
	}
}

func (g *GameRequest) KillSnakes(w World) {
	killList := make(map[int]bool)
	headList := make(map[string]*Snake)
	for i := 0; i < len(g.Board.Snakes); i++ {

		head := g.Board.Snakes[i].Body[0]

		// checking for head on collisions - have to do this twice to settle snakes of equal size
		for j := 0; j < 2; j++ {
			coordString := string([]rune{rune(head.X), rune(head.Y)})
			// fmt.Println("coord string")
			// fmt.Println(coordString)
			if _, isPresent := headList[coordString]; isPresent {
				sameSnake := headList[coordString].Id == g.Board.Snakes[i].Id
				// fmt.Println("Same snake")
				// fmt.Println(sameSnake)
				size := len(headList[coordString].Body)
				if !sameSnake && size >= len(g.Board.Snakes[i].Body) {
					// fmt.Println("Head on collision")
					killList[i] = true
				} else {
					headList[coordString] = &g.Board.Snakes[i]
				}
			} else {
				headList[coordString] = &g.Board.Snakes[i]
			}
		}

		if OutOfBounds(head, g.Board.Width) || w[head.Y][head.X].Kind == KindBlocker {
			// fmt.Println("Out of bounds!")
			killList[i] = true
		}
	}
	if len(killList) > 0 {
		newSnakeList := []Snake{}
		for i := 0; i < len(g.Board.Snakes); i++ {
			// snake isn't in killList
			if _, isPresent := killList[i]; !isPresent {
				newSnakeList = append(newSnakeList, g.Board.Snakes[i])
			}
		}
		g.Board.Snakes = newSnakeList
	}
	// fmt.Println(headList)
}

func OutOfBounds(head Coord, size int) bool {
	return head.X < 0 || head.X >= size || head.Y < 0 || head.Y >= size
}

func (s *Snake) RandomMove(w World) string {
	head := s.Body[0]
	// fmt.Print(w[head.Y][head.X].Kind)
	neighbors := w[head.Y][head.X].PathNeighbors()
	if len(neighbors) > 0 {
		neighbor := neighbors[rand.Intn(len(neighbors))]
		nT := neighbor.(*Tile)
		moveCoord := Coord{X: nT.Y, Y: nT.X}

		// fmt.Print("\n")
		// fmt.Print(moveCoord)
		move := InternalParseMove(head, moveCoord)
		return move
	}
	return "right"
}

func (oldWorld World) DeepCopyWorld() World {
	w := World{}
	for x, row := range oldWorld {
		for y, tile := range row {
			kind := tile.Kind
			w.SetTile(&Tile{
				Kind: kind,
			}, x, y)
		}
	}
	return w
}

func (s *Snake) Move(m string, eat bool) {
	direction := InternalParseDirection(m)
	oldHead := s.Body[0]
	newHead := Coord{X: oldHead.X + direction[0], Y: oldHead.Y + direction[1]}
	if eat {
		s.Body = append([]Coord{newHead}, s.Body...)
	} else {
		s.Body = append([]Coord{newHead}, s.Body[:len(s.Body)-1]...)
	}

}

func InternalParseDirection(m string) [2]int {
	result := [2]int{0, 0}
	switch m {
	case "left":
		result = [2]int{-1, 0}
	case "right":
		result = [2]int{1, 0}
	case "up":
		result = [2]int{0, -1}
	case "down":
		result = [2]int{0, 1}
	}
	return result
}

func InternalParseMove(head Coord, moveCoord Coord) string {
	var direction = [2]int{moveCoord.X - head.X, moveCoord.Y - head.Y}
	return InternalMoveMap(direction)
}

// needs to be tested - what happens if direction is malformed?
func InternalMoveMap(direction [2]int) string {
	if direction[0] == 0 {
		if direction[1] == 1 {
			return "down"
		} else if direction[1] == -1 {
			return "up"
		}
	} else if direction[0] == 1 {
		return "right"
	} else {
		return "left"
	}
	return "right"
}
