package runserver

import (
	"github.com/RealLiuSha/echo-admin/bootstrap"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/spf13/cobra"

	"go.uber.org/fx"
)

var configFile string
var casbinModel string

func init() {
	pf := StartCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c",
		"config/config.yaml", "this parameter is used to start the service application")

	pf.StringVarP(&casbinModel, "casbin_model", "m",
		"config/casbin_model.conf", "this parameter is used for the running configuration of casbin")

	cobra.MarkFlagRequired(pf, "config")
	cobra.MarkFlagRequired(pf, "casbin_model")
}

var StartCmd = &cobra.Command{
	Use:          "runserver",
	Short:        "Start API server",
	Example:      "{execfile} server -c config/config.yaml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		lib.SetConfigPath(configFile)
		lib.SetConfigCasbinModelPath(casbinModel)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	fx.New(bootstrap.Module, fx.NopLogger).Run()
}
