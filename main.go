package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//global flag
var logLevel string

var RootCmd = &cobra.Command{
	Use:   "cflow",
	Short: "Contracts in the Command Line",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set the log level
	},
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New Contract",
	RunE:  newContract,
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile into readable contract",
	RunE:  compileContract,
}

func init() {
	//parse flag and set config
	RootCmd.PersistentFlags().StringVar(&logLevel, "log_level", "info", "Log level")

	RootCmd.AddCommand(newCmd)
	RootCmd.AddCommand(compileCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
