package main

import (
	"testing"
)

var meanTest = []struct {
	in       []int
	expected float64
}{
	{[]int{-10, 0, 10, 20, 30}, 10.00},
	{[]int{8, 9, 10, 11, 12}, 10.00},
	{[]int{1, 10, 20, 30, 40}, 20.20},
	{[]int{2, 2, 3}, 2.3333333333333335},
}

func TestCalcMean(t *testing.T) {
	for _, tt := range meanTest {
		actual := calcMean(tt.in)
		if actual != tt.expected {
			t.Errorf("calcMean(%d): expected %f, actual %f", tt.in, tt.expected, actual)
		}

	}
}
