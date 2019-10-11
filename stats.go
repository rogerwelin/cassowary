package main

import (
	"math"
	"sort"
	"strconv"
	"time"
)

func calcMean(nums []float64) float64 {
	var total float64
	length := len(nums)
	for _, item := range nums {
		total += item
	}
	mean := total / float64(length)
	return mean
}

func calcMedian(nums []float64) float64 {
	sort.Float64s(nums)

	isEven := len(nums)%2 == 0
	mNumber := len(nums) / 2

	if !isEven {
		return nums[mNumber]
	}
	return (float64(nums[mNumber-1]) + float64(nums[mNumber])) / 2
}

func calcVarience(nums []float64) float64 {
	var variance float64
	mean := calcMean(nums)

	for index := range nums {
		variance += math.Pow(nums[index]-mean, 2)
	}

	return variance / float64(len(nums))
}

func calcStdDev(nums []float64) float64 {
	variance := calcVarience(nums)
	return math.Sqrt(variance)
}

func calc95Percentile(nums []float64) string {
	sort.Float64s(nums)
	nineFive := float64(len(nums)-1) * 0.95

	newSlice := nums[int(nineFive):]
	//return strconv.Itoa(newSlice[0])
	return strconv.FormatFloat(newSlice[0], 'f', 0, 64)
}

func requestsPerSecond(request int, duration time.Duration) string {
	convertedDuration := float64(duration) / float64(time.Second)
	return strconv.FormatFloat(float64(request)/convertedDuration, 'f', 6, 64)
}

func failedRequests(slice []int) string {
	non200 := 0

	for _, item := range slice {
		if item > 226 {
			non200++
		}
	}
	return strconv.Itoa(non200)
}
