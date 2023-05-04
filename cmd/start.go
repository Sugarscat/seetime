package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sugarscat/seetime/server"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start SeeTime server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	server.Loading()
}

func init() {
	RootCmd.AddCommand(StartCmd)
}
