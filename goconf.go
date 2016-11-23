// Generic JSON parser.
// TODO:
//  - Add support for other value type.
//  - Add support for list (slice) as type of value.

package goconf

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/golang/glog"
)

// JSONParser struct contains data read from any JSON source. With this struct,
// one can query value by key. It only supports string as type for the keys. For
// values, they can be any type but not slice.
type JSONParser struct {
	data interface{}
}

// NewJSONParserFromFile creates JSONParser object out of a JSON file. It reads
// the JSON file from the input file path, and store the data into the newly
// created JSONParser.
//
// The JSONParser only support keys in string type and values are generic
// interface{} type.
func NewJSONParserFromFile(filePath string) (*JSONParser, error) {
	// Read file to byte array.
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return NewJSONParserFromBytes(file)
}

// NewJSONParserFromBytes creates JSONParser object out of a byte array. It
// reads the byte array then store the data into the newly craeted JSONParser.
//
// The JSONParser only support keys in string type and values are generic
// interface{} type.
func NewJSONParserFromBytes(file []byte) (*JSONParser, error) {
	var data interface{}

	// Transform the file byte array to map.
	err := json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return &JSONParser{data}, nil
}

func (parser *JSONParser) getData(key string, defaultVal interface{}) (
	interface{}, error) {
	// Handle nil default value.
	if defaultVal == nil {
		defaultVal = ""
	}

	// Empty key.
	if len(key) == 0 {
		glog.Info("Empty key, returning default value parsed in.")
		return defaultVal, nil
	}

	// Parse the key to nested keys.
	keyPath := strings.Split(key, "/")
	var value interface{}
	value = parser.data

	// Parse the map to get the value.
	for i, s := range keyPath {
		if m, ok := value.(map[string]interface{}); ok {
			val, exists := m[s]
			if !exists {
				return nil, NoSuchKeyError{strings.Join(keyPath[:(i+1)], "/")}
			}
			value = val
		} else {
			return nil, InvalidValueTypeError{strings.Join(keyPath[:(i+1)], "/"), value}
		}
	}

	return value, nil
}

// GetString gets string value by key, the value will be converted from
// interface{} to string type if possible.
func (parser *JSONParser) GetString(key string, defaultVal string) (string, error) {
	// Get value.
	value, err := parser.getData(key, defaultVal)
	if err != nil {
		return defaultVal, err
	}

	// Convert to string.
	if str, ok := value.(string); ok {
		return str, nil
	}
	return "", InvalidValueTypeError{key, value}
}

// GetBool gets boolean value by key, the value will be converted to bool type
// if possible.
func (parser *JSONParser) GetBool(key string, defaultVal bool) (bool, error) {
	// Get value.
	value, err := parser.getData(key, defaultVal)
	if err != nil {
		return defaultVal, err
	}

	// Convert to string.
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, InvalidValueTypeError{key, value}
}
