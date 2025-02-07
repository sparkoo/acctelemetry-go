package acctelemetry_test

func convertToString(chars []uint16) string {
	var str string
	for _, val := range chars {
		str += string(rune(val))
	}
	return str
}
