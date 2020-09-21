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

	if os.Args[1] == handlerModeArg(httpMode) {
		runHTTPHandler()
	}
	if os.Args[1] == clientModeArg(trainMode, httpMode) {
		runHTTPTrain()
	}
	if os.Args[1] == clientModeArg(applyMode, httpMode) {
		runHTTPApply()
	}
	if os.Args[1] == clientModeArg(statsMode, httpMode) {
		runHTTPStats()
	}

	if os.Args[1] == handlerModeArg(grpcMode) {
		runGRPCHandler()
	}
	if os.Args[1] == clientModeArg(trainMode, grpcMode) {
		runGRPCTrain()
	}
	if os.Args[1] == clientModeArg(applyMode, grpcMode) {
		runGRPCApply()
	}
	if os.Args[1] == clientModeArg(statsMode, grpcMode) {
		runGRPCStats()
	}
}
