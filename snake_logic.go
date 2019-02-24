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

func FindClosestFood(g GameRequest) Coord {
	snakeHead := g.You.Body[0]
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

func LastResort(w World) string {
	fmt.Println("\nLAST RESORT")
	headTile := w.From()
	head := Coord{
		X: headTile.X,
		Y: headTile.Y,
	}
	neighbors := headTile.PathNeighbors()
	if len(neighbors) > 0 {
		neighbor := neighbors[0]
		nT := neighbor.(*Tile)
		moveCoord := Coord{X: nT.X, Y: nT.Y}
		move := ParseMove(head, moveCoord)
		return move
	}
	fmt.Println("\nDead end!")
	return "right"
}

func TileSafe(tile Coord, g GameRequest) bool {
	for i := 0; i < len(g.Board.Snakes); i++ {
		snake := g.Board.Snakes[i]
		if snake.Id != g.You.Id {
			snakeHead := snake.Body[0]
			snakeHead = FlipCoords(snakeHead)
			distance := ManhattanDistance(tile, snakeHead)
			fmt.Printf("\nOpponent snake: %v", snake)
			fmt.Printf("\nDestination: %v", tile)
			fmt.Printf("\nOpponent Head: %v", snakeHead)
			fmt.Printf("\nDistance: %v", distance)
			if distance <= 1 && len(snake.Body) >= len(g.You.Body) {
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

func TrappingMove(w World, g GameRequest, c Coord) bool {
	flippedC := FlipCoords(c)
	availableSpace := FloodFill(w, flippedC)
	fmt.Println("\nCHECKING FOR TRAPS")
	fmt.Println(c)
	fmt.Println(availableSpace)

	return availableSpace < len(g.You.Body)
}

// func IdleMove(g GameRequest) bool, string {

// }

// Try to eat food
func HungryMove(world World, g GameRequest) (bool, string) {
	foundMove := false
	head := g.You.Body[0]
	foodList := ListFoodDistances(head, g.Board.Food)
	foodList = SortByDistance(foodList)
	for i := 0; i < len(foodList); i++ {
		goalCoord := foodList[i].Loc
		world.SetGoal(goalCoord)
		fmt.Println("\n\n\nworld")
		fmt.Printf(StringifyWorld(world))
		p, _, found := Path(world.From(), world.To())
		if found {
			head := Coord{
				X: world.From().X,
				Y: world.From().Y,
			}
			cutPath := p[len(p)-2]
			moveTile := cutPath.(*Tile)
			moveCoord := Coord{X: moveTile.X, Y: moveTile.Y}
			if TileSafe(moveCoord, g) && !TrappingMove(world, g, moveCoord) {
				fmt.Println("SAFE MOVE")
				move := ParseMove(head, moveCoord)
				foundMove = true
				return foundMove, move
			}
		}
		world.StripGoal(goalCoord)
	}
	return foundMove, ""
}

// Chase own tail
func ScaredyMove(world World, g GameRequest) (bool, string) {
	foundMove := false
	tail := g.You.Body[len(g.You.Body)-1]

	// if head and tail are on same tile, default to something else
	// only relevent for first move
	if !world.IsEmpty(tail) {
		return foundMove, ""
	}

	// tail = FlipCoords(tail)
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
		if TileSafe(moveCoord, g) {
			move := ParseMove(head, moveCoord)
			foundMove = true
			return foundMove, move
		}
	}
	return foundMove, ""
}

func FindMove(g GameRequest) string {
	world := ParseWorldFromRequest(g)
	foundMove := false
	move := ""
	health := g.You.Health
	fmt.Println("HEALTH %v", health)
	if health > 30 {
		foundMove, move = ScaredyMove(world, g)
		if foundMove {
			return move
		}
	}
	foundMove, move = HungryMove(world, g)
	if foundMove {
		return move
	}
	return LastResort(world)
}
