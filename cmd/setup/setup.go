package setup

import (
	"os"

	"github.com/RealLiuSha/echo-admin/models"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/RealLiuSha/echo-admin/api/repository"
	"github.com/RealLiuSha/echo-admin/api/services"
	"github.com/RealLiuSha/echo-admin/lib"
	"github.com/RealLiuSha/echo-admin/pkg/file"
)

var configFile string
var menuFile string

func init() {
	pf := StartCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c",
		"config/config.yaml", "this parameter is used to start the service application")
	pf.StringVarP(&menuFile, "menu", "m",
		"config/menu.yaml", "this parameter is used to set the initialized menu data.")

	cobra.MarkFlagRequired(pf, "config")
	cobra.MarkFlagRequired(pf, "menu")
}

var StartCmd = &cobra.Command{
	Use:          "setup",
	Short:        "Set up data for the application",
	Example:      "{execfile} init -c config/settings.yml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		lib.SetConfigPath(configFile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		config := lib.NewConfig()
		logger := lib.NewLogger(config)
		db := lib.NewDatabase(config, logger)

		menuService := services.NewMenuService(
			logger,
			repository.NewMenuRepository(db, logger),
			repository.NewMenuActionRepository(db, logger),
			repository.NewMenuActionResourceRepository(db, logger),
		)

		if !file.IsFile(menuFile) {
			logger.Zap.Fatal("menu file does not exist")
		}

		fs, err := os.Open(menuFile)
		if err != nil {
			logger.Zap.Fatalf("menu file could not be opened: %v", err)
		}

		defer fs.Close()

		var menuTrees models.MenuTrees
		yd := yaml.NewDecoder(fs)
		if err = yd.Decode(&menuTrees); err != nil {
			logger.Zap.Fatalf("menu file decode error: %v", err)
		}

		if err = menuService.CreateMenus("", menuTrees); err != nil {
			logger.Zap.Fatalf("menu file init err: %v", err)
		}

		logger.Zap.Info("menu file import successfully")
	},
}
