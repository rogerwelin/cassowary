package main

import (
	"sort"
	"strconv"
)

func calcMean(nums []int) string {
	var total int
	length := len(nums)
	for _, item := range nums {
		total += item
	}
	mean := float64(total) / float64(length)
	return strconv.FormatFloat(mean, 'f', 2, 64)
}

func calcMedian(nums []int) string {
	sort.Ints(nums)

	even := len(nums)%2 == 0

	if !even {
		firstMiddleMost := (len(nums) - 1) / 2
		secondMiddleMost := (len(nums)-1)/2 + 1
		median := (firstMiddleMost + secondMiddleMost) / 2
		return strconv.Itoa(nums[median])
	}
	return strconv.Itoa(nums[(len(nums)-1)/2])
}

func return95Median(nums []int) string {
	sort.Ints(nums)
	nineFive := float64(len(nums)-1) * 0.95

	newSlice := nums[int(nineFive):]
	return strconv.Itoa(newSlice[0])
}
