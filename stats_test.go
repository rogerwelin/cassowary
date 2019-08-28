package main

import (
	"testing"
)

var testNums = []struct {
	in             []float64
	expectedMean   float64
	expectedMedian float64
	expectedStdDev float64
}{
	{[]float64{-10, 0, 10, 20, 30}, 10, 10, 14.142135623730951},
	{[]float64{8, 9, 10, 11, 12}, 10, 10, 1.4142135623730951},
	{[]float64{40, 10, 20, 30, 1}, 20.20, 20, 13.862178760930766},
	{[]float64{3, 2, 2}, 2.3333333333333335, 2, 0.4714045207910317},
	{[]float64{3, 2, 1, 9}, 3.75, 2.5, 3.112474899497183},
}

func TestCalcMean(t *testing.T) {
	for i, tt := range testNums {
		actual := calcMean(tt.in)
		if actual != tt.expectedMean {
			t.Errorf("test: %d, calcMean(%f): expected %f, actual %f", i+1, tt.in, tt.expectedMean, actual)
		}
	}
}

func TestCalcMedian(t *testing.T) {
	for i, tt := range testNums {
		actual := calcMedian(tt.in)
		if actual != tt.expectedMedian {
			t.Errorf("test: %d, calcMean(%f): expected %f, actual %f", i+1, tt.in, tt.expectedMedian, actual)
		}
	}
}

func TestCalcStdDev(t *testing.T) {
	for i, tt := range testNums {
		actual := calcStdDev(tt.in)
		if actual != tt.expectedStdDev {
			t.Errorf("test: %d, calcMean(%f): expected %f, actual %f", i+1, tt.in, tt.expectedStdDev, actual)
		}
	}
}
