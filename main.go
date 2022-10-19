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
var returnArray []node

var solveLen int

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
	x    int
	y    int
	move []string
}

func search(nodearr *[]node, diffx int, diffy int) string {

	na_len := len(*nodearr)

	newx := (*nodearr)[na_len-1].x + diffx
	newy := (*nodearr)[na_len-1].y + diffy

	//fmt.Printf("starting from: \n%v\n\tgoing to: %v, %v\n\texcluding: %v\n\tmin path: %v\n", *nodearr, newy, newx, visited, solveLen)

	if newx < 0 || newy < 0 || newy >= len(worldmap) || newx >= len(worldmap[newy]) || (solveLen > 0 && na_len >= solveLen) {
		return "invalid"
	}

	// look if the node as already been visited, in this case ignore it
	_, exists := visited[fmt.Sprintf("%d, %d", newy, newx)]
	if exists {
		return "invalid"
	}

	//fmt.Println("looking at: ", worldmap[newy][newx])
	switch worldmap[newy][newx] {

	case "G":
		//if the goal is hit then add the last element and exit
		*nodearr = append(*nodearr, node{
			x: newx,
			y: newy,
		})
		solveLen = len(*nodearr)
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

		//fmt.Println("\n-------------------------------------------------------\n")
		//printResult(*nodearr)

		dir, err := relPos((*nodearr)[na_len-1], (*nodearr)[na_len])

		if err != nil {
			fmt.Println("error in getting relative position")
		}

		(*nodearr)[na_len-1].move = append((*nodearr)[na_len-1].move, dir)

		// search north, west, east, south of that node

		// go east
		if search(nodearr, 1, 0) == "goal" {
			return "goal"
		}

		// go south
		if search(nodearr, 0, 1) == "goal" {
			return "goal"
		}

		// go west
		if search(nodearr, -1, 0) == "goal" {
			return "goal"
		}

		// go north
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

func refine(sol []node) {

	sol = sol[:len(sol)-1]

	for i := 0; i < len(sol); i++ {
		subNode := sol[:len(sol)-i]
		visited = make(map[string]bool)

		for _, t := range subNode[:len(subNode)-1] {
			visited[fmt.Sprintf("%d, %d", t.y, t.x)] = true
		}

		// go east
		if !Contains(subNode[len(subNode)-1].move, "east") {
			if search(&subNode, 1, 0) == "goal" {
				returnArray = cpyNodeArray(subNode)
				continue
			}
		}

		// go south
		if !Contains(subNode[len(subNode)-1].move, "south") {
			if search(&subNode, 0, 1) == "goal" {
				returnArray = cpyNodeArray(subNode)
				continue
			}
		}

		// go west
		if !Contains(subNode[len(subNode)-1].move, "west") {
			if search(&subNode, -1, 0) == "goal" {
				returnArray = cpyNodeArray(subNode)
				continue
			}
		}

		// go north
		if !Contains(subNode[len(subNode)-1].move, "north") {
			if search(&subNode, 0, -1) == "goal" {
				returnArray = cpyNodeArray(subNode)
				continue
			}
		}
	}

	if len(sol) > 2 {
		refine(cpyNodeArray(sol))
	}

}

func cpyNodeArray(i []node) []node {
	//printResult(i)
	cpy := make([]node, len(i))
	for i, t := range i {
		cpy[i] = node{
			x: t.x,
			y: t.y,
		}
	}
	return cpy
}

func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

func relPos(n1 node, n2 node) (string, error) {

	//fmt.Print(n1, "&", n2, "\t")

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
		//fmt.Print(x1, "<", x2)
		return "east", nil
	}

	if x1 > x2 {
		//fmt.Print(x1, ">", x2)
		return "west", nil
	}

	if y1 > y2 {
		//fmt.Print(y1, ">", y2)
		return "north", nil
	}

	if y1 < y2 {
		//fmt.Print(y1, "<", y2)
		return "south", nil
	}

	return "", fmt.Errorf("no condition has been met")
}

func prettyPrint(inp [][]string) {

	fmt.Print("   ")
	for i := 0; i < len(inp); i++ {
		secdec := i / 10
		if secdec > 0 {
			fmt.Print(secdec)
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()

	fmt.Print("   ")
	for i := 0; i < len(inp); i++ {
		fmt.Print(i % 10)
	}
	fmt.Print("\n\n")

	for i, y := range inp {

		fmt.Print(i)

		if i/10 < 1 {
			fmt.Print(" ")
		}

		fmt.Print(" ")

		for _, x := range y {
			fmt.Print(x)
		}
		fmt.Println()
	}

	fmt.Println()

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

	fmt.Println("\n\n...........................................\n\n")

	//tmpworldmap := worldmap

	tmpworldmap := make([][]string, len(worldmap))

	for i := range worldmap {
		tmpworldmap[i] = make([]string, len(worldmap[i]))
		copy(tmpworldmap[i], worldmap[i])
	}

	/*

		for i, val := range res {
			if i == 0 || i == len(res)-1 {
				continue
			}
			tmpworldmap[val.y][val.x] = "╬"
		}

	*/

	/*
	 */

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
			fmt.Printf("error during execution: %v\n", err)
			//return
		}

		dir2, err := relPos(val, res[i+1])

		if err != nil {
			fmt.Printf("error during execution: %v\n", err)
			//return
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
	/*
	 */

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

	flag.StringVar(&mapDir, "f", "", "specify text file containing map")
	flag.IntVar(&x, "x", 10, "specify labirinth lenght, default is 10")
	flag.IntVar(&y, "y", 10, "specify labirinth height, default is 10")
	flag.Float64Var(&p, "p", 0.8, "specify labirinth percolation, default is 0.8")

	flag.Parse()

	fmt.Println(mapDir)

	if mapDir != "" {
		err = loadMap(mapDir)
		if err != nil {
			fmt.Printf("Error while reading file: %v\n", err)
		}
	} else {
		worldmap = generateMapByPercolation(x, y, p)
	}

	prettyPrint(worldmap)

	res, err := solve()

	if err != nil {
		fmt.Printf("error during execution: %v\n", err)
		return
	}

	printResult(res)

	refine(res)

	printResult(returnArray)
}
