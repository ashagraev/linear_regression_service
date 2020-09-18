package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "regression_service/github.com/ashagraev/linear_regression"
	"time"
)

type grpcHandler struct {
	stats        pb.ServerStats
	requestStats chan pb.ServerStats

	modelsStorage *ModelsStorage
}

func newGRPCHandler(ctx context.Context) (*grpcHandler, error) {
	modelsStorage, err := NewModelsStorage(ctx)
	if err != nil {
		return nil, err
	}
	h := grpcHandler{requestStats: make(chan pb.ServerStats), modelsStorage: modelsStorage}
	go h.updateStatsLoop()
	return &h, nil
}

func (h *grpcHandler) Svc() *pb.RegressionSolverService {
	return &pb.RegressionSolverService{
		Solve: h.Solve,
		Calculate: h.Calculate,
	}
}

func (h *grpcHandler) updateStatsLoop() {
	for r := range h.requestStats {
		h.stats.TotalRequests += r.TotalRequests
		h.stats.TotalInstances += r.TotalInstances
		h.stats.SucceededRequests += r.SucceededRequests
	}
}

func (h *grpcHandler) Solve(ctx context.Context, request *pb.SolveRequest) (*pb.Solution, error) {
	var slr SimpleLinearRegression
	for _, instance := range request.Data.Instances {
		slr.AddWeightedInstance(instance.Argument, instance.Target, instance.Weight)
	}

	model := slr.Solve()
	solution := pb.Solution{
		Model: &pb.SimpleRegressionModel{
			Coefficient: model.Coefficient,
			Intercept: model.Intercept,
		},
		SumSquaredErrors: slr.SumSquaredErrors(),
	}

	if request.StoreModel {
		name, commitTime, err := h.modelsStorage.SaveSLRModel(ctx, model)
		if err != nil {
			solution.Error = fmt.Sprintf("%v", err)
		}
		solution.Name = name
		solution.CreationTime = fmt.Sprintf("%v", commitTime)
	}

	return &solution, nil
}

func (h *grpcHandler) Calculate(ctx context.Context, request *pb.CalculateRequest) (*pb.ModelValue, error) {
	requestInfo := pb.ServerStats{
		TotalRequests: 1,
	}
	defer func() {
		h.requestStats <- requestInfo
	}()

	modelValue := pb.ModelValue{}

	model, fromCache, err := h.modelsStorage.GetSLRModel(ctx, request.ModelName)
	if err != nil {
		modelValue.Error = fmt.Sprintf("error loading model %v: %v", request.ModelName, err)
		return &modelValue, err
	}

	modelValue.Argument = request.Argument
	modelValue.FromCache = fromCache
	modelValue.Value = model.Apply(request.Argument)
	modelValue.CalculationTime = fmt.Sprintf("%v", time.Now())
	modelValue.Model = &pb.SimpleRegressionModel{
		Name:        model.Name,
		Intercept:   model.Intercept,
		Coefficient: model.Coefficient,
	}
	requestInfo.SucceededRequests = 1

	return &modelValue, nil
}

func runGRPCHandler() {
	flag.Bool("grpc-server", true, "run the solving server")
	address := flag.String("address", "localhost:80", "grpc handler network address")
	flag.Parse()

	ctx := context.Background()
	h, err := newGRPCHandler(ctx)
	if err != nil {
		log.Fatal("cannot create handler: ", err)
	}

	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRegressionSolverService(grpcServer, h.Svc())
	grpcServer.Serve(lis)
}
