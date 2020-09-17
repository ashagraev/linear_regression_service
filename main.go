/*
linear_regression_service is a demonstration program which runs an http service
providing API for solving simple regression tasks
 */
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	var port = flag.String("port", "80", "http service port")
	flag.Parse()

	ctx := context.Background()
	h, err := newHandler(ctx)
	if err != nil {
		log.Fatal("cannot create handler: ", err)
	}

	http.Handle("/solve", http.HandlerFunc(h.handleSolveRequest))
	http.Handle("/apply", http.HandlerFunc(h.handleApplyRequest))
	http.Handle("/stats", http.HandlerFunc(h.handleStatsRequest))

	go h.updateStatsLoop()

	err = http.ListenAndServe(":" + *port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
