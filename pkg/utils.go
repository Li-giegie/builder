package pkg

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Must(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}

func Execute(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
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

func PrintObj(a any) {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err)
	}
	println(string(data))
}

func OpenYamlObj(name string, a any) error {
	f, err := os.Open(filepath.FromSlash(name))
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(a)
}

func SaveYamlObj(name string, a any) error {
	f, err := os.OpenFile(filepath.FromSlash(name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(a)
}

func MkDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, os.ModePerm)
		}
		return err
	}
	return nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DefaultStr(s ...string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != "" {
			return s[i]
		}
	}
	return ""
}
