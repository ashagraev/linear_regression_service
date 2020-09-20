package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"github.com/golang/protobuf/jsonpb"

	pb "linear_regression_service/github.com/ashagraev/linear_regression"
)

func newTrainingGRPCClient() *regressionClient {
	return newRegressionClient(trainMode, grpcMode)
}

func newCalculatingGRPCClient() *regressionClient {
	return newRegressionClient(applyMode, grpcMode)
}

func newStatsGRPCClient() *regressionClient {
	return newRegressionClient(statsMode, grpcMode)
}

func reportProtoJSON(m proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EmitDefaults: true,
	}

	simpleJSON, err := marshaler.MarshalToString(m)
	if err != nil {
		return "", fmt.Errorf("error marshalling the report message: %v", err)
	}

	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, []byte(simpleJSON), "", "    "); err != nil {
		return "", fmt.Errorf("error indenting the report message: %v", err)
	}

	return prettyJson.String(), nil
}

func createConnection(serverPath string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	return grpc.Dial(serverPath, opts...)
}

func (rc *regressionClient) requestGRPCTraining(ctx context.Context, pool *pb.Pool) (string, error) {
	conn, err := createConnection(rc.serverPath)
	if err != nil {
		return "", fmt.Errorf("cannot create grpc dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRegressionClient(conn)

	result, err := client.Train(ctx, &pb.TrainingRequest{
		Data:       pool,
		StoreModel: true,
	})
	if err != nil {
		return "", fmt.Errorf("error processing training request: %v", err)
	}

	return reportProtoJSON(result)
}

func (rc *regressionClient) requestGRPCCalculation(ctx context.Context, arg float64) (string, error) {
	conn, err := createConnection(rc.serverPath)
	if err != nil {
		return "", fmt.Errorf("cannot create grpc dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRegressionClient(conn)

	modelValue, err := client.Calculate(ctx, &pb.CalculateRequest{
		Argument: arg,
		ModelName: rc.modelName,
	})
	if err != nil {
		return "", fmt.Errorf("error processing calculation request: %v", err)
	}

	return reportProtoJSON(modelValue)
}

func (rc *regressionClient) requestGRPCStats(ctx context.Context) (string, error) {
	conn, err := createConnection(rc.serverPath)
	if err != nil {
		return "", fmt.Errorf("cannot create grpc dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRegressionClient(conn)
	stats, err := client.Stats(ctx, &pb.StatsRequest{})
	if err != nil {
		return "", fmt.Errorf("error processing stats request: %v", err)
	}

	return reportProtoJSON(stats)
}

func runGRPCTrain() {
	client := newTrainingGRPCClient()

	pool, err := loadProtoPoolFromTSV(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	result, err := client.requestGRPCTraining(ctx, pool)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func runGRPCApply() {
	client := newCalculatingGRPCClient()
	ctx := context.Background()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		arg, err := strconv.ParseFloat(text, 64)
		if err != nil {
			log.Fatalf("invalid float: %v", text)
		}

		result, err := client.requestGRPCCalculation(ctx, arg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	}
}

func runGRPCStats() {
	client := newStatsGRPCClient()
	ctx := context.Background()

	result, err := client.requestGRPCStats(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
