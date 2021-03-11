package client

import (
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func getBins(dataPoints int) int {
	b := int(math.RoundToEven(math.Sqrt(float64(dataPoints))))
	return b
}

// PlotHistogram outputs a histogram png
func (c *Cassowary) PlotHistogram(durations []float64) error {

	bins := getBins(len(durations))
	if bins == 0 {
		return nil
	}

	// remove outliers
	avg := calcMean(durations)
	stddev := calcStdDev(durations)
	outliers := avg + 1.95*stddev

	for i := len(durations) - 1; i >= 0; i-- {
		if durations[i] > outliers {
			durations = append(durations[:i], durations[i+1:]...)
		}
	}

	vs := make(plotter.Values, len(durations))
	for i, d := range durations {
		vs[i] = d
	}

	h, err := plotter.NewHist(vs, bins)
	if err != nil {
		return err
	}

	p := plot.New()

	p.Add(h)
	p.Title.Text = "Distribution"
	p.Y.Label.Text = "Requests"
	p.X.Label.Text = "ms"
	h.FillColor = color.RGBA{R: 70, G: 130, B: 180, A: 1}
	h.Color = color.Opaque

	if err := p.Save(512, 512, "hist.png"); err != nil {
		return err
	}

	return nil
}
