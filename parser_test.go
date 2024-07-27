package dotenv

import (
	"fmt"
	"os"
	"testing"
)

type Config struct {
	ApiPort  string
	User     string
	Password string
	Host     string
	Port     int
	Db       string `env:"db_name"`
}

func TestParser(t *testing.T) {
	var conf Config

	buffer, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}
	if err := Unmarshal(buffer, &conf); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("conf: %v\n", conf)

	slice, err := UnmarshalSlice(buffer)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range slice {
		fmt.Println(v.Key, "=>", v.Value)
	}
	value, ok := slice.Get("user")
	if ok {
		fmt.Printf("value: %v\n", value)
	}

	_, ok = slice.Get("other")
	if ok {
		t.Fatal("inappropriate value")
	}
}
