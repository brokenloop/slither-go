package main

import (
	"fmt"
	"strconv"
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
	// log string
}

func FindMoveSimulation(w World, g GameRequest) string {
	headCoord := g.You.Body[0]
	headTile := w[headCoord.Y][headCoord.X]
	neighbors := headTile.PathNeighbors()
	results := []SimulationResult{}
	for i := 0; i < len(neighbors); i++ {
		gCopy := DeepCopyRequest(g)
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
	// fmt.Println("\n\n" + StringifyWorld(w))
	// fmt.Println(g.You.Name)
	// fmt.Println("RESULT")
	// fmt.Println(results)
	// fmt.Println(bestResult.alive)
	// fmt.Println(bestResult.moves)
	// fmt.Println(bestResult.move)
	fmt.Println(bestResult)
	return bestResult.move
}

func ParseMoveFromNeighbor(head Coord, neighbor Pather) string {
	nT := neighbor.(*Tile)
	moveCoord := Coord{X: nT.Y, Y: nT.X}
	move := InternalParseMove(head, moveCoord)
	return move
}

func Simulate2(w World, g GameRequest, myId string, firstMove string) SimulationResult {
	// should this be here?!

	// simulations := ""

	// simulations = simulations + "\n" + firstMove
	// simList := strings.Split(StringifyWorld(w), "\n")
	// fmt.Println(firstMove)
	// fmt.Println()
	// fmt.Println(StringifyWorld(w))

	foodMap := g.MapFood()
	movesToSimulate := 5
	for j := 1; j < movesToSimulate+1; j++ {

		// fmt.Println()
		// fmt.Println(i)
		// fmt.Println(g)
		for i := 0; i < len(g.Board.Snakes); i++ {
			eat := false
			foodIndex := -1
			head := g.Board.Snakes[i].Body[0]
			w.SetHead(head)
			// fmt.Println("HEAD SET")
			// fmt.Println(head)
			if g.Board.Snakes[i].Id == myId {
				// fmt.Println("My Id")
				// fmt.Println(myId)
				if j == 1 {
					eat, foodIndex = CheckEat(firstMove, g.Board.Snakes[i].Body[i], foodMap)
					g.Board.Snakes[i].Move(firstMove, eat)
				} else {
					// move := g.Board.Snakes[i].RandomMove(w)
					found, move := HungryMove(w, g, i)
					if found == false {
						move = "right"
					}
					eat, foodIndex = CheckEat(move, g.Board.Snakes[i].Body[i], foodMap)
					g.Board.Snakes[i].Move(move, eat)
				}
			} else {
				found, move := HungryMove(w, g, i)
				if found == false {
					move = "right"
				}
				eat, foodIndex = CheckEat(move, g.Board.Snakes[i].Body[i], foodMap)
				g.Board.Snakes[i].Move(move, eat)
			}
			if eat {
				g.Board.Snakes[i].Health = 100
				// might get error if food is last in list - have to keep an eye on this
				if foodIndex >= len(g.Board.Food) {
					g.Board.Food = g.Board.Food[:len(g.Board.Food)]
				} else {
					g.Board.Food = append(g.Board.Food[:foodIndex], g.Board.Food[foodIndex+1:]...)
				}
				foodMap = g.MapFood()
				fmt.Println("\n\nFOOD LEFT")
				fmt.Println(len(g.Board.Food))
			} else {
				g.Board.Snakes[i].Health--
			}
			// w.StripHead(head)
		}
		g.KillSnakes(w)
		w = ParseWorldFromRequest(g)
		// simulations = simulations + "\n\n" + StringifyWorld(w)
		// simList = strings.Split(StringifyWorld(w), "\n")
		// fmt.Println(simList)
		// fmt.Println(StringifyWorld(w))
		if !g.SnakeAlive(myId) {
			fmt.Println("DEAD")
			return SimulationResult{
				alive: false,
				moves: j,
				move:  firstMove,
				// log:   simulations,
			}
		}

	}
	return SimulationResult{
		alive: true,
		moves: movesToSimulate,
		move:  firstMove,
		// log:   simulations,
	}
}

func CheckEat(move string, head Coord, foodMap map[string]int) (bool, int) {
	direction := InternalParseDirection(move)
	moveCoord := Coord{X: head.X + direction[0], Y: head.Y + direction[1]}
	coordString := moveCoord.StringifyCoord()
	if index, isPresent := foodMap[coordString]; isPresent {
		return true, index
	}
	return false, -1
}

func (g *GameRequest) MapFood() map[string]int {
	food := make(map[string]int)
	for i := 0; i < len(g.Board.Food); i++ {
		foodString := g.Board.Food[i].StringifyCoord()
		food[foodString] = i
	}
	return food
}

func (g *GameRequest) SnakeAlive(id string) bool {
	for i := 0; i < len(g.Board.Snakes); i++ {
		if g.Board.Snakes[i].Id == id {
			return true
		}
	}
	return false
}

func (g *GameRequest) KillSnakes(w World) {
	killList := make(map[int]bool)
	headList := make(map[string]*Snake)
	for i := 0; i < len(g.Board.Snakes); i++ {

		head := g.Board.Snakes[i].Body[0]
		if g.Board.Snakes[i].Health <= 0 {
			fmt.Println("STARVATION")
			killList[i] = true
		}

		// checking for head on collisions - have to do this twice to settle snakes of equal size
		for j := 0; j < 2; j++ {
			coordString := head.StringifyCoord()
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

func (c Coord) StringifyCoord() string {
	return strconv.Itoa(c.X) + strconv.Itoa(c.Y)
	// return string([]rune{rune(c.X), rune(c.Y)})
}

func OutOfBounds(head Coord, size int) bool {
	return head.X < 0 || head.X >= size || head.Y < 0 || head.Y >= size
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

func DeepCopyRequest(g GameRequest) GameRequest {
	ng := GameRequest{
		Game:  g.Game,
		Turn:  g.Turn,
		You:   DeepCopySnake(g.You),
		Board: DeepCopyBoard(g.Board),
	}
	return ng
}

func DeepCopySnake(s Snake) Snake {
	ns := Snake{
		Id:     s.Id,
		Name:   s.Name,
		Health: s.Health,
		Body:   make([]Coord, len(s.Body)),
	}
	for i := 0; i < len(s.Body); i++ {
		ns.Body[i] = DeepCopyCoord(s.Body[i])
	}
	return ns
}

func DeepCopyBoard(b Board) Board {
	nb := Board{
		Height: b.Height,
		Width:  b.Width,
		Food:   make([]Coord, len(b.Food)),
		Snakes: make([]Snake, len(b.Snakes)),
	}
	for i := 0; i < len(b.Food); i++ {
		nb.Food[i] = DeepCopyCoord(b.Food[i])
	}
	for i := 0; i < len(b.Snakes); i++ {
		nb.Snakes[i] = DeepCopySnake(b.Snakes[i])
	}
	return nb
}

func DeepCopyCoord(c Coord) Coord {
	nc := Coord{
		X: c.X,
		Y: c.Y,
	}
	return nc
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
