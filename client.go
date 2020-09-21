package main

import (
	"flag"
	"log"
)

type regressionClient struct {
	serverPath string
	modelName string
}

type operationMode int

const (
	calculateMode operationMode = iota
	trainMode
	statsMode
)

func clientMode(operation operationMode, protocol protocolMode) string {
	operationStr := ""
	switch operation {
	case calculateMode: operationStr = "calc"
	case trainMode: operationStr = "train"
	case statsMode: operationStr = "stats"
	}
	return protocolPrefix(protocol) + "-" + operationStr
}

func clientModeArg(operation operationMode, protocol protocolMode) string {
	return "--" + clientMode(operation, protocol)
}

func clientUsage(operation operationMode) string {
	switch operation {
	case calculateMode: return "calculate model value"
	case trainMode: return "train model"
	case statsMode: return "collect service execution stats"
	}
	log.Fatalf("unknown operation mode: %v", operation)
	return ""
}

func newRegressionClient(operation operationMode, protocol protocolMode) *regressionClient {
	flag.Bool(clientMode(operation, protocol), true, clientUsage(operation))
	var server = flag.String("server", "", "network path of the training server")
	var model = flag.String("model", "", "model name for calculation")
	flag.Parse()

	return &regressionClient{serverPath: *server, modelName: *model}
}
