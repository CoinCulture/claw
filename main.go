package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// global flags
var outputType string

var RootCmd = &cobra.Command{
	Use:   "claw",
	Short: "Command Line Law",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
	compileCmd.Flags().StringVar(&outputType, "output", "md", "Output type: md | html | pdf")

	RootCmd.AddCommand(newCmd)
	RootCmd.AddCommand(compileCmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
