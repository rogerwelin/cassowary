package client

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestBoxPlot(t *testing.T) {
	num := 100
	res := make([]float64, num)
	filename := "boxplot.png"
	cass := &Cassowary{}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < num; i++ {
		res = append(res, rand.Float64()*(200.0-10.0))
	}

	err := cass.PlotBoxplot(res)
	if err != nil {
		t.Errorf("Expected ok but got: %v", err)
	}
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		t.Errorf("Expected %s in current dir but got: %v", filename, err)
	}

	_ = os.Remove(filename)
}
