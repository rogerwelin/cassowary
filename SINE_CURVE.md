# 🌊 Sine Curve Load Testing Guide

Sine curve load testing varies the request rate over time following a mathematical sine wave, simulating realistic traffic patterns with natural peaks and valleys.

## Quick Start

```bash
# Start with presets - easiest way to get good results
$ ./cassowary run -u http://example.com -n 1000 -d 60 --sine-preset gentle
$ ./cassowary run -u http://example.com -n 1000 -d 60 --sine-preset normal  
$ ./cassowary run -u http://example.com -n 1000 -d 60 --sine-preset aggressive

# Preview the pattern first (no requests sent)
$ ./cassowary run -u http://example.com -n 1000 -d 60 --sine-preview

# Run actual sine curve test with new defaults (frequency=0.5, amplitude=0.5)
$ ./cassowary run -u http://example.com -n 1000 -d 60 --sine-curve

# Live animated preview (no requests sent)
$ ./cassowary run -u http://example.com -n 1000 -d 60 --sine-preview-live
```

## Parameters

| Parameter | Description | Default | Example Values |
|-----------|-------------|---------|----------------|
| `--sine-preset` | Use preset configuration | - | gentle, normal, aggressive |
| `--sine-frequency` | Complete cycles during test | 0.5 | 0.25 (very slow), 2.0 (fast) |
| `--sine-amplitude` | Rate variation intensity (0.0-1.0) | 0.5 | 0.2 (gentle), 0.8 (dramatic) |

### Presets
- **gentle**: `frequency=0.25, amplitude=0.3` → Gradual ramp-up to +30%
- **normal**: `frequency=0.5, amplitude=0.5` → Ramp up/down with ±50% variation  
- **aggressive**: `frequency=2.0, amplitude=0.7` → Rapid oscillations with ±70% variation

**Rate Formula**: `rate = baseline × (1 + amplitude × sin(...))`
- **Amplitude 0.5**: Rate varies from 50% to 150% of baseline
- **Amplitude 0.8**: Rate varies from 20% to 180% of baseline

## Common Patterns

```bash
# Gentle daily pattern
$ ./cassowary run -u URL -n 1000 -d 60 --sine-curve --sine-frequency 1.0 --sine-amplitude 0.3

# Stress test with spikes  
$ ./cassowary run -u URL -n 1000 -d 60 --sine-curve --sine-frequency 2.0 --sine-amplitude 0.8

# Gradual ramp-up/down
$ ./cassowary run -u URL -n 1000 -d 120 --sine-curve --sine-frequency 0.5 --sine-amplitude 0.6
```

## Visual Examples

**Single Wave (frequency=1.0, amplitude=0.5)**:
```
  15 ┤        ╭─╮        
  10 ┤───────╱───╲─────── Baseline
   5 ┤      ╱     ╲      
     └─────────────────
     0s    30s    60s
```

**Multiple Waves (frequency=3.0, amplitude=0.5)**:
```  
  15 ┤  ╭─╮  ╭─╮  ╭─╮  
  10 ┤─╱─╲╲╱╱─╲╲╱╱─╲─ Baseline
   5 ┤╱   ╲╱   ╲╱   ╲
     └─────────────────
     0s   20s   40s  60s
```

## Best Practices

1. **Always preview first**: Use `--sine-preview` to visualize the pattern
2. **Start gentle**: Begin with amplitude 0.3, then increase
3. **Use concurrency=1**: For most precise sine wave timing: `-c 1`
4. **Monitor your system**: Watch how your infrastructure responds to the varying load

## Troubleshooting

- **"Duration too short"**: Increase duration or reduce requests
- **Variations not visible**: Increase `--sine-amplitude`  
- **Timing imprecise**: Use `-c 1` for better precision
- **Need help visualizing**: Use `--sine-preview-live` to see the animated pattern

## Two Modes

**Duration-based** (with `-d`): Sends N requests over fixed time
**Request-based** (no `-d`): Sends N requests following sine timing (uses 10 req/s baseline)