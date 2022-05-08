package main

import "fmt"

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
	fmt.Println(findDuplicates([]int{4, 3, 2, 7, 8, 2, 3, 1}))
	fmt.Println(findDuplicates([]int{1, 1, 2}))
	fmt.Println(findDuplicates([]int{1}))
}
