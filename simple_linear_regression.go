package main

// SimpleLinearRegression provides interface for solving simple linear regression problems
// https://en.wikipedia.org/wiki/Simple_linear_regression
type SimpleLinearRegression struct {
	sumWeights KahanAdder

	featureMean float64
	featureDev float64

	targetMean float64
	targetDev float64

	covariance float64
}

// AddInstance() adds one training example for the model.
func (slr *SimpleLinearRegression) AddInstance(feature float64, target float64) {
	slr.AddWeightedInstance(feature, target, 1)
}

// AddWeightedInstance() adds one weighted training example for the model.
func (slr *SimpleLinearRegression) AddWeightedInstance(feature float64, target float64, weight float64) {
	slr.sumWeights.Add(weight)
	sumWeights := slr.sumWeights.Get()
	if sumWeights <= 0 {
		return
	}

	wfd := weight * (feature - slr.featureMean)
	slr.featureMean += wfd / sumWeights
	slr.featureDev += wfd * (feature - slr.featureMean)

	wtd := weight * (target - slr.targetMean)
	slr.targetMean += wtd / sumWeights
	slr.targetDev += wtd * (target - slr.targetMean)

	slr.covariance += wfd * (target - slr.targetMean)
}

// Solve() builds a regression model according to the collected training data.
func (slr *SimpleLinearRegression) Solve() *SimpleRegressionModel {
	if slr.featureDev == 0 {
		return &SimpleRegressionModel{Coefficient: 0, Intercept: slr.targetMean}
	}

	srm := SimpleRegressionModel{Coefficient: slr.covariance / slr.featureDev}
	srm.Intercept = slr.targetMean - srm.Coefficient * slr.featureMean
	return &srm
}

// SumSquaredErrors() returns sum of squared errors on training data for the resulting model
func (slr* SimpleLinearRegression) SumSquaredErrors() float64 {
	srm := slr.Solve()
	return srm.Coefficient * srm.Coefficient * slr.featureDev - 2 * srm.Coefficient * slr.covariance + slr.targetDev
}
