package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var worldmap [][]string
var visited = make(map[string]bool)

func loadMap(filename string) error {

	readFile, err := os.Open(filename)

	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		worldmap = append(worldmap, strings.Split(fileScanner.Text(), ""))
	}

	readFile.Close()

	return nil
}

type node struct {
	x int
	y int
}

func search(nodearr *[]node, diffx int, diffy int) string {

	na_len := len(*nodearr)

	newx := (*nodearr)[na_len-1].x + diffx
	newy := (*nodearr)[na_len-1].y + diffy

	/*
		fmt.Printf("looking at: %d, %d", newx, newy)
		bin := ""
		fmt.Scan(&bin)
	*/

	if newy >= len(worldmap) || newx >= len(worldmap[newy]) || newx < 0 || newy < 0 {
		return "invalid"
	}

	// look if the node as already been visited, in this case ignore it
	_, exists := visited[fmt.Sprintf("%d, %d", newy, newx)]
	if exists {
		return "invalid"
	}

	switch worldmap[newy][newx] {

	case "G":
		//if the goal is hit then add the last element and exit
		*nodearr = append(*nodearr, node{
			x: newx,
			y: newy,
		})
		return "goal"
	case "#":
		return "invalid"
	case " ":
		// add the spot to the node array

		*nodearr = append(*nodearr, node{
			x: newx,
			y: newy,
		})

		visited[fmt.Sprintf("%d, %d", newy, newx)] = true

		// search north, west, east, south of that node

		if search(nodearr, 1, 0) == "goal" {
			return "goal"
		}

		if search(nodearr, 0, 1) == "goal" {
			return "goal"
		}

		if search(nodearr, -1, 0) == "goal" {
			return "goal"
		}

		if search(nodearr, 0, -1) == "goal" {
			return "goal"
		}

		// in case the node returns all walls then remove the node

		*nodearr = (*nodearr)[:na_len]

		return "invalid"
	}

	return "invalid"
}

func solve() ([]node, error) {

	var startx, starty int
	var nodeArr []node

	//identify start coordinates

	for i, y := range worldmap {
		for t, x := range y {
			if x == "S" {
				startx = t
				starty = i
			}
		}

	}

	fmt.Println(startx+1, starty+1)

	// add start as first node

	nodeArr = append(nodeArr, node{
		x: startx,
		y: starty,
	})

	visited[fmt.Sprintf("%d, %d", starty, startx)] = true
	// recursively search starting from north of the node

	if search(&nodeArr, 1, 0) == "goal" {
		return nodeArr, nil
	}

	if search(&nodeArr, 0, 1) == "goal" {
		return nodeArr, nil
	}

	if search(&nodeArr, -1, 0) == "goal" {
		return nodeArr, nil
	}

	if search(&nodeArr, 0, -1) == "goal" {
		return nodeArr, nil
	}

	return nil, fmt.Errorf("no path found")
}

func relPos(n1 node, n2 node) (string, error) {

	x1 := n1.x
	y1 := n1.y
	x2 := n2.x
	y2 := n2.y

	if x1 != x2 && y1 != y2 {
		return "", fmt.Errorf("error finding relative position, one axis must be equal")
	}

	if x1 == x2 && y1 == y2 {
		return "", fmt.Errorf("error finding relative position, only one axis must be equal")
	}

	if x1 < x2 {
		return "east", nil
	}

	if x1 > x2 {
		return "west", nil
	}

	if y1 > y2 {
		return "south", nil
	}

	if y1 < y2 {
		return "north", nil
	}

	return "", fmt.Errorf("no condition has been met")
}

func prettyPrint(inp [][]string) {
	for _, y := range inp {
		for _, x := range y {
			fmt.Print(x)
		}
		fmt.Println()
	}
}

func main() {

	var err error

	err = loadMap("map.txt")

	if err != nil {
		fmt.Printf("Error while reading file: %v\n", err)
	}

	res, err := solve()

	if err != nil {
		fmt.Println("error during execution: ", err.Error)
		return
	}

	fmt.Printf("solved:\n%v", res)

	ste := "╔"

	nte := "╚"

	wts := "╗"

	wtn := "╝"

	nts := "║"

	wte := "═"

	for i, val := range res {
		if i == 0 || i == len(res)-1 {
			continue
		}

		dir1, err := relPos(val, res[i-1])

		if err != nil {
			fmt.Println("error during execution: ", err.Error)
			return

		}

		dir2, err := relPos(val, res[i+1])

		if err != nil {
			fmt.Println("error during execution: ", err.Error)
			return
		}

		fmt.Printf("looking at %v: \ncoming from %s\ngoing to %s\n", val, dir1, dir2)

		if (dir1 == "south" && dir2 == "east") || (dir2 == "south" && dir1 == "east") {
			worldmap[val.y][val.x] = ste
		} else if (dir1 == "nord" && dir2 == "east") || (dir2 == "nord" && dir1 == "east") {
			worldmap[val.y][val.x] = nte
		} else if (dir1 == "west" && dir2 == "south") || (dir2 == "west" && dir1 == "south") {
			worldmap[val.y][val.x] = wts
		} else if (dir1 == "west" && dir2 == "north") || (dir2 == "west" && dir1 == "north") {
			worldmap[val.y][val.x] = wtn
		} else if (dir1 == "north" && dir2 == "south") || (dir2 == "north" && dir1 == "south") {
			worldmap[val.y][val.x] = nts
		} else if (dir1 == "west" && dir2 == "east") || (dir2 == "west" && dir1 == "east") {
			worldmap[val.y][val.x] = wte
		}
	}

	prettyPrint(worldmap)
}
