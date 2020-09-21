package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "linear_regression_service/github.com/ashagraev/linear_regression"
	"log"
	"net"
	"sync"
	"time"
)

type grpcHandler struct {
	stats        pb.ServerStats
	requestStats chan pb.ServerStats

	statsMutex sync.Mutex

	modelsStorage *modelsStorage
}

func newGRPCHandler(ctx context.Context) (*grpcHandler, error) {
	modelsStorage, err := newModelsStorage(ctx)
	if err != nil {
		return nil, err
	}
	h := grpcHandler{requestStats: make(chan pb.ServerStats), modelsStorage: modelsStorage}
	go h.updateStatsLoop()
	return &h, nil
}

func (h* grpcHandler) getStats() pb.ServerStats {
	h.statsMutex.Lock()
	defer h.statsMutex.Unlock()

	return h.stats
}

func (h* grpcHandler) setStats(stats pb.ServerStats) {
	h.statsMutex.Lock()
	defer h.statsMutex.Unlock()

	h.stats = stats
}

func (h *grpcHandler) Svc() *pb.RegressionService {
	return &pb.RegressionService{
		Train: h.Train,
		Calculate: h.Calculate,
		Stats: h.Stats,
	}
}

func (h *grpcHandler) updateStatsLoop() {
	for r := range h.requestStats {
		stats := h.getStats()
		stats.TotalRequests += r.TotalRequests
		stats.TotalInstances += r.TotalInstances
		stats.SucceededRequests += r.SucceededRequests
		h.setStats(stats)
	}
}

func (h *grpcHandler) Train(ctx context.Context, request *pb.TrainingRequest) (*pb.TrainingResults, error) {
	var slr SimpleLinearRegression
	for _, instance := range request.Instances {
		slr.AddWeightedInstance(instance.Argument, instance.Target, instance.Weight)
	}

	model := slr.Train()
	result := pb.TrainingResults{
		Model: &pb.SimpleRegressionModel{
			Coefficient: model.Coefficient,
			Intercept: model.Intercept,
		},
		SumSquaredErrors: slr.SumSquaredErrors(),
	}

	if request.StoreModel {
		name, commitTime, err := h.modelsStorage.saveSLRModel(ctx, model)
		if err != nil {
			result.Error = fmt.Sprintf("%v", err)
		}
		result.Name = name
		result.CreationTime = fmt.Sprintf("%v", commitTime)
	}

	return &result, nil
}

func (h *grpcHandler) Calculate(ctx context.Context, request *pb.CalculateRequest) (*pb.ModelValue, error) {
	requestInfo := pb.ServerStats{
		TotalRequests: 1,
	}
	defer func() {
		h.requestStats <- requestInfo
	}()

	modelValue := pb.ModelValue{}

	model, fromCache, err := h.modelsStorage.getSLRModel(ctx, request.ModelName)
	if err != nil {
		modelValue.Error = fmt.Sprintf("error loading model %v: %v", request.ModelName, err)
		return &modelValue, err
	}

	modelValue.Argument = request.Argument
	modelValue.FromCache = fromCache
	modelValue.Value = model.Calculate(request.Argument)
	modelValue.CalculationTime = fmt.Sprintf("%v", time.Now())
	modelValue.Model = &pb.SimpleRegressionModel{
		Name:        model.Name,
		Intercept:   model.Intercept,
		Coefficient: model.Coefficient,
	}
	requestInfo.SucceededRequests = 1

	return &modelValue, nil
}

func (h *grpcHandler) Stats(_ context.Context, _ *pb.StatsRequest) (*pb.ServerStats, error) {
	stats := h.getStats()
	return &stats, nil
}

func runGRPCHandler() {
	ctx, err := handlerContext(grpcMode)
	if err != nil {
		log.Fatal("cannot create context: ", err)
	}

	h, err := newGRPCHandler(ctx)
	if err != nil {
		log.Fatal("cannot create handler: ", err)
	}

	address := ctx.Value("address")
	lis, err := net.Listen("tcp", address.(string))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRegressionService(grpcServer, h.Svc())
	grpcServer.Serve(lis)
}
