package main

import (
	"events-exercise/config"
	"events-exercise/internal/app"
	"path/filepath"

	"github.com/joho/godotenv"
)

// In the original task the candidate is asked to implement their solution in the main.go

// I would suggest providing some kind of incomplete implementation of the reader here but point to
// `internal/events/processor.go` for the implementation of the processor.

// Ultimately if the candidate goes on to complete the http service there shouldn't be much here...
func main() {
	_ = godotenv.Load(filepath.Join("..", "..", ".env"))

	cfg := config.LoadConfig()
	a := app.Load(cfg)
	a.Run()
}
