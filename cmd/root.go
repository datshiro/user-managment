package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var (
	rootCmd = &cobra.Command{
		Use:   "cake",
		Short: "Cake is a digital bank by VP Bank",
		Long:  " Cake is a digital bank by VP Bank",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
