package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

type requestStat struct {
	success   bool
	instances int
}

type handler struct {
	stats        ExecutionStats
	requestStats chan requestStat

	modelsStorage *ModelsStorage
}

func newHandler(ctx context.Context) (*handler, error) {
	modelsStorage, err := NewModelsStorage(ctx)
	if err != nil {
		return nil, err
	}
	h := handler{requestStats: make(chan requestStat), modelsStorage: modelsStorage}
	go h.updateStatsLoop()
	return &h, nil
}

func (h *handler) updateStatsLoop() {
	for r := range h.requestStats {
		h.stats.TotalRequests++
		h.stats.TotalInstances += r.instances
		if r.success {
			h.stats.SucceededRequests++
		}
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

func (h*handler) handleStatsRequest(w http.ResponseWriter, _ *http.Request) {
	reportJSON(h.stats, "stats", w)
}

func storeModelRequested(r *http.Request) bool {
	storeNeeded := r.URL.Query().Get("store")
	return storeNeeded == "1" || storeNeeded == "true"
}

func (h* handler) handleApplyRequest(w http.ResponseWriter, r *http.Request) {
	var requestInfo requestStat
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

	model, err := h.modelsStorage.GetSLRModel(r.Context(), modelName)
	if err != nil {
		reportFormatError(w, "error loading model %v: %v", modelName, err)
		return
	}

	arg, err := strconv.ParseFloat(argStr, 64)
	if err != nil {
		reportFormatError(w, "error converting arg parameter to float: %v", argStr)
		return
	}

	modelValue := model.Evaluate(arg)
	reportJSON(modelValue, modelName, w)
}

func (h *handler) handleSolveRequest(w http.ResponseWriter, r *http.Request) {
	var requestInfo requestStat
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
	requestInfo.instances = len(instances)

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

	solveResult := Solution{
		Model: slr.Solve(),
		SumSquaredErrors: slr.SumSquaredErrors(),
	}

	if storeModelRequested(r) {
		name, commitTime, err := h.modelsStorage.SaveSLRModel(r.Context(), solveResult.Model)
		if err != nil {
			solveResult.Error = fmt.Sprintf("%v", err)
		}
		solveResult.Name = name
		solveResult.CreationTime = commitTime
	}
	reportJSON(solveResult, "solution", w)

	requestInfo.success = true
}
