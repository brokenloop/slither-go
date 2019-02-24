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
	fmt.Printf("Shortest distance: %d", shortestDistance)
	fmt.Printf("Closest food: %d", result)
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
	fmt.Println("\n\n\nLAST RESORT")
	// headTile := w[head.Y][head.X]
	headTile := w.From()
	head := Coord{
		X: headTile.X,
		Y: headTile.Y,
	}
	neighbors := headTile.PathNeighbors()
	if len(neighbors) > 0 {
		neighbor := neighbors[0]
		nT := neighbor.(*Tile)
		// coords are swapped because of world mapping
		moveCoord := Coord{X: nT.X, Y: nT.Y}
		fmt.Printf("Head coord: %d\n", head)
		fmt.Printf("Move coord: %d\n", moveCoord)

		move := ParseMove(head, moveCoord)
		return move
	}
	fmt.Println("\nDead end!")
	return "right"
}

func FindMove(g GameRequest) string {
	world := ParseWorldFromRequest(g)
	head := g.You.Body[0]
	foodList := ListFoodDistances(head, g.Board.Food)
	foodList = SortByDistance(foodList)
	for i := 0; i < len(foodList); i++ {
		goalCoord := foodList[i].Loc
		world.SetGoal(goalCoord)
		fmt.Println("\nworld")
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
			move := ParseMove(head, moveCoord)
			return move
		}
	}
	return LastResort(world)
}
