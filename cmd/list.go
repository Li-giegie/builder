package cmd

import (
	"github.com/Li-giegie/builder/internal"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var (
	listCmdFlagCfgName string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list builder commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		eng, err := internal.NewEngine(filepath.FromSlash(rootCmdFlagCfgName))
		if err != nil {
			return err
		}
		keys := make([]string, 0, len(eng.Root.Command))
		for s := range eng.Root.Command {
			keys = append(keys, s)
		}
		sort.Strings(keys)
		println("builder commands:")
		for _, key := range keys {
			entry := eng.Root.Command[key]
			println(" ", key, "  \t", entry.Desc)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&listCmdFlagCfgName, "config", "c", "./.builder.yaml", "config file path")
	rootCmd.AddCommand(listCmd)
}
