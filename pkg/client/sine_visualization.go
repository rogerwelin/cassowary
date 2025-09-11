package client

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
)

const (
	// Visualization constants
	chartWidth       = 80   // Width of ASCII chart in characters
	chartHeight      = 20   // Height of ASCII chart in characters
	timeAxisPoints   = 10   // Number of time points to show on X-axis
	previewDuration  = 60.0 // Default preview duration in seconds for visualization
)

// SineVisualizationParams holds parameters for generating the sine wave preview
type SineVisualizationParams struct {
	Frequency      float64
	Amplitude      float64
	Duration       float64
	BaselineRate   float64
	HasDuration    bool
	Requests       int
}

// ShowSineVisualization displays an ASCII visualization of the sine wave pattern
func (c *Cassowary) ShowSineVisualization() {
	c.generateSineVisualization()
}

// generateSineVisualization creates an ASCII visualization of the sine wave pattern  
func (c *Cassowary) generateSineVisualization() {
	params := SineVisualizationParams{
		Frequency:    c.SineFrequency,
		Amplitude:    c.SineAmplitude,
		HasDuration:  c.Duration > 0,
		Requests:     c.Requests,
	}

	if params.HasDuration {
		params.Duration = float64(c.Duration)
		params.BaselineRate = float64(c.Requests) / params.Duration
	} else {
		params.Duration = previewDuration
		params.BaselineRate = defaultSineBaselineRate
	}

	// Generate ASCII chart
	chart := generateASCIIChart(params)
	
	// Print header with configuration
	printVisualizationHeader(params)
	
	// Print the chart
	fmt.Println(chart)
	
	// Print footer with rate information
	printVisualizationFooter(params)
}

// generateASCIIChart creates the ASCII representation of the sine wave
func generateASCIIChart(params SineVisualizationParams) string {
	var chart strings.Builder
	
	// Calculate rate bounds
	minRate := params.BaselineRate * (1.0 - params.Amplitude)
	maxRate := params.BaselineRate * (1.0 + params.Amplitude)
	
	// Ensure positive rates
	if minRate <= 0 {
		minRate = params.BaselineRate * minRateMultiplier
	}
	
	// Generate data points for the chart
	dataPoints := make([]float64, chartWidth)
	for i := 0; i < chartWidth; i++ {
		// Calculate time position (0 to 1)
		timeProgress := float64(i) / float64(chartWidth-1)
		
		var sinePhase float64
		if params.HasDuration {
			// With duration: sine wave over the full duration
			sinePhase = 2 * math.Pi * params.Frequency * timeProgress
		} else {
			// Without duration: sine wave over preview duration
			elapsed := timeProgress * params.Duration
			sinePhase = 2 * math.Pi * params.Frequency * elapsed / secondsPerMinute
		}
		
		sineFactor := math.Sin(sinePhase)
		rate := params.BaselineRate * (1.0 + params.Amplitude*sineFactor)
		
		if rate <= 0 {
			rate = params.BaselineRate * minRateMultiplier
		}
		
		dataPoints[i] = rate
	}
	
	// Generate the ASCII chart
	for row := 0; row < chartHeight; row++ {
		// Calculate the rate threshold for this row (from top to bottom)
		rowRate := maxRate - (float64(row)/float64(chartHeight-1))*(maxRate-minRate)
		
		// Y-axis label
		if row == 0 {
			chart.WriteString(fmt.Sprintf("%5.1f ┤", maxRate))
		} else if row == chartHeight/2 {
			chart.WriteString(fmt.Sprintf("%5.1f ┤", params.BaselineRate))
		} else if row == chartHeight-1 {
			chart.WriteString(fmt.Sprintf("%5.1f ┤", minRate))
		} else {
			chart.WriteString("      │")
		}
		
		// Plot the data points
		for i := 0; i < chartWidth; i++ {
			rate := dataPoints[i]
			
			if i == 0 {
				// First character - always part of the line
				if rate >= rowRate {
					if row == chartHeight/2 {
						chart.WriteString(color.CyanString("─"))
					} else {
						chart.WriteString(color.GreenString("╭"))
					}
				} else {
					chart.WriteString(" ")
				}
			} else {
				prevRate := dataPoints[i-1]
				
				// Determine what character to use based on rate position
				if rate >= rowRate && prevRate >= rowRate {
					// Both points are above threshold
					chart.WriteString(color.GreenString("─"))
				} else if rate >= rowRate && prevRate < rowRate {
					// Rising edge
					chart.WriteString(color.GreenString("╱"))
				} else if rate < rowRate && prevRate >= rowRate {
					// Falling edge  
					chart.WriteString(color.GreenString("╲"))
				} else {
					// Both points below threshold
					if row == chartHeight/2 {
						// Baseline reference line
						chart.WriteString(color.CyanString("─"))
					} else {
						chart.WriteString(" ")
					}
				}
			}
		}
		
		// Add rate info on the right for key rows
		if row == 0 {
			chart.WriteString(color.RedString(" ← Peak Rate"))
		} else if row == chartHeight/2 {
			chart.WriteString(color.CyanString(" ← Baseline"))
		} else if row == chartHeight-1 {
			chart.WriteString(color.YellowString(" ← Valley Rate"))
		}
		
		chart.WriteString("\n")
	}
	
	// Add time axis
	chart.WriteString("      └")
	for i := 0; i < chartWidth; i++ {
		chart.WriteString("─")
	}
	chart.WriteString("\n")
	
	// Add time markers
	timeMarkersLine := "       "
	lastPos := 0
	
	for i := 0; i <= timeAxisPoints; i++ {
		// Calculate position for this time marker
		timeProgress := float64(i) / float64(timeAxisPoints)
		timeValue := timeProgress * params.Duration
		timeStr := fmt.Sprintf("%.0fs", timeValue)
		
		// Calculate character position
		targetPos := int(float64(i) * float64(chartWidth) / float64(timeAxisPoints))
		
		// Add spaces to reach the target position
		spacesNeeded := targetPos - lastPos
		if i > 0 && spacesNeeded > 0 {
			timeMarkersLine += strings.Repeat(" ", spacesNeeded)
		}
		
		timeMarkersLine += timeStr
		lastPos = targetPos + len(timeStr)
	}
	
	chart.WriteString(timeMarkersLine)
	chart.WriteString("\n")
	
	return chart.String()
}

