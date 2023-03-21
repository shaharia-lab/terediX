package main

import (
	"context"
	"infrastructure-discovery/pkg/cmd"
	"log"
)

func main() {

	ctx := context.Background()

	if err := cmd.NewRootCmd("1.0").ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
