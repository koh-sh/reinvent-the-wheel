//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var sc = bufio.NewScanner(os.Stdin)

func scanI() int {
	sc.Scan()
	a, _ := strconv.Atoi(sc.Text())
	return a
}

func main() {
	var n int
	sc.Split(bufio.ScanWords)
	n = scanI()
	h := make([]int, n)
	highest := 0
	highest_i := 0
	for i := 0; i < n; i++ {
		h[i] = scanI()
	}
	for i := 0; i < n; i++ {
		if h[i] > highest {
			highest_i = i
			highest = h[i]
		}
	}
	fmt.Printf("%d\n", highest_i+1)
}
