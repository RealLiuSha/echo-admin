package cmd

import (
	"errors"
	"os"

	"github.com/RealLiuSha/echo-admin/cmd/migrate"
	"github.com/RealLiuSha/echo-admin/cmd/runserver"
	"github.com/RealLiuSha/echo-admin/cmd/setup"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runserver.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(setup.StartCmd)
}

var rootCmd = &cobra.Command{
	Use:          "echo-admin",
	Short:        "echo-admin",
	SilenceUsage: true,
	Long:         `echo-admin`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(
				"requires at least one arg, " +
					"you can view the available parameters through `--help`",
			)
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run:               func(cmd *cobra.Command, args []string) {},
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
