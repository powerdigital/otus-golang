package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

func main() {
	var reverted = stringutil.Reverse("Hello, OTUS!")
	fmt.Println(reverted)
}
