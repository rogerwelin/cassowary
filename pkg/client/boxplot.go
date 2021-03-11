package client

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// PlotBoxplot outputs a boxplot png
func (c *Cassowary) PlotBoxplot(durations []float64) error {
	p := plot.New()
	p.Title.Text = "Box Plot"

	vs := make(plotter.Values, len(durations))
	for i, d := range durations {
		vs[i] = d
	}

	box, err := plotter.NewBoxPlot(vg.Length(20), 0.0, vs)
	if err != nil {
		panic(err)
	}
	box.FillColor = color.RGBA{R: 70, G: 130, B: 180, A: 1}
	p.Add(box)

	if err := p.Save(512, 512, "boxplot.png"); err != nil {
		return err
	}

	return nil
}
