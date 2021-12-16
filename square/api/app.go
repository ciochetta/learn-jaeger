package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"

	sq "github.com/ciochetta/go-square/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func run() error {

	helloHandler := otelhttp.NewHandler(http.HandlerFunc(HandleHello), "Hello")
	http.Handle("/hello", helloHandler)
	squareHandler := otelhttp.NewHandler(http.HandlerFunc(HandleSquare), "Square")
	http.Handle("/square", squareHandler)
	log.Print("Listening on port 8080")

	err := http.ListenAndServe(":8080", nil)

	return err

}

func HandleHello(w http.ResponseWriter, r *http.Request) {
	var name string = r.URL.Query().Get("name")

	fmt.Fprintf(w, "Hello, %s!", name)
}

func HandleSquare(w http.ResponseWriter, r *http.Request) {

	var n int

	n, err := strconv.Atoi(r.URL.Query().Get("n"))

	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))

	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	defer conn.Close()

	c := sq.NewSquareClient(conn)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)

	defer cancel()

	res, err := c.GetSquare(ctx, &sq.GetSquareRequest{Number: int32(n)})

	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	fmt.Fprintf(w, "Square of %d is %d", n, res.GetNumber())

}
