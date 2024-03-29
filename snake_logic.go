package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Destination struct {
	Dist int
	Loc  Coord
}

func FindClosestFood(g GameRequest, snakeIndex int) Coord {
	// snakeHead := g.You.Body[0]
	snakeHead := g.Board.Snakes[snakeIndex].Body[0]
	availableFood := g.Board.Food
	shortestDistance := math.MaxInt32
	result := Coord{-1, -1}
	for i := 0; i < len(availableFood); i++ {
		food := availableFood[i]
		dist := ManhattanDistance(snakeHead, food)
		if dist < shortestDistance {
			shortestDistance = dist
			result = food
		}
	}
	// fmt.Printf("Shortest distance: %d", shortestDistance)
	// fmt.Printf("Closest food: %d", result)
	return result
}

func ListFoodDistances(head Coord, foodList []Coord) []Destination {
	result := []Destination{}
	for i := 0; i < len(foodList); i++ {
		food := foodList[i]
		distance := ManhattanDistance(head, food)
		newDest := Destination{
			Dist: distance,
			Loc:  food,
		}
		result = append(result, newDest)
	}
	return result
}

// quicksort list of destinations based on their distance
func SortByDistance(a []Destination) []Destination {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]

	for i, _ := range a {
		if a[i].Dist < a[right].Dist {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]

	SortByDistance(a[:left])
	SortByDistance(a[left+1:])

	return a
}

func PrintFindMove(f func(...interface{}), g GameRequest) {
	world := ParseWorldFromRequest(g)
	p, _, found := Path(world.From(), world.To())
	if !found {
		f("Could not find a path")
	} else {
		f("Resulting path\n", world.RenderPath(p))
	}
}

func LastResort(w World, g GameRequest, snakeIndex int) string {
	// fmt.Println("\nLAST RESORT")
	// head := g.Board.Snakes[snakeIndex].Body[0]
	// headTile := w[head.Y][head.X]
	// // headTile := w.From()
	// // head := Coord{
	// // 	X: headTile.X,
	// // 	Y: headTile.Y,
	// // }
	// neighbors := headTile.PathNeighbors()

	// if len(neighbors) > 0 {
	// 	fmt.Println("Possible moves: ", len(neighbors))
	// 	placeHolder := Coord{}
	// 	for i := 0; i < len(neighbors); i++ {
	// 		neighbor := neighbors[i]

	// 		nT := neighbor.(*Tile)
	// 		moveCoord := Coord{X: nT.Y, Y: nT.X}
	// 		placeHolder = moveCoord
	// 		if TileSafe(moveCoord, g, snakeIndex) {
	// 			move := ParseMove(head, moveCoord)
	// 			return move
	// 		}
	// 		fmt.Println("FALLBACK")
	// 	}
	// 	// flippedHead := FlipCoords(head)
	// 	flippedPH := FlipCoords(placeHolder)
	// 	move := ParseMove(head, flippedPH)
	// 	return move
	// }
	// fmt.Println("\nDead end!")
	// return "right"
	return g.Board.Snakes[snakeIndex].RandomMove(w)
}

func TileSafe(tile Coord, g GameRequest, snakeIndex int) bool {
	for i := 0; i < len(g.Board.Snakes); i++ {
		if i != snakeIndex {
			snake := g.Board.Snakes[i]
			snakeHead := snake.Body[0]
			snakeHead = FlipCoords(snakeHead)
			distance := ManhattanDistance(tile, snakeHead)
			fmt.Printf("\nOpponent snake: %v", snake)
			fmt.Printf("\nDestination: %v", tile)
			fmt.Printf("\nOpponent Head: %v", snakeHead)
			fmt.Printf("\nDistance: %v", distance)
			if distance <= 1 && len(snake.Body) >= len(g.Board.Snakes[snakeIndex].Body) {
				fmt.Println("\n\nTILE NOT SAFE")
				return false
			}
		}
	}
	return true
}

func FlipCoords(coord Coord) Coord {
	return Coord{X: coord.Y, Y: coord.X}
}

