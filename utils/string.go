package utils

import "fmt"

func ToString(str string, value any) string {
	strVal := fmt.Sprintf("%s%s", str, value)
	return strVal
}
