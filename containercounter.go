//go:build ignore

// https://github.com/koh-sh/iterm2-container-counter
package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	whale := "🐳 "
	stopped := "⚫ "
	switch_view := 10

	cmd := exec.Command("/usr/local/bin/docker", "container", "ls", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	outlines := strings.Split(out.String(), "\n")
	containers := outlines[1 : len(outlines)-1]
	all_num := len(containers)
	var run_num int
	for _, v := range containers {
		if strings.Contains(v, "Up ") {
			run_num++
		}
	}

	if all_num < switch_view {
		for i := 0; i < all_num; i++ {
			if i < run_num {
				fmt.Print(whale)
			} else {
				fmt.Print(stopped)
			}
		}
		fmt.Printf("\n")
	} else {
		stp_num := all_num - run_num
		fmt.Printf("%s * %d | %s * %d\n", whale, run_num, stopped, stp_num)
	}
}
