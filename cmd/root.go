package cmd

import (
	"context"
	"github.com/Li-giegie/builder/internal"
	"github.com/Li-giegie/builder/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var rootCmdFlagCfgName string
var rootCmd = &cobra.Command{
	Use:  "builder [flags] commands...",
	Args: cobra.OnlyValidArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		eng, err := internal.NewEngine(filepath.FromSlash(rootCmdFlagCfgName))
		if err != nil {
			return err
		}
		return eng.Execute(cmd.Context(), args)
	},
}

func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	err := rootCmd.Execute()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func init() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	userConfigFile := filepath.Join(userConfigDir, ".builder", "config.yaml")
	conf := map[string]any{"repoPath": filepath.Join(userCacheDir, "builder")}
	if pkg.IsExist(userConfigFile) {
		f, err := os.Open(userConfigFile)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		err = yaml.NewDecoder(f).Decode(&conf)
		_ = f.Close()
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}
	ctx := context.WithValue(context.TODO(), "repo", &internal.Repo{Root: conf["repoPath"].(string)})
	ctx = context.WithValue(ctx, "userCacheDir", userCacheDir)
	ctx = context.WithValue(ctx, "userConfigDir", filepath.Join(userConfigDir, ".builder"))
	ctx = context.WithValue(ctx, "userConfigFile", userConfigFile)
	ctx = context.WithValue(ctx, "conf", conf)
	rootCmd.SetContext(ctx)
	rootCmd.Flags().StringVarP(&rootCmdFlagCfgName, "config", "c", "./.builder.yaml", "config file path")
}
