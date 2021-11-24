package client

import (
	"math"
	"sort"
	"strconv"
	"time"
)

// return min, max, mean
func calcStat(nums []float64) (float64, float64, float64) {
	if len(nums) == 0 {
		return 0.0, 0.0, 0.0
	}

	var total float64
	var min float64 = math.MaxFloat64
	var max float64
	length := len(nums)
	for _, item := range nums {
		total += item
		if min > item {
			min = item
		}
		if max < item {
			max = item
		}
	}
	mean := total / float64(length)
	return min, max, mean
}

func calcMean(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}

	var total float64
	length := len(nums)
	for _, item := range nums {
		total += item
	}
	mean := total / float64(length)
	return mean
}

func calcMedian(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}

	sort.Float64s(nums)

	isEven := len(nums)%2 == 0
	mNumber := len(nums) / 2

	if !isEven {
		return nums[mNumber]
	}
	return (float64(nums[mNumber-1]) + float64(nums[mNumber])) / 2
}

func calcVarience(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}

	var variance float64
	mean := calcMean(nums)

	for index := range nums {
		variance += math.Pow(nums[index]-mean, 2)
	}

	return variance / float64(len(nums))
}

func calcStdDev(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}

	variance := calcVarience(nums)
	return math.Sqrt(variance)
}

func calc95Percentile(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}

	sort.Float64s(nums)
	nineFive := float64(len(nums)-1) * 0.95

	newSlice := nums[int(nineFive):]
	//return strconv.Itoa(newSlice[0])
	//return strconv.FormatFloat(newSlice[0], 'f', 0, 64)
	return math.Round(newSlice[0])
}

func calc99Percentile(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}

	sort.Float64s(nums)
	nineFive := float64(len(nums)-1) * 0.99

	newSlice := nums[int(nineFive):]
	//return strconv.Itoa(newSlice[0])
	//return strconv.FormatFloat(newSlice[0], 'f', 0, 64)
	return math.Round(newSlice[0])
}

func requestsPerSecond(request int, duration time.Duration) float64 {
	convertedDuration := float64(duration) / float64(time.Second)
	toS := strconv.FormatFloat(float64(request)/convertedDuration, 'f', 2, 64)
	return stringToFloat(toS)
}

func failedRequests(slice []int) int {
	non200 := 0

	for _, item := range slice {
		if item > 226 {
			non200++
		}
	}
	return non200
}

func respStatusMap(slice []int) (map[string]int, map[string]int) {
	successMap := make(map[string]int)
	failedMap := make(map[string]int)

	for _, item := range slice {
		if item > 226 {
			failedMap[strconv.Itoa(item)]++
		} else {
			successMap[strconv.Itoa(item)]++
		}
	}
	return successMap, failedMap
}
