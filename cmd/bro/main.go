package main

import (
	"crypto/rand"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.ToLower(rand.Text()))
}
