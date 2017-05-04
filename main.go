package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// global flag
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

func newContract(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("new expects two args: directory for new engagement and template path")
	}

	if err := newEngagement(args[0], args[1]); err != nil {
		return err
	}
	return nil
}

func compileContract(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("compile expects one arg: name")
	}

	if outputType == "" {
		return fmt.Errorf("must set output type with --output")
	}

	if err := generateContract(args[0], outputType); err != nil {
		return err
	}
	return nil
}
