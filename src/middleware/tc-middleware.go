package middleware

func MinInt(leftInt, rightInt int) int {
	if leftInt < rightInt {
		return leftInt
	}
	return rightInt
}
