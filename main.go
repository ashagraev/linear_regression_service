/*
linear_regression_service is a demonstration program which runs http and gRPC services
providing API for training simple regression models
 */
package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("choose mode: server, train or apply")
	}

	if os.Args[1] == "--http-server" {
		runHTTPHandler()
	}
	if os.Args[1] == "--http-train" {
		runHTTPTrain()
	}
	if os.Args[1] == "--http-apply" {
		runHTTPApply()
	}

	if os.Args[1] == "--grpc-server" {
		runGRPCHandler()
	}
	if os.Args[1] == "--grpc-train" {
		runGRPCTrain()
	}
	if os.Args[1] == "--grpc-apply" {
		runGRPCApply()
	}
}
