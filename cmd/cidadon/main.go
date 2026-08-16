package main

import (
	"cidadon/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
