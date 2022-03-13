package properties

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"properties/convert"
	"properties/utils"
	"properties/validator/props"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

// ############################################################
// Constants
// ############################################################

const commentDateLayout = "Mon Jan 02 15:04:05 MST 2006"

// ############################################################
// Interfaces
// ############################################################

type PropertyPair interface {
	MutablePair
	ValueStr() string
}

// ############################################################
// Structures
// ############################################################

type fileProperties struct {
	file      *os.File
	elements  map[string]string
	mutex     *sync.Mutex
	resolvers []func(data string) string
}

type loadState struct {
	buffering    bool
	lastProperty string
}

// ############################################################
// Generators
// ############################################################

func NewProperties() FProperties {
	return &fileProperties{
		file:     nil,
		elements: make(map[string]string),
		mutex:    new(sync.Mutex),
	}
}

func FromPair(pair Pair) PropertyPair {
	switch out := pair.(type) {
	case PropertyPair:
		return out
	default:
		return &keyPair{key: out.Key(), value: out.Value()}
	}
}

// ############################################################
// Implementations
// ############################################################

func (k *keyPair) ValueStr() string {
	return convert.RawValueString(k.Value())
}

// ############################################################
// Methods
// ############################################################

func (f *fileProperties) PutResolver(resolver func(data string) string) {
	// Check if already exists
	for _, v := range f.resolvers {
		// Storage functions identifier
		id1 := fmt.Sprintf("%v", reflect.ValueOf(v))
		id2 := fmt.Sprintf("%v", reflect.ValueOf(resolver))
		// Check identifiers
		if id1 == id2 {
			return
		}
	}
	// Add resolver
	f.resolvers = append(f.resolvers, resolver)
}

func (f *fileProperties) Load(reader *bufio.Reader) error {
	state := new(loadState)
	return f.load(reader, state)
}

func (f *fileProperties) Store(writer *bufio.Writer) error {
	return f.store(writer)
}

func (f *fileProperties) GetProperty(key string) (string, error) {
	if f.Contains(key) {
		return f.processValue(f.elements[key]), nil
	}
	return "", errors.New(fmt.Sprintf("property \"%s\" not found", key))
}

func (f *fileProperties) GetPropertyOrDefault(key string, defVal string) string {
	data, err := f.GetProperty(key)
	// Check error
	if err != nil {
		return f.processValue(defVal)
	}
	return data
}

func (f *fileProperties) Count() int {
	return len(f.elements)
}

func (f *fileProperties) IsEmpty() bool {
	return f.Count() == 0
}

func (f *fileProperties) Keys() []string {
	// Generate key array
	keys := make([]string, 0, f.Count())
	// Iterate all elements
	for key := range f.elements {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (f *fileProperties) Values() []string {
	// Generate value array
	values := make([]string, 0, f.Count())
	// Iterate elements
	for _, value := range f.elements {
		values = append(values, value)
	}
	return values
}

func (f *fileProperties) Pairs() []Pair {
	// Generate pair data
	pairs := make([]Pair, 0, f.Count())
	// Iterate elements
	for k, v := range f.elements {
		pairs = append(pairs, PairOf(k, v))
	}
	return pairs
}

func (f *fileProperties) Contains(key string) bool {
	_, ok := f.elements[key]
	return ok
}

func (f *fileProperties) Put(key string, value interface{}) string {
	// Lock mutex
	f.mutex.Lock()
	// Storage old value
	old := f.GetPropertyOrDefault(key, "")
	f.elements[key] = convert.RawValueString(value)
	// Unlock mutex
	defer f.mutex.Unlock()
	return old
}

func (f *fileProperties) PutAll(elements []Pair) error {
	// Lock mutex
	f.mutex.Lock()
	// Put all elements
	for _, item := range elements {
		f.elements[item.Key()] = convert.RawValueString(item.Value())
	}
	// Unlock mutex
	defer f.mutex.Unlock()
	return nil
}

func (f *fileProperties) Remove(key string) string {
	// Lock mutex
	f.mutex.Lock()
	// Store old value
	old := f.GetPropertyOrDefault(key, "")
	if f.Contains(key) {
		delete(f.elements, key)
	}
	// Unlock mutex
	defer f.mutex.Unlock()
	return old
}

func (f *fileProperties) Clear() {
	// Lock mutex
	f.mutex.Lock()
	// Clear map
	f.elements = make(map[string]string)
	// Unlock mutex
	defer f.mutex.Unlock()
}

func (f fileProperties) String() string {
	// Store temporal data
	keys := f.Keys()
	// Check if properties are empty
	if f.IsEmpty() {
		return "(0) {}"
	}
	// Generators
	result := new(strings.Builder)
	result.WriteString(fmt.Sprintf("(%d) {", f.Count()))
	// Iterate elements
	for i, v := range keys {
		data := fmt.Sprintf("[%s = %s]", v, f.elements[v])
		if i == len(keys)-1 {
			result.WriteString(data + "}")
		} else {
			result.WriteString(data + ", ")
		}
	}
	return result.String()
}

// ############################################################
// Internal methods
// ############################################################

func (f *fileProperties) load(r *bufio.Reader, state *loadState) error {
	// Generate result
	var result error = nil
	// Read content
	for {
		// Read line
		line, err := r.ReadString('\n')
		// Check if err is not null
		if err != nil {
			if err != io.EOF {
				result = err
			} else {
				// EOF almost contains line content
				f.processLine(line, state)
			}
			break
		}
		// Ignore all line comments
		f.processLine(line, state)
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func (f *fileProperties) store(w *bufio.Writer) error {
	// Get current date
	now := time.Now().Format(commentDateLayout)
	// Write start comment
	w.WriteString("#" + now + utils.LineSeparator)
	for _, item := range f.Keys() {
		key := item
		value := f.elements[key]
		// Store value
		w.WriteString(key + "=" + value + utils.LineSeparator)
	}
	// Flush writer
	out := w.Flush()
	return out
}

func (f *fileProperties) processLine(line string, state *loadState) {
	// Ignore all invalid lines
	if !props.IsLineValid(line) {
		return
	}
	// Process line data
	line = props.LineWithoutComment(line)
	// Check if state is buffering
	if state.buffering {
		// Check if value contains a multiline data
		indexC := props.ContinuationIndex(line)
		if indexC == utils.IndexNotFound {
			state.buffering = false
			indexC = len(line)
		}
		// Append information
		old := f.GetPropertyOrDefault(state.lastProperty, "")
		newVal := old + props.ValidContent(line)
		f.Put(state.lastProperty, newVal)
	} else {
		key, value, err := props.GenerateLinePair(line)
		// Check any error
		if err != nil {
			return
		}
		// Check if value contains a multiline data
		pair := FromPair(PairOf(key, value))
		indexC := props.ContinuationIndex(pair.ValueStr())
		if indexC != utils.IndexNotFound {
			state.lastProperty = pair.Key()
			state.buffering = true
		}
		// Insert data
		f.Put(pair.Key(), props.ValidContent(pair.ValueStr()))
	}
}

func (f *fileProperties) processValue(data string) string {
	for _, ff := range f.resolvers {
		data = ff(data)
	}
	return data
}
