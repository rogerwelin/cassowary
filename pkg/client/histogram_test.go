package client

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

var testBins = []struct {
	in           int
	expectedBins int
}{
	{10, 3},
	{100, 10},
	{250, 16},
	{500, 22},
	{750, 27},
	{1000, 32},
}

func TestBins(t *testing.T) {
	for i, tt := range testBins {
		actual := getBins(tt.in)
		if actual != tt.expectedBins {
			t.Errorf("test: %d, getBins(%d): expected %d, actual %d", i+1, tt.in, tt.expectedBins, actual)
		}
	}
}

func TestHistOutput(t *testing.T) {
	num := 100
	res := make([]float64, num)
	filename := "hist.png"
	cass := &Cassowary{}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < num; i++ {
		res = append(res, rand.Float64()*(200.0-10.0))
	}

	err := cass.PlotHistogram(res)
	if err != nil {
		t.Errorf("Expected ok but got: %v", err)
	}
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		t.Errorf("Expected %s in current dir but got: %v", filename, err)
	}

	_ = os.Remove(filename)
}
