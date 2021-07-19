package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Inject the version with GoReleaser
	Version = "unset"
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
