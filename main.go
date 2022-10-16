package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var E = " "
var F = "#"

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

	if newx < 0 || newy < 0 || newy >= len(worldmap) || newx >= len(worldmap[newy]) {
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

	return nil, errors.New("no path found")
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
		return "north", nil
	}

	if y1 < y2 {
		return "south", nil
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

func fillNodeConn(slc *string, p float64) {
	rand.Seed(time.Now().UnixNano())
	if rand.Float64() > p {
		*slc = F
	} else {
		*slc = E
	}
}

func generateMapByPercolation(x int, y int, p float64) [][]string {
	game_map := make([][]string, y*2+1, y*2+1)

	for i, _ := range game_map {
		game_map[i] = make([]string, x*2+1, x*2+1)
	}

	for i := 0; i < y*2+1; i++ {
		for t := 0; t < x*2+1; t++ {
			// print edges
			if i == 0 || i == y*2 || t == 0 || t == x*2 {
				game_map[i][t] = F
			} else if i%2 == 1 {
				if t%2 == 1 {
					game_map[i][t] = F
				} else {
					fillNodeConn(&game_map[i][t], p)
				}
			} else {
				if t%2 == 1 {
					fillNodeConn(&game_map[i][t], p)
				} else {
					//print node
					game_map[i][t] = E
				}
			}

		}
	}

	for {
		gy := rand.Int() % len(game_map)
		gx := rand.Int() % len(game_map[0])

		if game_map[gy][gx] == " " {
			game_map[gy][gx] = "G"
			break
		}
	}
	for {
		sy := rand.Int() % len(game_map)
		sx := rand.Int() % len(game_map[0])

		if game_map[sy][sx] == " " {
			game_map[sy][sx] = "S"
			break
		}
	}

	return game_map
}

func printResult(res []node) {

	tmpworldmap := worldmap

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

		if (dir1 == "north" && dir2 == "east") || (dir2 == "north" && dir1 == "east") {
			tmpworldmap[val.y][val.x] = nte
		} else if (dir1 == "south" && dir2 == "east") || (dir2 == "south" && dir1 == "east") {
			tmpworldmap[val.y][val.x] = ste
		} else if (dir1 == "north" && dir2 == "east") || (dir2 == "north" && dir1 == "east") {
			tmpworldmap[val.y][val.x] = nte
		} else if (dir1 == "west" && dir2 == "south") || (dir2 == "west" && dir1 == "south") {
			tmpworldmap[val.y][val.x] = wts
		} else if (dir1 == "west" && dir2 == "north") || (dir2 == "west" && dir1 == "north") {
			tmpworldmap[val.y][val.x] = wtn
		} else if (dir1 == "north" && dir2 == "south") || (dir2 == "north" && dir1 == "south") {
			tmpworldmap[val.y][val.x] = nts
		} else if (dir1 == "west" && dir2 == "east") || (dir2 == "west" && dir1 == "east") {
			tmpworldmap[val.y][val.x] = wte
		}
	}

	prettyPrint(tmpworldmap)
}

func main() {

	var (
		err    error
		x      int
		y      int
		p      float64
		mapDir string
	)

	fmt.Println("app is now running")

	flag.StringVar(&mapDir, "dir", "", "specify text file containing map")
	flag.IntVar(&x, "x", 10, "specify labirinth lenght, default is 10")
	flag.IntVar(&y, "y", 10, "specify labirinth height, default is 10")
	flag.Float64Var(&p, "p", 0.8, "specify labirinth percolation, default is 0.8")

	flag.Parse()

	if mapDir != "" {
		err = loadMap("map.txt")
		if err != nil {
			fmt.Printf("Error while reading file: %v\n", err)
		}
		return
	} else {
		worldmap = generateMapByPercolation(x, y, p)
	}

	prettyPrint(worldmap)

	fmt.Println("\n-------------------------------------------------------\n")

	res, err := solve()

	if err != nil {
		fmt.Println("error during execution: ", err)
		return
	}

	printResult(res)
}
