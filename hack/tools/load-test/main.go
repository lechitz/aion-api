package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	cfg := loadConfig(os.Args[1:])
	if err := run(context.Background(), cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load test failed: %v\n", err)
		os.Exit(1)
	}
}
