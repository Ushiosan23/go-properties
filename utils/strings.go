package utils

// IndexNotFound Constant used to determine if an element has not been found within an array.
const IndexNotFound = -1

// StrIndexOf Returns the index of the first match within a string.
func StrIndexOf(data string, item rune) int {
	// Generate result
	result := IndexNotFound
	// Iterate data
	for i, v := range data {
		if v == item {
			result = i
			break
		}
	}
	return result
}

// AnyStrIndexOf returns the index of the first match found within a string.
//
// The difference is that in this function you can use more than one rune to check the result.
// If at least one rune is found the function ends and returns the position of the match.
func AnyStrIndexOf(data string, items ...rune) int {
	// Generate result
	result := IndexNotFound
	fb := false
	// Iterate data
	for i, v := range data {
		for _, iv := range items {
			if v == iv {
				result = i
				fb = true
				break
			}
		}
		// Check child break
		if fb {
			break
		}
	}
	return result
}
