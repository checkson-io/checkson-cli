package cmd

import (
	"github.com/checkson-io/checkson-cli/output"
	"github.com/spf13/cobra"
)

var Verbose bool
var DevMode bool

func NewChecksonCommand(streams output.IOStreams) *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "checkson",
		Short: "command-line interface for Checkson",
		Long:  `A command-line interface for managing Checkson checks`,
	}

	rootCmd.AddCommand(newCompletionCmd())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newListChecksCmd())
	rootCmd.AddCommand(newStatusChecksCmd())
	rootCmd.AddCommand(newListRunsCmd())
	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newLogoutCmd())
	rootCmd.AddCommand(newLogsCmd())
	rootCmd.AddCommand(newCreateCheckCmd())
	rootCmd.AddCommand(newEditCheckCmd())
	rootCmd.AddCommand(newDeleteCheckCmd())
	rootCmd.AddCommand(newShowCheckCmd())
	rootCmd.AddCommand(newChannelCmd())

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dev-mode", false, "enable dev mode, communicating to local services")
	registerResourceCompletions(rootCmd)

	output.IoStreams = streams
	rootCmd.SetOut(streams.Out)
	rootCmd.SetErr(streams.ErrOut)
	return rootCmd
}
