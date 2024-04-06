package middleware

func MinInt(leftInt, rightInt int) int {
	if leftInt < rightInt {
		return leftInt
	}
	return rightInt
}

/*
This function is used to check if checkout is possible using the
left point and right point in a tracking center and returns the
number of possible times to checkout

NCheckoutPossible --> Number of Times checkout is possible
*/
func NCheckoutPossible(leftPoint int, rightPoint int, checkoutValue int) int {

	leftInt := leftPoint / checkoutValue
	rightInt := rightPoint / checkoutValue

	min := MinInt(leftInt, rightInt)

	return min
}
