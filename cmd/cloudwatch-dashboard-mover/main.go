package main

import (
	"context"
	"log"
	"os"

	mover "github.com/handlename/cloudwatch-dashboard-mover"
)

func main() {
	ctx := context.Background()

	input := mover.ParseMoverInput()
	m := mover.NewMover(input)
	out, err := m.Replace(ctx)
	if err != nil {
		log.Fatalf("failed to replace: %v", err)
	}

	os.Stdout.Write(out)
}
