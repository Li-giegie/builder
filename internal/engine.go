package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/Li-giegie/builder/pkg"
	"strings"
)

type Engine struct {
	Root  *Builder
	cache map[string]*Builder
}

func NewEngine(name string) (*Engine, error) {
	cfg, err := OpenBuilder(name)
	if err != nil {
		return nil, err
	}
	engine := &Engine{Root: cfg}
	engine.cache = map[string]*Builder{cfg.path: cfg}
	return engine, nil
}

func (e *Engine) Execute(ctx context.Context, commands []string) error {
	if len(commands) == 0 {
		if len(e.Root.DefaultCommand) == 0 {
			return errors.New("empty command")
		}
		commands = []string{e.Root.DefaultCommand}
	}
	var execCmds []string
	var ref = new(refCheck)
	for _, command := range commands {
		cmd, ok := e.Root.Command[command]
		if !ok {
			return &NotFoundCMdErr{
				NameSpace: e.Root.NameSpace,
				Command:   command,
			}
		}
		for _, s := range cmd.Shell {
			if s == "" {
				continue
			}
			// 引用关系
			if s[0] == '$' {
				result, err := e.ParseRef(ctx, e.Root, command, s[1:], ref)
				if err != nil {
					return err
				}
				ref.Clear()
				execCmds = append(execCmds, result...)
			} else {
				execCmds = append(execCmds, s)
			}
		}
	}
	for _, item := range execCmds {
		execCmd := pkg.ScanWorld(item)
		err := pkg.Execute(execCmd[0], execCmd[1:])
		if err != nil {
			return err
		}
	}
	return nil
}

type SyntaxErr string

func (s SyntaxErr) Error() string {
	return fmt.Sprintf("syntax error [ $command | $namespace.command ] err: %q", string(s))
}

type NotFoundCMdErr struct {
	NameSpace string
	Command   string
}

func (e *NotFoundCMdErr) Error() string {
	return fmt.Sprintf("namespace %q not found command %q", e.NameSpace, e.Command)
}

type CycleErr struct {
	SrcNameSpace string
	SrcCommand   string
	DstNameSpace string
	DstCommand   string
}

func (e *CycleErr) Error() string {
	return fmt.Sprintf("namespace: %q\n\tcommand: %q\nnamespace: %q\n\tcommand: %q \ncommand cycle not allowed", e.SrcNameSpace, e.SrcCommand, e.DstNameSpace, e.DstCommand)
}

func (e *Engine) ParseRef(ctx context.Context, root *Builder, parentCmd, shell string, ref *refCheck) ([]string, error) {
	args := strings.SplitN(shell, ".", 2)
	switch len(args) {
	case 1:
		if args[0] == "" {
			return nil, SyntaxErr(shell)
		}
		subCmd := args[0]
		subCmdObj, ok := root.Command[subCmd]
		if !ok {
			return nil, &NotFoundCMdErr{NameSpace: root.NameSpace, Command: subCmd}
		}
		if ref.LoadAndStore(root.NameSpace + "." + parentCmd + ":" + root.NameSpace + "." + subCmd) {
			return nil, &CycleErr{
				SrcNameSpace: root.NameSpace,
				SrcCommand:   parentCmd,
				DstNameSpace: root.NameSpace,
				DstCommand:   subCmd,
			}
		}
		result := make([]string, 0, len(subCmdObj.Shell))
		for _, s2 := range subCmdObj.Shell {
			if s2 == "" {
				continue
			}
			if s2[0] == '$' {
				ret, err := e.ParseRef(ctx, root, subCmd, s2[1:], ref)
				if err != nil {
					return nil, err
				}
				result = append(result, ret...)
			} else {
				result = append(result, s2)
			}
		}
		return result, nil
	case 2:
		if args[0] == "" || args[1] == "" {
			return nil, SyntaxErr(shell)
		}
	default:
		return nil, SyntaxErr(shell)
	}
	var (
		namespace = args[0]
		command   = args[1]
		path      string
	)
	if namespace == root.NameSpace {
		path = root.path
	} else {
		for _, s2 := range root.Import {
			if s2.Name == namespace {
				path = s2.Path.FromSlash()
			}
		}
	}
	if path == "" {
		repo := ctx.Value("repo").(*Repo)
		builder, ok := e.cache[repo.BuilderPath(namespace)]
		if !ok {
			var err error
			builder, err = repo.Find(namespace)
			if err != nil {
				return nil, fmt.Errorf("import %q not found %s", namespace, err.Error())
			}
		}
		path = builder.path
		e.cache[path] = builder
	}
	if ref.LoadAndStore(root.NameSpace + "." + parentCmd + ":" + namespace + "." + command) {
		return nil, &CycleErr{
			SrcNameSpace: root.NameSpace,
			SrcCommand:   parentCmd,
			DstNameSpace: namespace,
			DstCommand:   command,
		}
	}
	builder, ok := e.cache[path]
	if !ok {
		cfg, err := OpenBuilder(path)
		if err != nil {
			return nil, err
		}
		e.cache[path] = cfg
		builder = cfg
	}
	cmd, ok := builder.Command[command]
	if !ok {
		return nil, &NotFoundCMdErr{NameSpace: namespace, Command: command}
	}
	var result []string
	for _, s2 := range cmd.Shell {
		if s2 == "" {
			continue
		}
		if s2[0] == '$' {
			res, err := e.ParseRef(ctx, builder, command, s2[1:], ref)
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

type refCheck struct {
	path map[string]struct{}
}

func (r *refCheck) LoadAndStore(path string) bool {
	if r.path == nil {
		r.path = map[string]struct{}{}
	}
	_, pExist := r.path[path]
	if !pExist {
		r.path[path] = struct{}{}
		return false
	}
	return true
}

func (r *refCheck) Clear() {
	clear(r.path)
}
