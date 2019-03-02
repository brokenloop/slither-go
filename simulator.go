package main

import (
	"fmt"
	"math"
	"strconv"
)

// func (g GameRequest) MoveSnake(s Snake, m string) {
// 	g.Board.s
// }

type SimulationResult struct {
	simulationId int
	// whether the snake is alive at the end of this simulation
	alive     bool
	foodEaten int

	// how many moves the simulation ran
	moves int

	// the first move of this simulation
	move string

	// cause of death
	// cause string
	// log string
}

// type CompositeResult struct {
// 	numDeaths    int
// 	numMoves     int
// 	foodEaten    int
// 	snakesKilled int
// }

func FindMoveSimulation(w World, g GameRequest) string {
	myIndex := FindSnakeIndex(g, g.You.Id)
	// headCoord := g.You.Body[0]
	// headTile := w[headCoord.Y][headCoord.X]
	// moves := headTile.GetAvailableMoves()
	allMoves := g.GetAllAvailableMoves(w)
	// results := []SimulationResult{}
	results := make(map[string]map[int]SimulationResult)
	numSnakes := len(g.Board.Snakes)
	numSimulations := 0
	// double check that you can index food like this
	for i := 0; i < len(allMoves[myIndex]); i++ {

		myMove := allMoves[myIndex][i]
		if len(myMove) <= 0 {
			// fmt.Println("DEAD END")
			myMove = "left"
		}
		// for each of the enemies possible moves, do a simulation
		// Doing this should ensure that head to head collisions actually happen in simulations
		// check max length in moves to prevent uneccesary thread spawning?
		for j := 0; j < 4; j++ {
			precursorMoves := make([]string, numSnakes)
			precursorMoves[myIndex] = myMove

			for k := 0; k < numSnakes; k++ {
				if k == myIndex {
					continue
				}
				opponentMoves := allMoves[k]
				var opponentMove string
				if len(opponentMoves) <= 0 {
					opponentMove = "right"
				} else {
					opponentMove = opponentMoves[j%len(opponentMoves)]
				}
				precursorMoves[k] = opponentMove
			}
			numSimulations++
			simulationId := numSimulations
			gCopy := DeepCopyRequest(g)
			wCopy := ParseWorldFromRequest(g)
			result := Simulate2(simulationId, wCopy, gCopy, g.You.Id, myIndex, precursorMoves)

			if results[myMove] == nil {
				results[myMove] = make(map[int]SimulationResult)
			}
			results[myMove][simulationId] = result
		}

	}

	worstResultsPerMove := DecomposeResultsMap(results)
	bestResult := ChooseBestResult(worstResultsPerMove)
	fmt.Println()
	fmt.Println("Move", g.Turn)
	fmt.Println(worstResultsPerMove)
	return bestResult.move
	// resultMap := make(map[string]CompositeResult)
	// for i := 0; i < len(results); i++ {
	// 	result := results[i]
	// 	resultMap = EvaluateResult(resultMap, result)
	// }

	// // choose best move
	// bestResult := SimulationResult{
	// 	alive: false,
	// 	moves: 0,
	// 	move:  "",
	// }
	// for i := 0; i < len(results); i++ {
	// 	result := results[i]
	// 	if result.alive && !bestResult.alive {
	// 		bestResult = result
	// 	} else if result.moves > bestResult.moves {
	// 		bestResult = result
	// 	}
	// }
	// fmt.Println("\n\n" + StringifyWorld(w))
	// fmt.Println(g.You.Name)
	// fmt.Println("RESULT")
	// fmt.Println(results)
	// fmt.Println(bestResult.alive)
	// fmt.Println(bestResult.moves)
	// fmt.Println(bestResult.move)
	// fmt.Println(bestResult)
	// return bestResult.move
}

func FirstResultIsBetter(r1 SimulationResult, r2 SimulationResult) bool {
	// alive
	// if r1.alive && !r2.alive {
	// 	return true
	// } else if !r1.alive && r2.alive {
	// 	return false
	// }
	// if r1.moves > r2.moves {
	// 	return true
	// } else if r1.moves < r2.moves {
	// 	return false
	// }
	return ScoreResult(r1) > ScoreResult(r2)
}

func ScoreResult(r SimulationResult) int {
	return r.foodEaten*200 + r.moves*150
}

func DecomposeResultsMap(resultMap map[string]map[int]SimulationResult) map[string]SimulationResult {
	worstResultPerMove := make(map[string]SimulationResult)
	for move, _ := range resultMap {
		worstResult := SimulationResult{
			alive: true,
			moves: math.MaxInt32,
		}
		for simulationId, _ := range resultMap[move] {
			result := resultMap[move][simulationId]
			if FirstResultIsBetter(worstResult, result) {
				worstResult = result
			}
		}
		worstResultPerMove[move] = worstResult
	}
	return worstResultPerMove
}

func ChooseBestResult(resultMap map[string]SimulationResult) SimulationResult {
	bestResult := SimulationResult{
		alive: false,
		moves: 0,
	}
	for move, _ := range resultMap {
		result := resultMap[move]
		if !FirstResultIsBetter(bestResult, result) {
			bestResult = result
		}
	}
	return bestResult
}

// func ChooseBestResult(resultMap map[string]CompositeResult) string {
// 	bestResult := CompositeResult{}

