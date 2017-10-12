// +build matr

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/euforic/matr"
)

func main() {
	// Create new Matr instance
	m := matr.New()

	// Override the default Handler
	m.Handle("_default_", func(ctx context.Context) (context.Context, error) {
		fmt.Print("\nAvailable Commands:\n\n",
			"	proto   Proto generates go files from protobuf\n",
			"	test   Test runs all go tests\n",
			"	bench   Bench runs all the go benchmarks\n",
			"	docker   Docker builds static go binary then builds a docker image for it\n","\n",
		)
		return ctx, nil
	})
	
	// Proto generates go files from protobuf
	m.Handle("proto", Proto)
	
	// Test runs all go tests
	m.Handle("test", Test)
	
	// Bench runs all the go benchmarks
	m.Handle("bench", Bench)
	
	// Docker builds static go binary then builds a docker image for it
	m.Handle("docker", Docker)

	// Run Matr
	if err := m.Run(context.Background(), os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
