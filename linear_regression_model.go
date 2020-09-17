package main

import (
	"errors"
	"time"
)

// SimpleRegressionModel represents simple regression model with one linear coefficient and the intercept.
// It has the following form: f(x) = a * x + b
type SimpleRegressionModel struct {
	Coefficient float64
	Intercept float64
}

// Solution stores simple linear regression solution information
type Solution struct {
	// Model is a simple regression model which fits the training data best.
	Model *SimpleRegressionModel

	// SumSquaredErrors stores the model's sum of squared errors over the training data.
	SumSquaredErrors float64

	// Name stores the name of the model stored in Spanner database
	Name			string	`json:"Name,omitempty"`

	// Error stores the error message
	Error			string	`json:"Error,omitempty"`

	// CreationTime stores the creation time of the stored model
	CreationTime	*time.Time	`json:"CreationTime,omitempty"`
}

// Apply() returns the result of applying the model to an argument.
func (srm *SimpleRegressionModel) Apply(arg float64) float64 {
	return srm.Coefficient * arg + srm.Intercept
}

// ToFloatArray() converts a simple regression model to an array of float parameters.
func (srm *SimpleRegressionModel) ToFloatArray() []float64 {
	return []float64{srm.Coefficient, srm.Intercept}
}

// NewSimpleRegressionModel() converts an array of float parameters to a simple regression model
func NewSimpleRegressionModel(params []float64) (*SimpleRegressionModel, error) {
	if len(params) != 2 {
		return nil, errors.New("simple regression model must have exactly two params")
	}

	return &SimpleRegressionModel{Coefficient: params[0], Intercept: params[1]}, nil
}
