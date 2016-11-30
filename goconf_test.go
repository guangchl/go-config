// Unit testing for Generic JSON parser.

package goconf

import (
	"flag"
	"reflect"
	"testing"

	"github.com/golang/glog"
)

var data []byte

func init() {
	flag.Parse()
	data = []byte(`
		{
			"object": {
				"Database": {
					"host": "localhost",
					"name": "go",
					"pass": true,
					"type": "mysql",
					"user": "root"
				},
				"buffer_size": 10
			}
		}`)
}

func TestGetString(t *testing.T) {
	glog.Info("***Testing get string from JSON.***")

	parser, err := NewJSONParserFromBytes(data)
	if err != nil {
		glog.Fatal(err)
	}

	str, err := parser.GetString("object/Database/name", "defaultName")
	if err != nil {
		glog.Fatal(err)
	}

	if str != "go" {
		glog.Fatalf("Assertion failure, expected: go, result: %s.", str)
	}
}

func TestGetBool(t *testing.T) {
	glog.Info("***Testing get boolean from JSON.***")

	parser, err := NewJSONParserFromBytes(data)
	if err != nil {
		glog.Fatal(err)
	}

	b, err := parser.GetBool("object/Database/pass", false)
	if err != nil {
		glog.Fatal(err)
	}

	if !b {
		glog.Fatal("Assertion failure, expected: true.")
	}
}

func TestGetInt(t *testing.T) {
	glog.Info("***Testing get boolean from JSON.***")

	parser, err := NewJSONParserFromBytes(data)
	if err != nil {
		glog.Fatal(err)
	}

	i, err := parser.GetInt("object/buffer_size", -1)
	if err != nil {
		glog.Fatal(err)
	}

	if i != 10 {
		glog.Fatal("Assertion failure, expected: 10, result: %s.", i)
	}
}

func TestInvalidKey(t *testing.T) {
	glog.Info("***Testing get invalid key from JSON.***")

	parser, err := NewJSONParserFromBytes(data)
	if err != nil {
		glog.Fatal(err)
	}

	_, err = parser.GetBool("wrongkey", false)
	if err == nil || reflect.TypeOf(err).Name() != "NoSuchKeyError" {
		glog.Fatalf("Assertion failure, expected: NoSuchKeyError, result: %s.", reflect.TypeOf(err))
	}
}

func TestInvalidValueType(t *testing.T) {
	glog.Info("***Testing invalid value type from JSON.***")

	parser, err := NewJSONParserFromBytes(data)
	if err != nil {
		glog.Fatal(err)
	}

	_, err = parser.GetString("object/Database/pass", "false")
	if err == nil || reflect.TypeOf(err).Name() != "InvalidValueTypeError" {
		glog.Fatalf("Assertion failure, expected: InvalidValueTypeError, result: %s.", reflect.TypeOf(err))
	}
}
