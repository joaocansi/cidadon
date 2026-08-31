package main

import (
	"cidadon/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		panic(err)
	}
}
