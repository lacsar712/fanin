package cat

// Join appends b onto a without copying when cap(a) is large enough.
func Join(a, b []int) []int {
	return append(a, b...)
}
