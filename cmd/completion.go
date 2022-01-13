package cmd

import (
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

func newCompletionCmd() *cobra.Command {
	var cmdCompletion = &cobra.Command{
		Use: "completion [bash|zsh|fish|powershell]",
		Long: `To load completions:

Bash:

$ source <(checkson completion bash)

# To load completions for each session, execute once:
Linux:
  $ checkson completion bash > /etc/bash_completion.d/checkson
MacOS:
  $ checkson completion bash > /usr/local/etc/bash_completion.d/checkson

Zsh:

$ source <(checkson completion zsh)

# To load completions for each session, execute once:
$ checkson completion zsh > "${fpath[1]}/_checkson"

Fish:

$ checkson completion fish | source

# To load completions for each session, execute once:
$ checkson completion fish > ~/.config/fish/completions/checkson.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var err error

			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				err = cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				err = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				err = errors.Errorf("Unsupported shell type %q.", args[0])
			}

			if err != nil {
				output.Fail(errors.Wrap(err, "unable to generate shell completion"))
			}
		},
	}

	return cmdCompletion
}
