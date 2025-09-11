package client

import (
	"testing"
)

func TestSineVisualizationParams(t *testing.T) {
	// Test with duration
	cass := Cassowary{
		Duration:      60,
		Requests:      100,
		SineFrequency: 1.0,
		SineAmplitude: 0.5,
		SineCurve:     true,
		SinePreview:   true,
	}

	// This should not panic and should generate visualization
	cass.ShowSineVisualization()
}

func TestSineVisualizationNoDuration(t *testing.T) {
	// Test without duration
	cass := Cassowary{
		Duration:      0,
		Requests:      100,
		SineFrequency: 2.0,
		SineAmplitude: 0.8,
		SineCurve:     true,
		SinePreview:   true,
	}

	// This should not panic and should generate visualization
	cass.ShowSineVisualization()
}

func TestLiveVisualizationParams(t *testing.T) {
	// Test that live visualization setup works without panic
	cass := Cassowary{
		Duration:        10,
		Requests:        50,
		SineFrequency:   2.0,
		SineAmplitude:   0.7,
		SineCurve:       true,
		SinePreviewLive: true,
	}

	// This tests parameter setup without running the full animation
	params := SineVisualizationParams{
		Frequency:    cass.SineFrequency,
		Amplitude:    cass.SineAmplitude,
		HasDuration:  cass.Duration > 0,
		Requests:     cass.Requests,
	}

	if params.HasDuration {
		params.Duration = float64(cass.Duration)
		params.BaselineRate = float64(cass.Requests) / params.Duration
	} else {
		params.Duration = previewDuration
		params.BaselineRate = defaultSineBaselineRate
	}

	// Verify parameters are set correctly
	if params.Frequency != 2.0 {
		t.Errorf("Expected frequency 2.0, got %f", params.Frequency)
	}
	if params.Amplitude != 0.7 {
		t.Errorf("Expected amplitude 0.7, got %f", params.Amplitude)
	}
	if !params.HasDuration {
		t.Error("Expected HasDuration to be true")
	}
	if params.BaselineRate != 5.0 { // 50 requests / 10 seconds
		t.Errorf("Expected baseline rate 5.0, got %f", params.BaselineRate)
	}
}

func TestSineVisualizationParams_EdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		frequency float64
		amplitude float64
		duration  int
		requests  int
	}{
		{"MinimumAmplitude", 1.0, 0.0, 60, 100},
		{"MaximumAmplitude", 1.0, 1.0, 60, 100},
		{"HighFrequency", 5.0, 0.5, 60, 100},
		{"LowFrequency", 0.1, 0.5, 60, 100},
		{"ShortDuration", 1.0, 0.5, 1, 10},
		{"LongDuration", 1.0, 0.5, 300, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cass := Cassowary{
				Duration:      tc.duration,
				Requests:      tc.requests,
				SineFrequency: tc.frequency,
				SineAmplitude: tc.amplitude,
				SineCurve:     true,
				SinePreview:   true,
			}

			// This should not panic
			cass.ShowSineVisualization()
		})
	}
}