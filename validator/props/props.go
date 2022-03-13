package props

import (
	"errors"
	"properties/utils"
	"strings"
)

var (
	commentElements  = []rune{'#', '!'}
	continuationRune = '\\'
)

// ############################################################
// Getters
// ############################################################

// CommentElements Returns the elements considered comments within the content
func CommentElements() []rune {
	return commentElements
}

// ContinuationRune Returns the rune that determines if the content continues on the next line.
func ContinuationRune() rune {
	return continuationRune
}

// ############################################################
// Validators
// ############################################################

// IsLineFullComment Determine if passed data is a comment line
func IsLineFullComment(data string) bool {
	data = strings.TrimSpace(data)
	elements := commentStrElements()
	// Evaluate
	return strings.HasPrefix(data, elements[0]) || strings.HasPrefix(data, elements[1])
}

// IsLineValid Determine if passed data is a valid property line data
func IsLineValid(data string) bool {
	data = strings.TrimSpace(data)
	// Check empty lines
	if len(data) == 0 {
		return false
	}
	// Check if line is a comment
	return !IsLineFullComment(data)
}

// ############################################################
// Utilities
// ############################################################

// CommentIndex Returns the last position of a comment character.
func CommentIndex(line string) int {
	return utils.AnyStrIndexOf(line, commentElements...)
}

// ContinuationIndex Returns the position of the backslash to indicate
// that the content continues on the next line.
func ContinuationIndex(line string) int {
	return utils.StrIndexOf(line, continuationRune)
}

// ############################################################
// Cleaners
// ############################################################

// LineWithoutComment Returns only valid content within a line.
//
// Content recognized as a comment is ignored.
func LineWithoutComment(line string) string {
	line = strings.TrimSpace(line)
	// Check if contains any comment character
	index := CommentIndex(line)
	// Return clean line
	if index == utils.IndexNotFound {
		return line
	}
	return line[:index]
}

// ValidContent Returns only the valid part of the information.
//
// The backslash is ignored.
func ValidContent(data string) string {
	data = strings.TrimSpace(data)
	// Check if contains backslash
	index := ContinuationIndex(data)
	// Return valid data
	if index == utils.IndexNotFound {
		return data
	}
	return data[:index]
}

// GenerateLinePair Get line pair element
func GenerateLinePair(data string) (string, interface{}, error) {
	data = strings.TrimSpace(data)
	// Get index pair
	index := utils.StrIndexOf(data, '=')
	// Check if data is not valid
	if index == utils.IndexNotFound {
		return "", nil, errors.New("invalid data pair")
	}
	// Generate result
	key := strings.TrimSpace(data[:index])
	value := RealValue(strings.TrimSpace(data[index+1:]))

	// Check key size
	if len(key) == 0 {
		return "", nil, errors.New("key cannot be empty")
	}
	// Return result
	return key, value, nil
}

// RealValue Returns a real pair value
func RealValue(data string) interface{} {
	if len(data) == 0 {
		return nil
	}
	return data
}

// ############################################################
// Internal functions
// ############################################################

// commentStrElements Returns the elements considered comments within the content
func commentStrElements() []string {
	result := make([]string, len(commentElements))
	for i, v := range commentElements {
		result[i] = string(v)
	}
	return result
}
