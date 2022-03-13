package properties

import (
	"bufio"
	"fmt"
	"github.com/ushiosan23/go-properties/validator/pairs"
)

// ############################################################
// Interfaces
// ############################################################

// Pair Element that stores information represented by key value.
// The information in this item is read-only
type Pair interface {

	// Key Get key name
	Key() string

	// Value Get pair value
	Value() interface{}
}

// MutablePair Element that stores information represented by key value.
// The information in this item can be changed
type MutablePair interface {

	// Pair Inherit interface
	Pair

	// SetValue Change element value
	SetValue(data interface{}) interface{}
}

// Properties Element that saves information in key value format
// and that can be saved in a props file.
type Properties interface {

	// PutResolver Special method used to extend the functionality of the object.
	// Behind the scenes this method implements an extension that looks for element patterns
	// and transforms those patterns into solid values. You can implement a plugin which detects
	// ${[A-Z]} patterns and replaces them with environment variables.
	PutResolver(resolver func(data string) string)

	// GetProperty Returns the indicated property or an error if the property does not exist.
	GetProperty(key string) (string, error)

	// GetPropertyOrDefault Returns the indicated property or a defVal if the property does not exist.
	// This function does not emit any error.
	GetPropertyOrDefault(key string, defVal string) string

	// Count Returns the size of the items stored inside this object.
	Count() int

	// IsEmpty Check if the current object is empty or if it has any meaningful data
	IsEmpty() bool

	// Keys Returns all the keys of the elements
	Keys() []string

	// Values Return all the values of the elements
	Values() []string

	// Pairs Return all the element pairs
	Pairs() []Pair

	// Contains Checks if any of the keys match the given value.
	Contains(key string) bool

	// Put Insert a new element to the object
	Put(key string, value interface{}) string

	// PutAll Insert multiple elements to the object
	PutAll(elements []Pair) error

	// Remove Delete an item from the save data
	Remove(key string) string

	// Clear Delete all saved items
	Clear()
}

// FProperties Element that determines that the content used comes from an external source.
type FProperties interface {

	// Properties Inherit interface
	Properties

	// Load Loads content from an external source.
	// The function is responsible for determining if the content is valid or not.
	Load(reader *bufio.Reader) error

	// Store Saves the information to a source external to the object.
	Store(writer *bufio.Writer) error
}

// ############################################################
// Structures
// ############################################################

type keyPair struct {
	key   string
	value interface{}
}

// ############################################################
// Generators
// ############################################################

// MutablePairOf Generates a new mutable pair element.
//
// A key must be supplied for the pair to be correct,
// and the value can be anything including the <nil> value.
func MutablePairOf(key string, value interface{}) MutablePair {
	err := pairs.ValidateKey(key)
	// Check error
	if err != nil {
		panic(err)
	}
	// Returns a new pair
	return &keyPair{key: key, value: value}
}

// PairOf Generates a new pair element.
// A key must be supplied for the pair to be correct,
// and the value can be anything including the 'nil' value.
func PairOf(key string, value interface{}) Pair {
	return MutablePairOf(key, value)
}

// ############################################################
// Methods
// ############################################################

// Key Get key name
func (k keyPair) Key() string {
	return k.key
}

// Value Get key value. This element can be nil
func (k keyPair) Value() interface{} {
	return k.value
}

// SetValue Change element value. Can be nil
func (k *keyPair) SetValue(data interface{}) interface{} {
	old := k.value
	// Change value
	k.value = data
	// Return old value
	return old
}

// String Structure string representation
func (k keyPair) String() string {
	return fmt.Sprintf("[%s = %v]", k.key, k.value)
}
