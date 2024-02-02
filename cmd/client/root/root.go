package root

import (
	"github.com/faireal/kobe/cmd/client/adhoc"
	"github.com/faireal/kobe/cmd/client/playbook"
	"github.com/faireal/kobe/cmd/client/project"
	"github.com/faireal/kobe/cmd/client/task"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "kobe",
	Short: "A kobe client cli",
}

func init() {
	Cmd.AddCommand(project.Cmd)
	Cmd.AddCommand(playbook.Cmd)
	Cmd.AddCommand(task.Cmd)
	Cmd.AddCommand(adhoc.Cmd)
}
