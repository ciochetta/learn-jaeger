package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"strconv"

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

	client, err := rpc.Dial("tcp", ":1234")

	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	defer client.Close()

	var result int
	err = client.Call("Service.Square", n, &result)

	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	fmt.Fprintf(w, "Square of %d is %d", n, result)
}
