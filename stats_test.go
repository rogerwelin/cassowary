package main

import (
	"testing"
)

var testNums = []struct {
	in             []int
	expectedMean   float64
	expectedMedian float64
}{
	//	{[]int{-10, 0, 10, 20, 30}, 10, 10},
	//	{[]int{8, 9, 10, 11, 12}, 10, 10},
	//	{[]int{40, 10, 20, 30, 1}, 20.20, 20},
	{[]int{3, 2, 2}, 2.3333333333333335, 2},
	{[]int{3, 2, 1, 9}, 3.75, 2.5},
}

func TestCalcMean(t *testing.T) {
	for _, tt := range testNums {
		actual := calcMean(tt.in)
		if actual != tt.expectedMean {
			t.Errorf("calcMean(%d): expected %f, actual %f", tt.in, tt.expectedMean, actual)
		}
	}
}

func TestCalcMedian(t *testing.T) {
	for _, tt := range testNums {
		actual := calcMedian(tt.in)
		if actual != tt.expectedMedian {
			t.Errorf("calcMean(%d): expected %f, actual %f", tt.in, tt.expectedMedian, actual)
		}
	}
}
