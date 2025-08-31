package cmd

import (
	"fmt"
	"github.com/Li-giegie/builder/internal"
	"github.com/Li-giegie/builder/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "out repository builder namespace list",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("repo").(*internal.Repo)
		return repo.List()
	},
}

var (
	repoLoadCmdFlagNamespace string
	repoLoadCmdFlagForce     bool
)

var repoLoadCmd = &cobra.Command{
	Use:   "load [flags] Filename",
	Args:  cobra.MinimumNArgs(1),
	Short: "load builder file to repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("repo").(*internal.Repo)
		return repo.Load(args[0], repoLoadCmdFlagNamespace, repoLoadCmdFlagForce)
	},
}
var (
	repoSaveCmdFlagOut string
	repoSaveCmdForce   bool
)

var repoSaveCmd = &cobra.Command{
	Use:   "save [flags] Namespace",
	Args:  cobra.MinimumNArgs(1),
	Short: "save repository builder file to path",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("repo").(*internal.Repo)
		return repo.Save(args[0], repoSaveCmdFlagOut, repoSaveCmdForce)
	},
}
var repoRemoveCmd = &cobra.Command{
	Use:   "remove Namespace",
	Args:  cobra.MinimumNArgs(1),
	Short: "remove repository builder",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := cmd.Context().Value("repo").(*internal.Repo)
		return repo.Remove(args[0])
	},
}
var repoPathCmd = &cobra.Command{
	Use:   "path NewPath",
	Args:  cobra.MinimumNArgs(1),
	Short: "set repository path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := filepath.FromSlash(args[0])
		info, err := os.Stat(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		} else {
			if !info.IsDir() {
				return fmt.Errorf("%s is not a directory", path)
			}
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		if err = pkg.MkDir(ctx.Value("userConfigDir").(string)); err != nil {
			return err
		}
		conf := ctx.Value("conf").(map[string]any)
		conf["repoPath"] = path
		f, err := os.OpenFile(ctx.Value("userConfigFile").(string), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		return yaml.NewEncoder(f).Encode(conf)
	},
}

func init() {
	repoLoadCmd.Flags().BoolVarP(&repoLoadCmdFlagForce, "force", "f", false, "force overwrite")
	repoLoadCmd.Flags().StringVarP(&repoLoadCmdFlagNamespace, "namespace", "n", "", "namespace")
	repoSaveCmd.Flags().StringVarP(&repoSaveCmdFlagOut, "out", "o", "./", "out file")
	repoSaveCmd.Flags().BoolVarP(&repoSaveCmdForce, "force", "f", false, "force overwrite")
	repoCmd.AddCommand(repoLoadCmd, repoSaveCmd, repoRemoveCmd, repoPathCmd)
	rootCmd.AddCommand(repoCmd)
}
