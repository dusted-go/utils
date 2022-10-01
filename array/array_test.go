package array

import (
	"testing"

	"github.com/dusted-go/utils/assert"
)

func Test_ContainsMoreThan(t *testing.T) {

	testCases := []struct {
		Expected bool
		Array    []int
		MoreThan []int
	}{
		{Expected: true, Array: []int{1, 2, 3, 4}, MoreThan: []int{1, 2}},
		{Expected: true, Array: []int{1, 2, 3, 4}, MoreThan: []int{1, 4}},
		{Expected: true, Array: []int{1, 2, 3, 4}, MoreThan: []int{3}},
		{Expected: true, Array: []int{1, 2, 3, 4}, MoreThan: []int{1}},
		{Expected: false, Array: []int{1}, MoreThan: []int{1}},
		{Expected: false, Array: []int{1, 2, 3}, MoreThan: []int{1, 3, 2}},
	}

	for _, testCase := range testCases {
		actual := ContainsMoreThan(testCase.Array, testCase.MoreThan...)
		assert.Equal(t, testCase.Expected, actual)
	}
}
