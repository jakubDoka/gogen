package test2

//def(
//rules Max<int>

func Max(a, b int) int {
	if a > b {
		return b
	}
	return a
}

//)

//def(
//rules Min<int>

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//)

//def(
//rules Clamp<int>
//dep Max<int, Max>
//dep Min<int, Min>

func Clamp(val, min, max int) int {
	return Max(min, Min(max, val))
}

//)
