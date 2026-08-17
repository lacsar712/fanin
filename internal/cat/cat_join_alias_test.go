package cat_test

import (
	"reflect"
	"testing"

	"example.com/fanin/internal/cat"
)

// TestJoinNoWriteBackToInputs is an independent regression test for the
// aliasing bug where Join/concatenation could share a backing array with one
// of the input slices, so writing to the result polluted the caller's slice.
func TestJoinNoWriteBackToInputs(t *testing.T) {
	// a deliberately has spare capacity so a naive append(a, b...) would reuse
	// a's backing array.
	a := make([]int, 2, 8)
	a[0], a[1] = 1, 2
	aSnapshot := append([]int(nil), a...)

	b := []int{3, 4}
	bSnapshot := append([]int(nil), b...)

	joined := cat.Join(a, b)

	// Mutating the result must not touch either input slice.
	joined[0] = 90
	joined[1] = 91
	joined[2] = 92
	joined[3] = 93

	if !reflect.DeepEqual(a, aSnapshot) {
		t.Fatalf("Join result aliased a: got %v, want %v", a, aSnapshot)
	}
	if !reflect.DeepEqual(b, bSnapshot) {
		t.Fatalf("Join result aliased b: got %v, want %v", b, bSnapshot)
	}

	// Result content (after our mutation) should be independent and correct in
	// length only here; content correctness is checked in TestJoinContent.
	if len(joined) != 4 {
		t.Fatalf("len(joined) = %d, want 4", len(joined))
	}
}

// TestJoinContent checks the concatenation result and several edge cases.
func TestJoinContent(t *testing.T) {
	cases := []struct {
		name string
		a, b []int
		want []int
	}{
		{"both non-empty", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"empty a", nil, []int{3, 4}, []int{3, 4}},
		{"empty b", []int{1, 2}, nil, []int{1, 2}},
		{"both empty", nil, nil, []int{}},
		{"single elements", []int{1}, []int{2}, []int{1, 2}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cat.Join(tc.a, tc.b)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Join(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// TestJoinFreshCapacity ensures the result's capacity is exactly its length,
// i.e. no spare capacity is carried over from a's backing array.
func TestJoinFreshCapacity(t *testing.T) {
	a := make([]int, 2, 8)
	a[0], a[1] = 1, 2
	joined := cat.Join(a, []int{3, 4})
	if cap(joined) != len(joined) {
		t.Fatalf("cap(joined) = %d, want %d (no spare capacity)", cap(joined), len(joined))
	}
}
