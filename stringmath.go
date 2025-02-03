package gvgo

// Plus pluses source and extra.
// If one of them are not number or number with invalid format, it will return source.
// It ignores higher digits, so "001" + "002" = "003".
func Plus(source, extra string) string {
	if !isNum(source) || !isNum(extra) {
		return source
	}
	if isBadNum(source) || isBadNum(extra) {
		return source
	}

	res := ""
	carry := 0
	i := len(source) - 1
	j := len(extra) - 1

	for i >= 0 || j >= 0 || carry > 0 {
		digit1 := 0
		if i >= 0 {
			digit1 = int(source[i] - '0')
			i--
		}

		digit2 := 0
		if j >= 0 {
			digit2 = int(extra[j] - '0')
			j--
		}

		sum := digit1 + digit2 + carry
		carry = sum / 10
		digit := sum % 10
		res = string(rune(digit+'0')) + res // prepend the digit
	}

	return res
}
