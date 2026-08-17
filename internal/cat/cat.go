package cat

// Join returns a new slice holding the elements of a followed by the elements of b.
//
// The returned slice never aliases either input's backing array, so mutating it
// cannot pollute the caller's slices. This holds even when cap(a) is large
// enough that a naive append(a, b...) would reuse a's backing array.
func Join(a, b []int) []int {
	out := make([]int, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	return out
}
