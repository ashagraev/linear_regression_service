package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func newTrainingHTTPClient() *regressionClient {
	return newRegressionClient("train","http","train model")
}

func newCalculatingHTTPClient() *regressionClient {
	return newRegressionClient("apply","http","apply model")
}

func (rc *regressionClient) requestHTTPTraining(instances [][]float64) (string, error) {
	data, err := json.Marshal(instances)
	if err != nil {
		return "", fmt.Errorf("can't marshal instances: %v", err)
	}

	dataReader := bytes.NewReader(data)
	resp, err := http.Post(rc.serverPath + "/train?store=1", "application/json", dataReader)
	if err != nil {
		return "", fmt.Errorf("error processing /train: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't load /train answer: %v", err)
	}

	return string(body), nil
}

func (rc *regressionClient) requestHTTPCalculation(arg float64) (string, error) {
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

func runHTTPTrain() {
	client := newTrainingHTTPClient()

	instances, err := loadInstancesFromTSV(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.requestHTTPTraining(instances)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func runHTTPApply() {
	client := newCalculatingHTTPClient()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		arg, err := strconv.ParseFloat(text, 64)
		if err != nil {
			log.Fatalf("invalid float: %v", text)
		}

		result, err := client.requestHTTPCalculation(arg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	}
}
