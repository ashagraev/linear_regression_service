/*
linear_regression_service is a demonstration program which runs an http service
providing API for solving simple regression tasks
 */
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func runHandler() {
	flag.Bool("server", true, "run the solving server")
	port := flag.String("port", "80", "run http handler using this port")
	flag.Parse()

	ctx := context.Background()
	h, err := newHandler(ctx)
	if err != nil {
		log.Fatal("cannot create handler: ", err)
	}

	http.Handle("/solve", http.HandlerFunc(h.handleSolveRequest))
	http.Handle("/apply", http.HandlerFunc(h.handleApplyRequest))
	http.Handle("/stats", http.HandlerFunc(h.handleStatsRequest))

	err = http.ListenAndServe(":" + *port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func runTrain() {
	client := NewTrainingClient()

	instances, err := loadInstancesFromTSV(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.requestTraining(instances)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func runApply() {
	client := NewCalculatingClient()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		arg, err := strconv.ParseFloat(text, 64)
		if err != nil {
			log.Fatalf("invalid float: %v", text)
		}

		result, err := client.requestCalculation(arg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("choose mode: server, train or apply")
	}

	if os.Args[1] == "--server" {
		runHandler()
	}
	if os.Args[1] == "--train" {
		runTrain()
	}
	if os.Args[1] == "--apply" {
		runApply()
	}
}
