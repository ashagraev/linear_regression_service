package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// ExecutionStats stores all-time execution statistics for the service.
type ExecutionStats struct {
	// SucceededRequests stores the number of successfully processed requests.
	SucceededRequests int

	// TotalRequests stores the total number of received requests.
	TotalRequests  int

	// TotalInstances stores the total number of instances used while learning models.
	TotalInstances int
}

type httpHandler struct {
	stats        ExecutionStats
	requestStats chan ExecutionStats

	statsMutex sync.Mutex

	modelsStorage *modelsStorage
}

func newHTTPHandler(ctx context.Context) (*httpHandler, error) {
	modelsStorage, err := newModelsStorage(ctx)
	if err != nil {
		return nil, err
	}
	h := httpHandler{requestStats: make(chan ExecutionStats), modelsStorage: modelsStorage}
	go h.updateStatsLoop()
	return &h, nil
}

func (h* httpHandler) getStats() ExecutionStats {
	h.statsMutex.Lock()
	defer h.statsMutex.Unlock()

	return h.stats
}

func (h* httpHandler) setStats(stats ExecutionStats) {
	h.statsMutex.Lock()
	defer h.statsMutex.Unlock()

	h.stats = stats
}

func (h *httpHandler) updateStatsLoop() {
	for r := range h.requestStats {
		stats := h.getStats()
		stats.TotalRequests += r.TotalRequests
		stats.TotalInstances += r.TotalInstances
		stats.SucceededRequests += r.SucceededRequests
		h.setStats(stats)
	}
}

func reportError(w http.ResponseWriter, message string) {
	w.WriteHeader(500)
	io.WriteString(w, message)
	fmt.Fprintln(os.Stderr, message)
}

func reportFormatError(w http.ResponseWriter, format string, args ...interface{}) {
	w.WriteHeader(500)
	io.WriteString(w, fmt.Sprintf(format, args))
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args))
}

func reportJSON(value interface{}, name string, w http.ResponseWriter) {
	simpleJSON, err := json.Marshal(value)
	if err != nil {
		reportFormatError(w, "could not marshal %v", name)
		return
	}

	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, simpleJSON, "", "    "); err != nil {
		reportFormatError(w, "could indent %v json", name)
		return
	}

	w.WriteHeader(200)
	io.WriteString(w, prettyJson.String())
}

func (h *httpHandler) handleStatsRequest(w http.ResponseWriter, _ *http.Request) {
	reportJSON(h.getStats(), "stats", w)
}

func storeModelRequested(r *http.Request) bool {
	storeNeeded := r.URL.Query().Get("store")
	return storeNeeded == "1" || storeNeeded == "true"
}

func (h* httpHandler) handleCalculationRequest(w http.ResponseWriter, r *http.Request) {
	requestInfo := ExecutionStats{
		TotalRequests: 1,
	}
	defer func() {
		h.requestStats <- requestInfo
	}()

	argStr := r.URL.Query().Get("arg")
	if len(argStr) == 0 {
		reportError(w, "arg key is required")
		return
	}

	modelName := r.URL.Query().Get("model")
	if len(modelName) == 0 {
		reportError(w, "model key is required")
		return
	}

	model, fromCache, err := h.modelsStorage.getSLRModel(r.Context(), modelName)
	if err != nil {
		reportFormatError(w, "error loading model %v: %v", modelName, err)
		return
	}

	arg, err := strconv.ParseFloat(argStr, 64)
	if err != nil {
		reportFormatError(w, "error converting arg parameter to float: %v", argStr)
		return
	}

	modelValue := ModelValue{
		Value:           model.Calculate(arg),
		Argument:        arg,
		Model:           model,
		CalculationTime: time.Now(),
		FromCache:       fromCache,
	}
	requestInfo.SucceededRequests = 1

	reportJSON(modelValue, modelName, w)
}

func (h *httpHandler) handleTrainingRequest(w http.ResponseWriter, r *http.Request) {
	requestInfo := ExecutionStats{
		TotalRequests: 1,
	}
	defer func() {
		h.requestStats <- requestInfo
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		reportError(w,"could not load request's body")
		return
	}

	var instances [][]float64
	if err := json.Unmarshal(body, &instances); err != nil {
		reportError(w,"could not load json")
		return
	}
	requestInfo.TotalInstances = len(instances)

	var slr SimpleLinearRegression
	for idx, instance := range instances {
		if len(instance) == 2 {
			slr.AddInstance(instance[0], instance[1])
		} else if len(instance) == 3 {
			slr.AddWeightedInstance(instance[0], instance[1], instance[2])
		} else {
			reportFormatError(w,"error processing instance #%v: must contain exactly two elements", idx)
			return
		}
	}

	trainingResults := TrainingResults{
		Model: slr.Train(),
		SumSquaredErrors: slr.SumSquaredErrors(),
	}

	if storeModelRequested(r) {
		name, commitTime, err := h.modelsStorage.saveSLRModel(r.Context(), trainingResults.Model)
		if err != nil {
			trainingResults.Error = fmt.Sprintf("%v", err)
		}
		trainingResults.Name = name
		trainingResults.CreationTime = commitTime
	}
	reportJSON(trainingResults, "training results", w)

	requestInfo.SucceededRequests = 1
}

func runHTTPHandler() {
	ctx, err := handlerContext(httpMode)
	if err != nil {
		log.Fatal("cannot create context: ", err)
	}

	h, err := newHTTPHandler(ctx)
	if err != nil {
		log.Fatal("cannot create handler: ", err)
	}

	http.Handle("/train", http.HandlerFunc(h.handleTrainingRequest))
	http.Handle("/calc", http.HandlerFunc(h.handleCalculationRequest))
	http.Handle("/stats", http.HandlerFunc(h.handleStatsRequest))

	port := ctx.Value("port")
	err = http.ListenAndServe(":" + port.(string), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
