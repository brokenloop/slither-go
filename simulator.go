package main

import (
	"fmt"
	"math/rand"
)

// func (g GameRequest) MoveSnake(s Snake, m string) {
// 	g.Board.s
// }

func Simulate(w World, g GameRequest) {
	for i := 0; i < 10; i++ {
		for i := 0; i < len(g.Board.Snakes); i++ {
			g.Board.Snakes[i].Move("up", false)
		}
		KillSnakes(w, g)
		w = ParseWorldFromRequest(g)
		fmt.Println(StringifyWorld(w))
	}
}

// need to update this to check head on collisions - bigger snake wins
func KillSnakes(w World, g GameRequest) {
	killList := make(map[int]bool)
	for i := 0; i < len(g.Board.Snakes); i++ {

		head := g.Board.Snakes[i].Body[0]
		// if head.X < 0 || head.X >= g.Board.Width || head.Y < 0 || head.Y >= g.Board.Width {
		// 	if w[head.Y][head.X].Kind == KindBlocker {
		// 		killList[i] = true
		// 	}
		// }
		if OutOfBounds(head, g.Board.Width) || w[head.Y][head.X].Kind == KindBlocker {
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
	}
}

func OutOfBounds(head Coord, size int) bool {
	return head.X < 0 || head.X >= size || head.Y < 0 || head.Y >= size
}

func (s *Snake) RandomMove(w World) string {
	head := s.Body[0]
	fmt.Print(w[head.Y][head.X].Kind)
	neighbors := w[head.Y][head.X].PathNeighbors()
	neighbor := neighbors[rand.Intn(len(neighbors))]
	nT := neighbor.(*Tile)
	moveCoord := Coord{X: nT.Y, Y: nT.X}

	fmt.Print("\n")
	fmt.Print(moveCoord)
	move := ParseMove(head, moveCoord)
	return move
}

func DeepCopyWorld(oldWorld World) World {
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
	direction := ParseDirection(m)
	oldHead := s.Body[0]
	newHead := Coord{X: oldHead.X + direction[0], Y: oldHead.Y + direction[1]}
	if eat {
		s.Body = append([]Coord{newHead}, s.Body...)
	} else {
		s.Body = append([]Coord{newHead}, s.Body[:len(s.Body)-1]...)
	}

}

func ParseDirection(m string) [2]int {
	result := [2]int{0, 0}
	switch m {
	case "up":
		result = [2]int{-1, 0}
	case "down":
		result = [2]int{1, 0}
	case "left":
		result = [2]int{0, -1}
	case "right":
		result = [2]int{0, 1}
	}
	return result
}
