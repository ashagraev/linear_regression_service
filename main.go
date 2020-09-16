/*
linear_regression_service is a demonstration program which runs an http service
providing API for solving simple regression tasks
 */
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var port = flag.String("port", "80", "http service port")
	flag.Parse()

	var h handler
	h.init()

	http.Handle("/solve", http.HandlerFunc(h.handleSolveRequest))
	http.Handle("/stats", http.HandlerFunc(h.handleStatsRequest))

	go h.updateStatsLoop()

	err := http.ListenAndServe(":" + *port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
