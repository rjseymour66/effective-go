package example

import "testing"

func TestAdd(t *testing.T) {
	type testCase struct {
		a        int
		b        int
		expected int
	}
	tt := map[string]testCase{
		"test one": {
			a:        4,
			b:        5,
			expected: 9,
		},
		"test two": {
			a:        -4,
			b:        15,
			expected: 11,
		},
		"test three": {
			a:        5,
			b:        1,
			expected: 6,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// act
			got := add(tc.a, tc.b)

			// assert
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}
