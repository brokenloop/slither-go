package main

// func (g GameRequest) MoveSnake(s Snake, m string) {
// 	g.Board.s
// }

func (s Snake) Move(m string) {
	direction := ParseDirection(m)
	oldHead := s.Body[0]
	newHead := Coord{X: oldHead.X + direction[0], Y: oldHead.Y + direction[1]}
	s.Body = append([]Coord{newHead}, s.Body...)
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
 