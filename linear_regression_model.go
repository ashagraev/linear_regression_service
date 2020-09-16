package main

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
}

// Apply() returns the result of applying the model to an argument.
func (srm *SimpleRegressionModel) Apply(arg float64) float64 {
	return srm.Coefficient * arg + srm.Intercept
}
