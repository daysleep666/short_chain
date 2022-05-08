package main

import (
	"fmt"

	"go.uber.org/atomic"
)

func findDuplicates(nums []int) (result []int) {
	for i := 0; i < len(nums); {
		if nums[i]-1 == i || nums[i] == -1 {
			i++
			continue
		}
		idx := nums[i] - 1
		num := nums[idx]
		if nums[i] == num {
			result = append(result, num)
			nums[i] = -1
			i++
		} else {
			nums[i], nums[idx] = nums[idx], nums[i]
		}
	}
	return
}

func main() {
	var number atomic.Int64
	for {

		fmt.Println(number.Inc() & 4095)
	}

}
