package main

import (
	"fmt"
	"sort"
	"strconv"
)

func calcMean(nums []int) float64 {
	var total int
	length := len(nums)
	for _, item := range nums {
		total += item
	}
	mean := float64(total) / float64(length)
	return mean
}

func calcMedian(nums []int) float64 {
	sort.Ints(nums)

	isEven := len(nums)%2 == 0
	mNumber := len(nums) / 2

	if !isEven {
		fmt.Println("is odd")
		return float64(nums[mNumber])
	}
	return (float64(nums[mNumber-1]) + float64(nums[mNumber])) / 2
}

func return95Median(nums []int) string {
	sort.Ints(nums)
	nineFive := float64(len(nums)-1) * 0.95

	newSlice := nums[int(nineFive):]
	return strconv.Itoa(newSlice[0])
}