// printVisualizationHeader prints the header information for the visualization
func printVisualizationHeader(params SineVisualizationParams) {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	fmt.Printf("\n%s\n", cyan("🌊 Sine Wave Load Pattern Preview"))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	if params.HasDuration {
		fmt.Printf("Mode: %s | ", yellow("Duration-based"))
		fmt.Printf("Duration: %s | ", yellow(fmt.Sprintf("%.0fs", params.Duration)))
		fmt.Printf("Requests: %s\n", yellow(fmt.Sprintf("%d", params.Requests)))
	} else {
		fmt.Printf("Mode: %s | ", yellow("Request-based"))
		fmt.Printf("Requests: %s | ", yellow(fmt.Sprintf("%d", params.Requests)))
		fmt.Printf("Preview Duration: %s\n", yellow(fmt.Sprintf("%.0fs", params.Duration)))
	}
	
	fmt.Printf("Frequency: %s cycles | ", yellow(fmt.Sprintf("%.1f", params.Frequency)))
	fmt.Printf("Amplitude: %s | ", yellow(fmt.Sprintf("%.1f", params.Amplitude)))
	fmt.Printf("Baseline Rate: %s\n\n", yellow(fmt.Sprintf("%.1f req/s", params.BaselineRate)))
	
	fmt.Printf("Request Rate Over Time:\n\n")
}

// printVisualizationFooter prints footer information and usage tips
func printVisualizationFooter(params SineVisualizationParams) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	
	minRate := params.BaselineRate * (1.0 - params.Amplitude)
	maxRate := params.BaselineRate * (1.0 + params.Amplitude)
	
	if minRate <= 0 {
		minRate = params.BaselineRate * minRateMultiplier
	}
	
	fmt.Printf("\n📊 Pattern Analysis:\n")
	fmt.Printf("• %s: %.1f req/s (%.0f%% of baseline)\n", 
		red("Peak Rate"), maxRate, (maxRate/params.BaselineRate)*100)
	fmt.Printf("• %s: %.1f req/s\n", 
		cyan("Baseline Rate"), params.BaselineRate)  
	fmt.Printf("• %s: %.1f req/s (%.0f%% of baseline)\n", 
		yellow("Valley Rate"), minRate, (minRate/params.BaselineRate)*100)
	
	if params.HasDuration {
		cycleTime := params.Duration / params.Frequency
		fmt.Printf("• %s: %.1f seconds per complete cycle\n", 
			green("Cycle Duration"), cycleTime)
	} else {
		cycleTime := secondsPerMinute / params.Frequency  
		fmt.Printf("• %s: %.1f seconds per complete cycle\n", 
			green("Cycle Duration"), cycleTime)
	}
	
	fmt.Printf("\n💡 Tips:\n")
	fmt.Printf("• Use %s for gentler variation\n", cyan("--sine-amplitude 0.2"))
	fmt.Printf("• Use %s for more dramatic changes\n", cyan("--sine-amplitude 0.8"))
	fmt.Printf("• Use %s for faster oscillations\n", cyan("--sine-frequency 2.0"))
	fmt.Printf("• Use %s for slower waves\n", cyan("--sine-frequency 0.5"))
	
	fmt.Printf("\n🚀 Ready to run the actual test? Remove %s flag!\n", cyan("--sine-preview"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}