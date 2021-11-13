package mmp

import (
	"github.com/nlamot/nero/cmd/mmp/credentials"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "mmp",
		Short: "MMP related actions",
		Long:  `Automates MMP related administrative actions`,
	}
)

func init() {
	Cmd.AddCommand(credentials.Cmd)
}
