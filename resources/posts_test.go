package resources

import "testing"

func TestSlicing(t *testing.T) {
	limit := 10
	offset := 3
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	start := offset
	end := offset + limit
	if end > len(numbers) {
		end = len(numbers)
	}
	sliced := numbers[start:end]
	t.Logf("sliced: %+v", sliced)
	t.Fail()
}