func TrappingMove(w World, g GameRequest, snakeIndex int, c Coord) bool {
	flippedC := FlipCoords(c)
	availableSpace := FloodFill(w, flippedC)
	// fmt.Println("\nCHECKING FOR TRAPS")
	// fmt.Println(c)
	// fmt.Println(availableSpace)

	return availableSpace < len(g.Board.Snakes[snakeIndex].Body)
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

func IdleMove(w World, g GameRequest, snakeIndex int) (bool, string) {
	fmt.Println("\nIDLE MOVE")
	head := g.Board.Snakes[snakeIndex].Body[0]
	headTile := w[head.Y][head.X]
	// headTile := w.From()
	// head := Coord{
	// 	X: headTile.X,
	// 	Y: headTile.Y,
	// }
	neighbors := headTile.PathNeighbors()

	if len(neighbors) > 0 {
		largestSpace := 0
		firstOption := neighbors[0].(*Tile)
		bestMove := Coord{X: firstOption.X, Y: firstOption.Y}
		for i := 0; i < len(neighbors); i++ {
			neighbor := neighbors[i]

			nT := neighbor.(*Tile)
			moveCoord := Coord{X: nT.X, Y: nT.Y}

			if TileSafe(moveCoord, g, snakeIndex) {
				flippedC := FlipCoords(moveCoord)
				availableSpace := FloodFill(w, flippedC)
				fmt.Println("Move: ", moveCoord)
				fmt.Println("Available space: ", availableSpace)
				if availableSpace > largestSpace {
					largestSpace = availableSpace
					bestMove = moveCoord
				}
			}
		}
		flippedHead := FlipCoords(head)
		move := ParseMove(flippedHead, bestMove)
		return true, move
	}
	fmt.Println("DEAD END")
	return false, "right"
}

// Try to eat food
func HungryMove(world World, g GameRequest, snakeIndex int) (bool, string) {
	foundMove := false
	// head := g.You.Body[0]
	head := g.Board.Snakes[snakeIndex].Body[0]
	foodList := ListFoodDistances(head, g.Board.Food)
	foodList = SortByDistance(foodList)
	for i := 0; i < len(foodList); i++ {
		goalCoord := foodList[i].Loc
		world.SetGoal(goalCoord)
		world.SetHead(head)
		if head == goalCoord {
			continue
		}
		// fromTile := world.From()
		// fmt.Println("fromTile")
		// fmt.Println(fromTile)
		// toTile := world.To()
		// fmt.Println("toTile")
		// fmt.Println(toTile)
		p, _, found := Path(world.From(), world.To())
		if found {
			head := Coord{
				X: world.From().X,
				Y: world.From().Y,
			}
			cutPath := p[len(p)-2]
			moveTile := cutPath.(*Tile)
			moveCoord := Coord{X: moveTile.X, Y: moveTile.Y}
			if TileSafe(moveCoord, g, snakeIndex) && !TrappingMove(world, g, snakeIndex, moveCoord) {
				// fmt.Println("SAFE MOVE")
				move := ParseMove(head, moveCoord)
				foundMove = true
				return foundMove, move
			}
		}
		world.StripGoal(goalCoord)
		world.StripHead(head)
	}
	return foundMove, ""
}

// Chase own tail
func ScaredyMove(world World, g GameRequest, snakeIndex int) (bool, string) {
	foundMove := false
	// tail := g.You.Body[len(g.You.Body)-1]
	tail := g.Board.Snakes[snakeIndex].Body[len(g.Board.Snakes[snakeIndex].Body)-1]

	// if head and tail are on same tile, default to something else
	// only relevent for first move
	if !world.IsEmpty(tail) {
		return foundMove, ""
	}
	fmt.Printf("TAILCOORD %v", tail)
	world.SetGoal(tail)
	fmt.Println()
	fmt.Printf(StringifyWorld(world))
	p, _, found := Path(world.From(), world.To())
	world.StripGoal(tail)
	if found {
		head := Coord{
			X: world.From().X,
			Y: world.From().Y,
		}
		cutPath := p[len(p)-2]
		moveTile := cutPath.(*Tile)
		moveCoord := Coord{X: moveTile.X, Y: moveTile.Y}
		if TileSafe(moveCoord, g, snakeIndex) {
			move := ParseMove(head, moveCoord)
			foundMove = true
			return foundMove, move
		}
	}
	return foundMove, ""
}

func FindSnakeIndex(g GameRequest, snakeId string) int {

	index := -1
	for i := 0; i < len(g.Board.Snakes); i++ {
		if g.Board.Snakes[i].Id == snakeId {
			index = i
		}
	}
	return index
}

func FindMove(g GameRequest) string {
	snakeIndex := FindSnakeIndex(g, g.You.Id)
	world := ParseWorldFromRequest(g)
	world.SetHead(g.You.Body[0])
	fmt.Println("\n\n\nworld")
	fmt.Printf(StringifyWorld(world))
	foundMove := false
	move := ""
	health := g.You.Health
	if health > 50 {
		foundMove, move = ScaredyMove(world, g, snakeIndex)
		if foundMove {
			return move
		}
	}
	foundMove, move = HungryMove(world, g, snakeIndex)
	if foundMove {
		return move
	}
	foundMove, move = IdleMove(world, g, snakeIndex)
	if foundMove {
		return move
	}
	// is this needed?
	return LastResort(world, g, snakeIndex)
}

func FindMoveBySimulation(g GameRequest) string {
	world := ParseWorldFromRequest(g)
	return FindMoveSimulation(world, g)
}
