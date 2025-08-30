package internal

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type PathType string

func (p PathType) FromSlash() string {
	return filepath.FromSlash(string(p))
}

type Builder struct {
	path      string
	Version   string
	NameSpace string
	Import    []struct {
		Path PathType
		Name string
	} `yaml:",omitempty"`
	DefaultCommand string `yaml:"default_command,omitempty"`
	Command        map[string]struct {
		Desc  string
		Shell []string
	} `yaml:",inline"`
}

func (b *Builder) GenerateFile(name string, force ...bool) error {
	_, err := os.Stat(name)
	if !os.IsNotExist(err) && (len(force) < 0 || !force[0]) {
		return errors.New(name + " file already exists")
	}
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(b)
}

func OpenBuilder(name string) (*Builder, error) {
	name = filepath.FromSlash(name)
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open builder file %q: %v", name, err)
	}
	defer f.Close()
	var builder Builder
	if err = yaml.NewDecoder(f).Decode(&builder); err != nil {
		return nil, fmt.Errorf("builder file invalid %q: %v", name, err)
	}
	builder.path = name
	return &builder, err
}

func DefaultBuilder(out string, force bool) error {
	b := &Builder{
		Version:   "1.0",
		NameSpace: "default",
		Import: []struct {
			Path PathType
			Name string
		}{
			{
				Path: "{{file path}}",
				Name: "{{reNamespace}}",
			},
		},
		DefaultCommand: "hello",
		Command: map[string]struct {
			Desc  string
			Shell []string
		}{
			"hello": {
				Desc:  "hello command",
				Shell: []string{`echo "hello world"`},
			},
		},
	}
	return b.GenerateFile(out, force)
}
