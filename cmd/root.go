package cmd

import (
	"github.com/Li-giegie/builder/internal"
	"github.com/spf13/cobra"
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
		return eng.Execute(args)
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
	rootCmd.Flags().StringVarP(&rootCmdFlagCfgName, "config", "c", "./.builder.yaml", "config file path")
}
