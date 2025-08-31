package internal

import (
	"errors"
	"github.com/Li-giegie/builder/pkg"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Repo struct {
	Root string
}

func (r *Repo) Load(name, namespace string, force bool) error {
	srcFile := name
	builder, err := OpenBuilder(srcFile)
	if err != nil {
		return err
	}
	if builder.NameSpace = pkg.DefaultStr(namespace, builder.NameSpace); builder.NameSpace == "" {
		return errors.New(srcFile + " namespace is required")
	}
	if err = pkg.MkDir(r.Root); err != nil {
		return err
	}
	builderFileName := filepath.Join(r.Root, builder.NameSpace+".yaml")
	if pkg.IsExist(builderFileName) && !force {
		return errors.New("namespace \"" + builder.NameSpace + "\" repository already exists")
	}
	return pkg.SaveYamlObj(builderFileName, builder)
}

func (r *Repo) Save(namespace, out string, force bool) error {
	srcPath := filepath.Join(r.Root, namespace+".yaml")
	if !pkg.IsExist(srcPath) {
		return errors.New(namespace + " is not found in the repository")
	}
	dir, file := filepath.Split(out)
	if dir == "" {
		dir = "./"
	} else {
		if err := pkg.MkDir(dir); err != nil {
			return err
		}
	}
	if file == "" {
		file = namespace + ".yaml"
	}
	outFile := filepath.Join(dir, file)
	if pkg.IsExist(outFile) && !force {
		return errors.New("out path \"" + namespace + ".yaml\" already exists")
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func (r *Repo) Remove(namespace string) error {
	err := os.Remove(filepath.Join(r.Root, namespace+".yaml"))
	if err != nil {
		if os.IsExist(err) {
			return err
		}
		return errors.New(namespace + " is not found in the repository")
	}
	return nil
}

func (r *Repo) List() error {
	dirs, err := os.ReadDir(r.Root)
	if err != nil {
		if os.IsExist(err) {
			return err
		}
		return nil
	}
	for _, d := range dirs {
		if !d.IsDir() {
			println(strings.TrimSuffix(d.Name(), ".yaml"))
		}
	}
	return nil
}

func (r *Repo) Find(namespace string) (*Builder, error) {
	file := r.BuilderPath(namespace)
	if !pkg.IsExist(file) {
		return nil, errors.New(namespace + " is not found in the repository")
	}
	return OpenBuilder(file)
}

func (r *Repo) BuilderPath(name string) string {
	return filepath.Join(r.Root, name+".yaml")
}
