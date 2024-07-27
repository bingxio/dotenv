package dotenv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Kv struct {
	Key   string
	Value string
}

type KvSlice []Kv

func (s KvSlice) Get(key string) (string, bool) {
	for _, v := range s {
		if v.Key == key {
			return v.Value, true
		}
	}
	return "", false
}

var (
	buffer []byte = nil
	line          = 1
	ip            = 0
)

func makeError(text string) error { return fmt.Errorf(text, line) }
func end() bool                   { return ip >= len(buffer) }

func now() byte {
	if ip >= len(buffer) {
		return '\x00'
	}
	return buffer[ip]
}

func isLetter() bool {
	c := now()
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func whiteSpace() {
	for now() == ' ' {
		ip++
	}
	for now() == '\n' {
		ip++
		line++
	}
}

func markComment() {
	if now() == '#' {
		for !end() && now() != '\n' {
			ip++
		}
		ip++
		line++
	}
}

func markKeyValue() (*Kv, error) {
	key := []byte{}

	for isLetter() {
		key = append(key, now())
		ip++
	}
	if len(key) == 0 {
		return nil, makeError("key is empty (line %d)")
	}
	whiteSpace()

	if now() != '=' {
		return nil, makeError("expect equal(=) sign after key at line %d")
	}
	ip++
	if now() == '\n' || end() {
		return nil, makeError("value is empty (line %d)")
	}
	whiteSpace()

	value := []byte{}

	for !end() && now() != '\n' {
		if now() != '"' && now() != '\'' {
			value = append(value, now())
		}
		ip++
	}
	line++

	kv := Kv{
		Key:   string(key),
		Value: string(value),
	}
	return &kv, nil
}

func unmarshal() (KvSlice, error) {
	slice := KvSlice{}

	for !end() {
		if whiteSpace(); ip >= len(buffer) {
			return slice, nil
		}
		if now() == '#' {
			markComment()
			continue
		}
		v, err := markKeyValue()
		if err != nil {
			return nil, err
		}
		slice = append(slice, *v)
		ip++
	}
	return slice, nil
}

func restore(src []byte) {
	buffer = src
	line = 1
	ip = 0
}

func Unmarshal(src []byte, dst any) error {
	restore(src)
	slice, err := unmarshal()

	if err != nil {
		return err
	}
	v := reflect.ValueOf(dst)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("encoding structure type is incorrect")
	}
	v = v.Elem() // value pointed to by the pointer

	for i := 0; i < v.NumField(); i++ {
		fieldType := v.Type().Field(i)
		fieldValue := v.Field(i)

		if !unicode.IsUpper(rune(fieldType.Name[0])) {
			continue
		}
		err := encodeToStruct(slice, fieldType, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnmarshalSlice(src []byte) (KvSlice, error) {
	restore(src)
	return unmarshal()
}

func encodeToStruct(
	slice KvSlice,
	ft reflect.StructField,
	fv reflect.Value,
) error {
	var ok bool
	var value string
	err := fmt.Errorf("undefined key '%s'", ft.Name)

	if ft.Tag == "" {
		value, ok = slice.Get(strings.ToLower(ft.Name))
		if !ok {
			value, ok = slice.Get(strings.ToUpper(ft.Name))
		}
		if !ok {
			return err
		}
	} else {
		value, ok = slice.Get(ft.Tag.Get("env"))
		if !ok {
			return err
		}
	}
	return valueType(value, ft.Type, fv)
}

func valueType(value string, p reflect.Type, val reflect.Value) error {
	switch p.Kind() {
	// string
	case reflect.String:
		val.SetString(value)
		return nil
		// uint
	case reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		v, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return err
		}
		val.SetUint(v)
		return nil
		// int
	case reflect.Int,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		v, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return err
		}
		val.SetInt(v)
		return nil
		// float
	case reflect.Float32, reflect.Float64:

		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		val.SetFloat(v)
		return nil
	}
	return errors.New("only supports parsing number and string")
}
