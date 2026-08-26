package main

import (
	"cidadon/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		panic(err)
	}
}
