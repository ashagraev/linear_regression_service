/*
linear_regression_service is a demonstration program which runs http and gRPC services
providing API for training simple regression models.
 */
package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("choose mode: server, train or calc")
	}

	if os.Args[1] == handlerModeArg(httpMode) {
		runHTTPHandler()
	}
	if os.Args[1] == clientModeArg(trainMode, httpMode) {
		runHTTPTraining()
	}
	if os.Args[1] == clientModeArg(calculateMode, httpMode) {
		runHTTPCalculation()
	}
	if os.Args[1] == clientModeArg(statsMode, httpMode) {
		runHTTPStats()
	}

	if os.Args[1] == handlerModeArg(grpcMode) {
		runGRPCHandler()
	}
	if os.Args[1] == clientModeArg(trainMode, grpcMode) {
		runGRPCTraining()
	}
	if os.Args[1] == clientModeArg(calculateMode, grpcMode) {
		runGRPCCalculation()
	}
	if os.Args[1] == clientModeArg(statsMode, grpcMode) {
		runGRPCStats()
	}
}
