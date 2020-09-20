package main

import (
	"context"
	"errors"
	"flag"
	"log"
)

type protocolMode int

const (
	httpMode protocolMode = iota
	grpcMode
)

func protocolPrefix(mode protocolMode) string {
	switch mode {
	case httpMode: return "http"
	case grpcMode: return "grpc"
	}
	log.Fatalf("unknown protocol mode: %v", mode)
	return ""
}

func handlerMode(mode protocolMode) string {
	return protocolPrefix(mode) + "-server"
}

func handlerModeArg(mode protocolMode) string {
	return "--" + handlerMode(mode)
}

func handlerContext(mode protocolMode) (context.Context, error) {
	flag.Bool(handlerMode(mode), true, "run the regression service")

	project := flag.String("spanner-project", "", "Spanner project name")
	instance := flag.String("spanner-instance", "", "Spanner instance name")
	database := flag.String("spanner-database", "", "Spanner database name")

	var port, address *string
	if mode == httpMode {
		port = flag.String("port", "8080", "run the http handler using this port")
	}
	if mode == grpcMode {
		address = flag.String("address", "localhost:8081", "run the grpc handler using this address")
	}
	flag.Parse()

	if len(*project) == 0 {
		return nil, errors.New("choose the spanner project (--spanner-project)")
	}
	if len(*instance) == 0 {
		return nil, errors.New("choose the spanner instance (--spanner-instance)")
	}
	if len(*database) == 0 {
		return nil, errors.New("choose the spanner database (--spanner-database)")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "project", *project)
	ctx = context.WithValue(ctx, "instance", *instance)
	ctx = context.WithValue(ctx, "database", *database)

	if port != nil {
		ctx = context.WithValue(ctx, "port", *port)
	} else if mode == httpMode {
		return nil, errors.New("choose the port for the http daemon (--port)")
	}

	if address != nil {
		ctx = context.WithValue(ctx, "address", *address)
	} else if mode == grpcMode {
		return nil, errors.New("choose the port for the http daemon (--address)")
	}

	return ctx, nil
}
