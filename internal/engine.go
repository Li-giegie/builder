package internal

import (
	"errors"
	"fmt"
	"github.com/Li-giegie/builder/pkg"
	"strings"
	"sync"
)

type Engine struct {
	Root  *Builder
	cache sync.Map
}

func NewEngine(name string) (*Engine, error) {
	cfg, err := OpenBuilder(name)
	if err != nil {
		return nil, err
	}
	engine := &Engine{Root: cfg}
	engine.cache.Store(cfg.path, cfg)
	return engine, nil
}

func (e *Engine) Execute(commands []string) error {
	if len(commands) == 0 {
		if len(e.Root.DefaultCommand) == 0 {
			return errors.New("empty command")
		}
		commands = []string{e.Root.DefaultCommand}
	}
	for _, command := range commands {
		cmd, ok := e.Root.Command[command]
		if !ok {
			return fmt.Errorf("not found command %q", command)
		}
		for _, s := range cmd.Shell {
			if s == "" {
				continue
			}
			var execCmds []string
			// 引用关系
			if s[0] == '$' {
				result, err := e.ParseRef(e.Root, s, nil)
				if err != nil {
					return err
				}
				execCmds = append(execCmds, result...)
			} else {
				execCmds = append(execCmds, s)
			}
			for _, item := range execCmds {
				execCmd := pkg.ScanWorld(item)
				err := pkg.Execute(execCmd[0], execCmd[1:])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (e *Engine) ParseRef(root *Builder, s string, refPath map[string]map[string]struct{}) ([]string, error) {
	args := strings.Split(s, ".")
	if len(args) != 2 || args[0] == "" || args[1] == "" {
		return nil, fmt.Errorf("syntax error [namespace.command] err: %q", s)
	}
	var (
		namespace = strings.TrimPrefix(args[0], "$")
		command   = args[1]
		path      string
	)
	for _, s2 := range root.Import {
		if s2.Name == namespace {
			path = s2.Path.FromSlash()
		}
	}
	if path == "" {
		return nil, fmt.Errorf("import %q not found", namespace)
	}
	if _, ok := refPath[root.path][path]; ok {
		return nil, fmt.Errorf("imported packages appear with circular references \n%q\n\t%q", root.path, path)
	}
	val, ok := e.cache.Load(path)
	if !ok {
		cfg, err := OpenBuilder(path)
		if err != nil {
			return nil, err
		}
		e.cache.Store(path, cfg)
		val = cfg
	}
	if refPath == nil {
		refPath = make(map[string]map[string]struct{})
	}
	refPathItem := refPath[root.path]
	if refPathItem == nil {
		refPathItem = make(map[string]struct{})
		refPath[root.path] = refPathItem
	}
	refPathItem[path] = struct{}{}

	builder := val.(*Builder)
	cmd, ok := builder.Command[command]
	if !ok {
		return nil, errors.New("namespace " + namespace + " command" + command + " not found command")
	}
	var result []string
	for _, s2 := range cmd.Shell {
		if s2 == "" {
			continue
		}
		if strings.HasPrefix(s2, "$") {
			res, err := e.ParseRef(builder, s2, refPath)
			if err != nil {
				return nil, err
			}
			result = append(result, res...)
		} else {
			result = append(result, s2)
		}
	}
	return result, nil
}
