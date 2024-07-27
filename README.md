# dotenv

A Go language parsing tool for key value pair configuration files.

Used in configuration file parsing of back-end small applications. Only the following features are supported, please consider using them.

### grammar

This tool parses by line, where one line can be a key value, a comment, or a blank line.

```
_id=20
name = "Kate"
# Comment
```

- Keys can be uppercase and lowercase English letters and underscores.
- The equal sign cannot be missing.
- Values are parsed as strings by default, and single and double quotes are automatically ignored.

### value types

Structure unmarshalling of the following value types is supported.

```go
string
int     | int8    | int16  | int32  | int64
uint    | uint8   | uint16 | uint32 | uint64
float32 | float64
```

### install

```
go get github.com/bingxio/dotenv
```

### example

See: [parser_test.go](https://github.com/bingxio/dotenv/blob/main/parser_test.go)
