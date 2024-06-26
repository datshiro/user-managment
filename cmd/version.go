package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of app",
	Long:  `All software has versions. This is an api server version`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("App 1.0.0")
  },
}
