package testdemo

// IsSorted reports whether data is sorted.
func IsSorted(data []int) bool {
	n := len(data)
	if n == 0 || n == 1 {
		return true
	}
	i := 0
	for i < n-1 && data[i] <= data[i+1] {
		i = i + 1
	}
	return i == n-1
}
