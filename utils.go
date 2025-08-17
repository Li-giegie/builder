package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Must(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}

func printExit(code int, args ...string) {
	buf := make([]byte, 0, len(args)*5)
	for _, arg := range args {
		buf = append(buf, arg...)
		buf = append(buf, ' ')
	}
	buf = append(buf, '\n')
	os.Stderr.Write(buf)
	os.Exit(code)
}

func Execute(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func ScanWorld(str string) []string {
	var builder bytes.Buffer
	var item = make([]string, 0, 16)
	var tag1, tag2 bool
	r := strings.NewReader(str)
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if builder.Len() > 0 {
					item = append(item, builder.String())
				}
				builder.Reset()
				break
			}
			return nil
		}
		switch char {
		case '"':
			if !tag2 {
				tag1 = !tag1
			}
			builder.WriteRune(char)
		case '\'':
			if !tag1 {
				tag2 = !tag2
			}
			builder.WriteRune(char)
		case ' ':
			if tag1 || tag2 {
				builder.WriteRune(char)
			} else {
				item = append(item, builder.String())
				builder.Reset()
			}
		default:
			builder.WriteRune(char)
		}
	}
	var result = make([]string, 0, len(item))
	for _, s := range item {
		if len(s) == 0 {
			continue
		}
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			s = s[1 : len(s)-1]
		} else if s[0] == '"' && s[len(s)-1] == '"' {
			s = s[1 : len(s)-1]
		}
		result = append(result, s)
	}
	return result
}
