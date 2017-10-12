// +build matr

package main

import (
	"context"

	"github.com/euforic/matr/tlkn"
)

// Proto generates go files from protobuf
func Proto(ctx context.Context) (context.Context, error) {
	return ctx, tlkn.Bash(`protoc -I . -I ${GOPATH}/src/ \
		--proto_path=${GOPATH}/src/github.com/gogo/protobuf:. \
		--gogofast_out=plugins=grpc,\
		Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:. \
		proto/*.proto
	`)
}

// Test runs all go tests
func Test(ctx context.Context) (context.Context, error) {
	return ctx, tlkn.Bash(`go test -v ./...`)
}

// Bench runs all the go benchmarks
func Bench(ctx context.Context) (context.Context, error) {
	return ctx, tlkn.Bash(`go test -bench=. -benchmem -benchtime 10s ./...`)
}

// Docker builds static go binary then builds a docker image for it
func Docker(ctx context.Context) (context.Context, error) {
	return ctx, tlkn.Bash(`
	GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o build/kobs ./kobs/.
	docker build . -t euforic/kobs:latest
	`)
}
