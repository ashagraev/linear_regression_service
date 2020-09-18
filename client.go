package main

import "flag"

type regressionClient struct {
	serverPath string
	modelName string
}

func newRegressionClient(mode string, protocol string, usage string) *regressionClient {
	flag.Bool(protocol + "-" + mode, true, usage)
	var server = flag.String("server", "", "network path of the training server")
	var model = flag.String("model", "", "model name for calculation")
	flag.Parse()

	return &regressionClient{serverPath: *server, modelName: *model}
}
