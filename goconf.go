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

// getData gets value by served key, if any error, return default value with
// the error.
func (parser *JSONParser) getData(key string, defaultVal interface{}) (
	interface{}, error) {

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
				return defaultVal, NoSuchKeyError{strings.Join(keyPath[:(i+1)], "/")}
			}
			value = val
		} else {
			return defaultVal, NoSuchKeyError{strings.Join(keyPath[:(i+1)], "/")}
		}
	}

	return value, nil
}

// GetString gets string value by served key, the value will be converted from
// interface{} to string type if possible.
//
// If the key is not found, return defaultVal passed in with NoSuchKeyError. If
// the value type is not string, return defaultVal passed in with
// InvalidValueTypeError.
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

// GetBool gets boolean value by served key, the value will be converted from
// interface{} to bool type if possible.
//
// If the key is not found, return defaultVal passed in with NoSuchKeyError. If
// the value type is not bool, return defaultVal passed in with
// InvalidValueTypeError.
func (parser *JSONParser) GetBool(key string, defaultVal bool) (bool, error) {
	// Get value.
	value, err := parser.getData(key, defaultVal)
	if err != nil {
		return defaultVal, err
	}

	// Convert to bool.
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, InvalidValueTypeError{key, value}
}

// GetFloat gets float64 value by served key, the value will be converted from
// interface{} to float64 type if possible.
//
// If the key is not found, return defaultVal passed in with NoSuchKeyError. If
// the value type is not float64, return defaultVal passed in with
// InvalidValueTypeError.
func (parser *JSONParser) GetFloat64(key string, defaultVal float64) (float64, error) {
	// Get value.
	value, err := parser.getData(key, defaultVal)
	if err != nil {
		return defaultVal, err
	}

	// Convert to float64.
	if f, ok := value.(float64); ok {
		return f, nil
	}
	return defaultVal, InvalidValueTypeError{key, value}
}

// GetInt gets int value by served key, the value will be converted from
// interface{} to int type if possible.
//
// If the key is not found, return defaultVal passed in with NoSuchKeyError. If
// the value type is not int, return defaultVal passed in with
// InvalidValueTypeError.
func (parser *JSONParser) GetInt(key string, defaultVal int) (int, error) {
	// Get value. json unmarshall parse numerical data to float64 by default.
	value, err := parser.GetFloat64(key, float64(defaultVal))
	if err != nil {
		return defaultVal, err
	}

	// Convert to int.
	i := int(value)
	if float64(i) == value {
		return i, nil
	}
	return defaultVal, InvalidValueTypeError{key, value}
}
