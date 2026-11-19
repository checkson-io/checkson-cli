package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"testing"
)

func TestResourceCompletionRegistration(t *testing.T) {
	root := NewChecksonCommand(output.DefaultIOStreams())
	for _, path := range [][]string{{"show"}, {"delete"}, {"runs"}, {"logs"}, {"channels", "show"}, {"channels", "delete"}} {
		command, _, err := root.Find(path)
		if err != nil || command.ValidArgsFunction == nil {
			t.Fatalf("missing completion for %v: %v", path, err)
		}
		names, directive := command.ValidArgsFunction(command, []string{"already-completed"}, "")
		if len(names) != 0 || directive != cobra.ShellCompDirectiveNoFileComp {
			t.Fatalf("unexpected extra argument completions for %v", path)
		}
	}
	create, _, _ := root.Find([]string{"create"})
	if completion, ok := create.GetFlagCompletionFunc("channel"); !ok || completion == nil {
		t.Fatal("missing channel flag completion")
	}
	// Also register edit when the independent edit feature is present.
	edit := &cobra.Command{Use: "edit CHECK"}
	root.AddCommand(edit)
	registerResourceCompletions(root)
	if edit.ValidArgsFunction == nil {
		t.Fatal("missing edit completion")
	}
}
