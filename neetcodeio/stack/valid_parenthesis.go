package stack

func isValid(s string) bool {
	bracketAndCo := make([]rune, 0)

	// what is c? byte or rune?
	for _, c := range s {
		if isOpenBracketAndCo(c) {
			bracketAndCo = append(bracketAndCo, c)
		}

		length := len(bracketAndCo) - 1

		if c == ']' {
			if length >= 0 && bracketAndCo[length] == '[' {
				bracketAndCo = bracketAndCo[:length]
			} else {
				return false
			}
		}
		if c == '}' {
			if length >= 0 && bracketAndCo[length] == '{' {
				bracketAndCo = bracketAndCo[:length]
			} else {
				return false
			}
		}
		if c == ')' {
			if length >= 0 && bracketAndCo[length] == '(' {
				bracketAndCo = bracketAndCo[:length]
			} else {
				return false
			}
		}
	}

	return len(bracketAndCo) == 0
}

func isOpenBracketAndCo(c rune) bool {
	return c == '[' || c == '{' || c == '('
}
