# Notes

**Author**: Angel Dhakal
**Date**: 2024-08-15
**File**: insertion_sort.go

## Problem
Insertion sort an array

## Solution
start at index 1, put element of i in temp and create a space at ith index, and compare the temp to the previous element,
and if pe>temp, shift the pe to the right once(in the gap), update the gap to be -1 than last time, and
again check if space>0 and pe>temp and continue. and if space is 0 or temp>pe then put temp in the final space

**Time Complexity**: O(n^2)
**Space Complexity**: O(1)

## Note


## Code
```go
package example

func InsertionSort(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		space := i
		temp := arr[i]
		for space > 0 && arr[space-1] > temp {
			arr[space] = arr[space-1]
			space = space - 1
		}
		arr[space] = temp
	}
}

```
