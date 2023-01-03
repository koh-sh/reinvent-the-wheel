//go:build ignore

// print date which is prime number within the year

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	year, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	startdate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)

	for i := 0; i < 366; i++ {
		date := startdate.AddDate(0, 0, i)
		d := date.Format("20060102")
		if d[:4] != strconv.Itoa(year) {
			break
		}
		n, err := strconv.Atoi(d)
		if err != nil {
			log.Fatal(err)
		}
		if isPrime(n) {
			fmt.Printf("%d\n", n)
		}

	}
}
