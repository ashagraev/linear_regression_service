package main

import (
	"errors"
	"time"
)

// SimpleRegressionModel represents simple regression model with one linear coefficient and the intercept.
// It has the following form: f(x) = a * x + b
type SimpleRegressionModel struct {
	Name string `json:"Name,omitempty"`

	Coefficient float64
	Intercept float64
}

// TrainingResults stores the results of simple linear regression model training.
type TrainingResults struct {
	// Model is a simple regression model which fits the training data best.
	Model *SimpleRegressionModel

	// SumSquaredErrors stores the model's sum of squared errors over the training data.
	SumSquaredErrors float64

	// Name stores the name of the model stored in Spanner database
	Name			string	`json:"Name,omitempty"`

	// Error stores the error message
	Error			string	`json:"Error,omitempty"`

	// CreationTime stores the creation time of the stored model
	CreationTime	time.Time	`json:"CreationTime,omitempty"`
}

// ModelValue stores the information about model calculation over the given argument
type ModelValue struct {
	// Value stores the calculated model value
	Value float64

	// Argument stores the given argument value
	Argument float64

	// Model stores the requested model
	Model *SimpleRegressionModel

	// FromCache reports whether the model was taken from local cache
	FromCache bool

	// CalculationTime stores the moment of calculation
	CalculationTime time.Time
}

// Apply() returns the result of applying the model to an argument.
func (srm *SimpleRegressionModel) Apply(arg float64) float64 {
	return srm.Coefficient * arg + srm.Intercept
}

// Apply() returns the result of applying the model to an argument.
func (srm *SimpleRegressionModel) Evaluate(arg float64) ModelValue {
	return ModelValue{
		Value: srm.Apply(arg),
		Argument: arg,
		Model : srm,
		CalculationTime: time.Now(),
	}
}

// ToFloatArray() converts a simple regression model to an array of float parameters.
func (srm *SimpleRegressionModel) ToFloatArray() []float64 {
	return []float64{srm.Coefficient, srm.Intercept}
}

// NewSimpleRegressionModel() converts an array of float parameters to a simple regression model
func NewSimpleRegressionModel(params []float64, name string) (*SimpleRegressionModel, error) {
	if len(params) != 2 {
		return nil, errors.New("simple regression model must have exactly two params")
	}

	return &SimpleRegressionModel{Name: name, Coefficient: params[0], Intercept: params[1]}, nil
}