// }

// func EvaluateResult(resultMap map[string]CompositeResult, r SimulationResult) map[string]CompositeResult {
// 	var compositeResult CompositeResult
// 	value, isPresent := resultMap[r.move]
// 	if isPresent {
// 		compositeResult = value
// 	} else {
// 		compositeResult = CompositeResult{}
// 	}
// 	// compositeResult.foodEaten =  compositeResult.foodEaten + r.
// 	if !r.alive {
// 		compositeResult.numDeaths++
// 	}
// 	compositeResult.numMoves += r.moves
// 	resultMap[r.move] = compositeResult
// 	return resultMap
// }

func (t Tile) GetAvailableMoves() []string {
	tileCoord := Coord{X: t.Y, Y: t.X}
	neighbors := t.PathNeighbors()

	moves := make([]string, len(neighbors))
	for i := 0; i < len(neighbors); i++ {
		moves[i] = ParseMoveFromNeighbor(tileCoord, neighbors[i])
	}
	return moves
}

func (g GameRequest) GetAllAvailableMoves(w World) [][]string {
	result := make([][]string, len(g.Board.Snakes))
	for i := 0; i < len(g.Board.Snakes); i++ {
		snakeHead := g.Board.Snakes[i].Body[0]
		snakeHeadTile := w[snakeHead.Y][snakeHead.X]
		result[i] = snakeHeadTile.GetAvailableMoves()
	}
	return result
}

func ParseMoveFromNeighbor(head Coord, neighbor Pather) string {
	nT := neighbor.(*Tile)
	moveCoord := Coord{X: nT.Y, Y: nT.X}
	move := InternalParseMove(head, moveCoord)
	return move
}

func Simulate2(simulationId int, w World, g GameRequest, myId string, myIndex int, precursorMoves []string) SimulationResult {
	// should this be here?!
	// simulations := ""

	// simulations = simulations + "\n" + firstMove
	// simList := strings.Split(StringifyWorld(w), "\n")
	// fmt.Println(firstMove)
	// fmt.Println()
	// fmt.Println(StringifyWorld(w))

	foodMap := g.MapFood()
	foodEaten := 0
	movesToSimulate := 20
	for j := 1; j < movesToSimulate+1; j++ {

		// fmt.Println()
		// fmt.Println(i)
		// fmt.Println(g)
		for i := 0; i < len(g.Board.Snakes); i++ {
			eat := false
			foodIndex := -1
			// head := g.Board.Snakes[i].Body[0]
			// w.SetHead(head)
			// fmt.Println("HEAD SET")
			// fmt.Println(head)
			var move string
			var found bool
			if j == 1 {
				move = precursorMoves[i]
			} else {
				found, move = HungryMove(w, g, i)
				if found == false {
					move = "right"
				}
			}
			eat, foodIndex = CheckEat(move, g.Board.Snakes[i].Body[i], foodMap)
			g.Board.Snakes[i].Move(move, eat)
			if eat {
				foodEaten++
				g.Board.Snakes[i].Health = 100
				// might get error if food is last in list - have to keep an eye on this
				g.Board.Food = append(g.Board.Food[:foodIndex], g.Board.Food[foodIndex+1:]...)
				foodMap = g.MapFood()
				// fmt.Println("\n\nFOOD LEFT")
				// fmt.Println(len(g.Board.Food))
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
			// fmt.Println("DEAD")
			return SimulationResult{
				alive:        false,
				moves:        j,
				move:         precursorMoves[myIndex],
				foodEaten:    foodEaten,
				simulationId: simulationId,
				// log:   simulations,
			}
		}

	}
	return SimulationResult{
		alive:        true,
		moves:        movesToSimulate,
		move:         precursorMoves[myIndex],
		foodEaten:    foodEaten,
		simulationId: simulationId,
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
	headList := make(map[string][]*Snake)
	for i := 0; i < len(g.Board.Snakes); i++ {

		head := g.Board.Snakes[i].Body[0]
		coordString := head.StringifyCoord()
		// if _, isPresent := headList[coordString]; !isPresent {
		headList[coordString] = append(headList[coordString], &g.Board.Snakes[i])
		// }

		if g.Board.Snakes[i].Health <= 0 {
			// fmt.Println("STARVATION")
			killList[i] = true
		}

		if OutOfBounds(head, g.Board.Width) || w[head.Y][head.X].Kind == KindBlocker {
			// fmt.Println("Out of bounds!")
			killList[i] = true
		}
	}
	// head on collisions
	for coordString, _ := range headList {
		if len(headList[coordString]) > 1 {
			for i := 0; i < len(headList[coordString])-1; i++ {
				s1 := headList[coordString][i]
				s2 := headList[coordString][i+1]
				len1 := len(s1.Body)
				len2 := len(s2.Body)
				if len1 < len2 {
					index := FindSnakeIndex(*g, s1.Id)
					killList[index] = true
				} else if len1 > len2 {
					index := FindSnakeIndex(*g, s2.Id)
					killList[index] = true
				} else {
					index1 := FindSnakeIndex(*g, s1.Id)
					index2 := FindSnakeIndex(*g, s2.Id)
					killList[index1] = true
					killList[index2] = true
				}
			}
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
