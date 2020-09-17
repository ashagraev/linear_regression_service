package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type regressionClient struct {
	serverPath string
	modelName string
}

func NewRegressionClient(mode string, usage string) *regressionClient {
	flag.Bool(mode, true, usage)
	var server = flag.String("server", "", "network path of the solving server")
	var model = flag.String("model", "", "model name for calculation")
	flag.Parse()

	return &regressionClient{serverPath: *server, modelName: *model}
}

func NewTrainingClient() *regressionClient {
	return NewRegressionClient("train","train model")
}

func NewCalculatingClient() *regressionClient {
	return NewRegressionClient("apply","apply model")
}

func (rc *regressionClient) requestTraining(instances [][]float64) (string, error) {
	data, err := json.Marshal(instances)
	if err != nil {
		return "", fmt.Errorf("can't marshal instances: %v", err)
	}

	dataReader := bytes.NewReader(data)
	resp, err := http.Post(rc.serverPath + "/solve?store=1", "application/json", dataReader)
	if err != nil {
		return "", fmt.Errorf("error processing /solve: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't load /solve answer: %v", err)
	}

	return string(body), nil
}

func (rc *regressionClient) requestCalculation(arg float64) (string, error) {
	url := fmt.Sprintf("%v/apply?model=%v&arg=%v", rc.serverPath, rc.modelName, arg)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error processing /apply: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't load /apply answer: %v", err)
	}

	return string(body), nil
}

func loadInstancesFromTSV(reader io.Reader) ([][]float64, error){
	var instances [][]float64

	lineIdx := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var instance []float64
		for _, s := range strings.Fields(scanner.Text()) {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid float: %v", err)
			}
			instance = append(instance, v)
		}
		if len(instance) == 0 {
			continue
		}
		if len(instance) != 2 {
			return nil, fmt.Errorf("bad number of tokens: %v, line %v", len(instance), lineIdx)
		}

		instances = append(instances, instance)
		lineIdx++
	}

	if len(instances) == 0 {
		return nil, errors.New("no instances loaded")
	}

	return instances, nil
}