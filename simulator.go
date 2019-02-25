package main

// func (g GameRequest) MoveSnake(s Snake, m string) {

// }

func parseDirection(m string) [2]int {
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
