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

	pb "regression_service/github.com/ashagraev/linear_regression"
)

func NewTrainingGRPCClient() *regressionClient {
	return NewRegressionClient("train","grpc","train model")
}

func NewCalculatingGRPCClient() *regressionClient {
	return NewRegressionClient("apply","grpc","apply model")
}

func reportProtoJSON(m proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{}

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

func (rc *regressionClient) requestGRPCTraining(ctx context.Context, pool *pb.Pool) (string, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(rc.serverPath, opts...)
	if err != nil {
		return "", fmt.Errorf("cannot create grpc dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRegressionSolverClient(conn)

	solution, err := client.Solve(ctx, &pb.SolveRequest{
		Data:       pool,
		StoreModel: true,
	})
	if err != nil {
		return "", fmt.Errorf("error processing training request: %v", err)
	}

	return reportProtoJSON(solution)
}

func (rc *regressionClient) requestGRPCCalculation(ctx context.Context, arg float64) (string, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(rc.serverPath, opts...)
	if err != nil {
		return "", fmt.Errorf("cannot create grpc dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRegressionSolverClient(conn)

	modelValue, err := client.Calculate(ctx, &pb.CalculateRequest{
		Argument: arg,
		ModelName: rc.modelName,
	})
	if err != nil {
		return "", fmt.Errorf("error processing calculation request: %v", err)
	}

	return reportProtoJSON(modelValue)
}

func runGRPCTrain() {
	client := NewTrainingGRPCClient()

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
	client := NewCalculatingGRPCClient()
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
