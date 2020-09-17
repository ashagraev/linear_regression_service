package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
)

type ModelsStorage struct {
	spannerClient *spanner.Client
}

func NewModelsStorage(ctx context.Context) (*ModelsStorage, error) {
	spannerClient, err := spanner.NewClient(ctx, "projects/thematic-cider-289114/instances/machine-learning/databases/models")
	if err != nil {
		return nil, fmt.Errorf("spanner.NewClient() error: %v", err)
	}

	return &ModelsStorage{spannerClient: spannerClient}, nil
}

func randomModelName() (string, error) {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("cannot create random model name: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b), err
}

func (ms *ModelsStorage) SaveSLRModel(ctx context.Context, model *SimpleRegressionModel) (string, time.Time, error) {
	name, err := randomModelName()
	if err != nil {
		return "", time.Time{}, err
	}

	commitTS, err := ms.spannerClient.Apply(ctx, []*spanner.Mutation{
		spanner.Insert("slr_models",
			[]string{"name", "params", "creation_time"},
			[]interface{}{name, model.ToFloatArray(), spanner.CommitTimestamp})})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("cannot save model to Spanner: %v", err)
	}

	return name, commitTS, err
}

func (ms *ModelsStorage) GetSLRModel(ctx context.Context, name string) (*SimpleRegressionModel, error) {
	row, err := ms.spannerClient.Single().ReadRow(ctx, "slr_models",
		spanner.Key{name}, []string{"params"})
	if err != nil {
		return nil, fmt.Errorf("error loading model from Spanner: %v", err)
	}
	var params []float64
	if err = row.Columns(&params); err != nil {
		return nil, fmt.Errorf("error loading parameters from Spanner row: %v", err)
	}

	model, err := NewSimpleRegressionModel(params, name)
	if err != nil {
		return nil, err
	}

	return model, nil
}
