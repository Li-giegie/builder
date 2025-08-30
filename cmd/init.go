package cmd

import (
	"github.com/Li-giegie/builder/internal"
	"github.com/spf13/cobra"
)

var (
	initCmdFlagOut   string
	initCmdFlagForce bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a builder config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.DefaultBuilder(initCmdFlagOut, initCmdFlagForce)
	},
}

func init() {
	initCmd.Flags().StringVarP(&initCmdFlagOut, "out", "o", "./.builder.yaml", "output file name")
	initCmd.Flags().BoolVarP(&initCmdFlagForce, "force", "f", false, "force overwrite")
	rootCmd.AddCommand(initCmd)
}
