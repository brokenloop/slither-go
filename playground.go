package main

import (
	"fmt"
)

func testPath(worldInput string, expectedDist float64) {
	world := ParseWorld(worldInput)

	fmt.Println("Input world\n", world.RenderPath([]Pather{}))
	p, dist, found := Path(world.From(), world.To())
	if !found {
		fmt.Println("Could not find a path")
	} else {
		fmt.Println("Resulting path\n", world.RenderPath(p))
	}
	if !found && expectedDist >= 0 {
		fmt.Println("Could not find a path")
	}
	if found && dist != expectedDist {
		fmt.Println("Expected dist to be %v but got %v", expectedDist, dist)
	}
}

// path.
func TestStraightLine() {
	testPath(`
.....~......
.....MM.....
.F........T.
....MMM.....
............
`, 9)
}
