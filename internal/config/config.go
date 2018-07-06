package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

var Opts struct {
	// Port for website (with ':')
	Port string `default:":500"`
	// Debug mode
	Debug bool `default:"false"`
	// Offset - number of records on single screen
	Offset int `default:"100"`
	// CheckVer - should the program check check is there a new version
	CheckVer bool `default:"true"`
	// IsWriteMode - can program edit databases
	IsWriteMode bool `default:"true"`
	// OpenBrowser - should the program open a browser automatically
	OpenBrowser bool `default:"true"`
	// NeatWindow - should the program open the special neat window
	NeatWindow bool `default:"true"`
}

type field struct {
	name  string
	value interface{}
}

// ParseConfig parses flags like -port, -debug, -offset and etc.
// If there's no any flags, it tries to parse config file "config.ini"
func ParseConfig() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case string:
				err = fmt.Errorf(r.(string))
			case error:
				err = r.(error)
			default:
				err = fmt.Errorf("undefined error")
			}
		}
	}()

	setDefaultValues()

	if len(os.Args) > 1 {
		parseFlags()
	}

	return nil
}

// setDefaultValues sets default values of Opts's fields.
// If tag default was missed it panics
// If type of field isn't [int, string, bool] it panics
func setDefaultValues() {
	var defValues []field
	// For getting tags
	t := reflect.TypeOf(Opts)

	// Opts is always struct, so we shouldn't check if tp.Kind() == reflect.Struct
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		def := f.Tag.Get("default")
		if def == "" {
			panicf("default tag of field %s is empty", t.Field(i).Name)
		}

		switch f.Type.Kind() {
		case reflect.Bool:
			b := func() bool {
				if def == "true" {
					return true
				}
				return false
			}()
			defValues = append(defValues, field{name: f.Name, value: b})
		case reflect.String:
			defValues = append(defValues, field{name: f.Name, value: def})
		case reflect.Int:
			i, _ := strconv.ParseInt(def, 10, 64)
			defValues = append(defValues, field{name: f.Name, value: int(i)})
		default:
			panicf("Bad type of a field of Opts. Type: %s", f.Type.Kind())
		}
	}

	setValues(defValues)
}

func setValues(values []field) {
	opts := reflect.ValueOf(&Opts).Elem()

	for _, v := range values {
		f := opts.FieldByName(v.name)
		if f.IsValid() {
			if f.Kind() != reflect.TypeOf(v.value).Kind() {
				panicf("Different types of field and value: field type - %s, value type - %s", f.Kind().String(), reflect.TypeOf(v.value).Kind().String())
			}

			f.Set(reflect.ValueOf(v.value))
		}
	}
}

func panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v))
}
