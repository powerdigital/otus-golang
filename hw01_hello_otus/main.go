package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	reverted := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(reverted)
}
