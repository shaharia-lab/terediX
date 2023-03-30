// Main package
package main

import (
	"context"
	"log"

	"teredix/pkg/cmd"
)

func main() {

	ctx := context.Background()

	if err := cmd.NewRootCmd("1.0").ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
