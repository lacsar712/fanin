package cat_test

import (
	"testing"
	"example.com/fanin/internal/cat"
)

func TestJoinNoAlias(t *testing.T) {
	base := make([]int, 2, 8)
	base[0], base[1] = 1, 2
	alias := base[:2]
	joined := cat.Join(alias, []int{3, 4})
	joined[0] = 9
	if alias[0] == 9 {
		t.Fatal("join aliased caller slice")
	}
	if len(joined) != 4 || joined[1] != 2 || joined[2] != 3 {
		t.Fatalf("got %v", joined)
	}
}
