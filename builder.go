package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sort"
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
	if err == nil && len(force) > 0 && force[0] {
		return errors.New(name + " file already exists")
	}
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(b)
}

func (b *Builder) DefaultPrint() {
	keys := make([]string, 0, len(b.Command))
	for k := range b.Command {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		println(k, "\t", b.Command[k].Desc)
	}
	println("help\t", "get help")
	println("init\t", "Initialize the builder configuration file")
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

func DefaultBuilder() {
	Must(func() error {
		b := &Builder{
			Version:   "1.0",
			NameSpace: "default",
			Import: []struct {
				Path PathType
				Name string
			}{
				{
					Path: "./a/b",
					Name: "utils",
				},
			},
			DefaultCommand: "build",
			Command: map[string]struct {
				Desc  string
				Shell []string
			}{
				"build": {
					Desc:  "build command",
					Shell: []string{"go build -o main.go"},
				},
			},
		}
		b.DefaultPrint()
		return b.GenerateFile(".builder.yaml")
	})
}
