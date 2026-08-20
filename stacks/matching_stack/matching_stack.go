package matchingstack

func areBracketsProperlyMatched(code_snippet string) bool {
	// Write your code here
	// () [] {} -> appar in any order like func main() {} or [{a, b}, {c, d}] or func([()])
	// the closing must match the last opened.
	// func main() {} -> [{}] -> [, {, }
	// stack
	// state -> all last opened elements [, { and (

	stack := make([]rune, 0, len(code_snippet))

	var closingElements = map[rune]rune{
		']': '[',
		'}': '{',
		')': '(',
	}

	// now edge cases? [ -> invalid
	// ([[[[[{} ]]]]])
	// now edge cases? empty...
	// vai guardar na staick.. e depois mais nada
	for _, c := range code_snippet {
		// store this
		// append adds at the end -> need to pop from last to first
		// will try with slice operations
		// [{}] -> [, {...
		if c == '[' || c == '{' || c == '(' {
			stack = append(stack, c)
		}

		v, ok := closingElements[c]

		// opa, temos um }..
		if ok {
			if len(stack) == 0 {
				return false
			}
			// pega ai o ultimo add mano -> }
			// ei, } me da {.. sao iguais? YES. tira o ultimo, segue o rumo no proximo
			// ei, ]! me da [.. sao iguais?
			tail := stack[len(stack)-1]
			if tail == v {
				// all good, the closing matches the last openend
				// if good, need to pop htat
				stack = stack[:len(stack)-1]
			} else {
				// nop, bad, should return false
				return false
			}
		}
	}

	return len(stack) == 0
}
