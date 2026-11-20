package cmd

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/checkson-io/checkson-cli/operations/auth"
	"github.com/checkson-io/checkson-cli/services"
	"github.com/spf13/cobra"
)

func resourceCompletion(channels, flag bool) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, prefix string) ([]string, cobra.ShellCompDirective) {
		directive := cobra.ShellCompDirectiveNoFileComp
		if !flag && len(args) > 0 {
			return nil, directive
		}
		token, err := auth.ReadAuthToken()
		if err != nil {
			return nil, directive
		}
		devMode, _ := cmd.Flags().GetBool("dev-mode")
		ctx, cancel := context.WithTimeout(cmd.Context(), 2*time.Second)
		defer cancel()
		names, err := services.ListCompletionNames(ctx, token, devMode, channels)
		if err != nil {
			return nil, directive
		}
		var matches []string
		for _, name := range names {
			if strings.HasPrefix(name, prefix) {
				matches = append(matches, name)
			}
		}
		sort.Strings(matches)
		return matches, directive
	}
}

func registerResourceCompletions(root *cobra.Command) {
	for _, command := range root.Commands() {
		switch command.Name() {
		case "show", "edit", "delete", "runs", "logs":
			command.ValidArgsFunction = resourceCompletion(false, false)
		case "create":
			command.ValidArgsFunction = cobra.NoFileCompletions
			_ = command.RegisterFlagCompletionFunc("channel", resourceCompletion(true, true))
		case "channels":
			for _, child := range command.Commands() {
				switch child.Name() {
				case "show", "delete":
					child.ValidArgsFunction = resourceCompletion(true, false)
				case "create":
					child.ValidArgsFunction = cobra.NoFileCompletions
				}
			}
		}
	}
}
