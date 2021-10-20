package credentials

import (
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "credentials",
		Short: "Credentials related actions",
		Long: `Credentials related action`,
	}
)


func init() {
	Cmd.AddCommand(listCmd)
}
