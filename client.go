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
	applyMode operationMode = iota
	trainMode
)

func clientMode(operation operationMode, protocol protocolMode) string {
	operationStr := ""
	switch operation {
	case applyMode: operationStr = "apply"
	case trainMode: operationStr = "train"
	}
	return protocolPrefix(protocol) + "-" + operationStr
}

func clientModeArg(operation operationMode, protocol protocolMode) string {
	return "--" + clientMode(operation, protocol)
}

func clientUsage(operation operationMode) string {
	switch operation {
	case applyMode: return "apply model"
	case trainMode: return "train model"
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
