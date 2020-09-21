package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	pb "linear_regression_service/github.com/ashagraev/linear_regression"
)

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

func loadProtoInstancesFromTSV(reader io.Reader) ([]*pb.Instance, error){
	var instances []*pb.Instance

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

		instances = append(instances, &pb.Instance{Argument: instance[0], Target: instance[1], Weight: 1})
		lineIdx++
	}

	if len(instances) == 0 {
		return nil, errors.New("no instances loaded")
	}

	return instances, nil
}
