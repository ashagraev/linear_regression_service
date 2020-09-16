package main

// KahanAdder implements Kahan's summation algorithm https://en.wikipedia.org/wiki/Kahan_summation_algorithm.
type KahanAdder struct {
	Sum      float64
	Residual float64
}

// Add() adds a single value to the sum.
func (ka *KahanAdder) Add(value float64) *KahanAdder {
	y := value - ka.Residual
	t := ka.Sum + y
	ka.Residual = (t - ka.Sum) - y
	ka.Sum = t
	return ka
}

// Get() gets the corrected value of the sum.
func (ka *KahanAdder) Get() float64 {
	return ka.Sum + ka.Residual
}
