package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "v0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("azverify version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
