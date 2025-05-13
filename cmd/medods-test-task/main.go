package main

import (
	"context"

	"medods-test-task/internal/server"
)

func main() {
	server.Run(context.Background())
}
