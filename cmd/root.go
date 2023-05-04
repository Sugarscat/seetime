package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "SeeTime",
	Short: "A software that automates scheduled tasks.",
	Long: `A program that supports the execution of scheduled tasks,
built with love by Sugarscat and friends in Go.
Complete documentation is available at https://seetime.sugarscat.com/`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {}
