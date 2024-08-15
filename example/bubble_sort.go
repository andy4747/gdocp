/*
Author: Angel Dhakal
Date: 2024-08-14
File: bubble_sort.go
Problem: Bubbble sort the given array
Solution: {
Iterate through the array, compare adjacent elements, swap the elements if out of order, repeat until no swaps needed.
after each iteration of i, the leargest/smalles elements (in asc/desc respectively order) is always at the end of the array,
so, no need to iterate over the sorted portion.
}
Note:{
...
...
}
Time Complexity: O(n^2)
Space Complexity: O(1)
*/

package example

func BubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}
