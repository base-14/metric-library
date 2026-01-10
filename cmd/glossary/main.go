package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("OTel Glossary - Metric Discovery Platform")
	fmt.Println("Phase 0 complete. Ready for Phase 1 implementation.")
	return nil
}
